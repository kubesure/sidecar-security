package middleware

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"

	api "github.com/kubesure/sidecar-security/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var customerDataSvc = os.Getenv("CUSTOMER_DATA_SVC")
var fraudCheckTCPSvc = os.Getenv("CUSTOMER_DATA_SVC")

//initializes logurs with info level
func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

//Customer type represents a customer
type Customer struct {
	accountNumer string
	CIF          int64
}

//FraudCheck message is sent to fraud checker to run fraud checks on the request.
type fraudCheck struct {
	smhSegIDVersion string
	smhMsgVersion   string
	smhTranType     string
	smhCustType     string
	smhActType      string
	smhSource       string
	fromAccount     string
	clientIP        string
	customerID      string
}

//TCP message from FraudCheck is parsed into fraudCheckRes
type fraudCheckRes struct {
	isOk bool
}

//TimeoutHandler is a customer timeout handler which return 504 when
//middlewares or origin does not respond with http.Server.WriteTimeout
func TimeoutHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGatewayTimeout)
	})
}

//Logger middleware logs orgin's request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("logging: request middleware")
		next.ServeHTTP(w, r)
		log.Infof("logging: response middleware")
	})
}

//Auth middleware authenticates request
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Authenticating request: ")
		if r.Header.Get("user") != "foo" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Infof("Auth: Pass")
		next.ServeHTTP(w, r)

	})
}

//Final middleware forwards request to orgin
func Final(proxy *httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Passing call to origin")
		proxy.ServeHTTP(w, r)
	})
}

//FraudChecker middleware checks if the request is fradulent.
func FraudChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Fraud checking request...")

		//reterive customerData from GRCP service Customer.getCustomer
		c, cerr := customerData(r)
		if cerr != nil {
			log.Errorf("Error getting customer data %v", cerr)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		//Make the message for Fraud checking TCP service
		msg, merr := makeFTCPMessage(r, c)
		if merr != nil {
			log.Errorf("error while making tcp message %v", merr)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		//Create a TCP connection to Fraud checking service
		conn, cerr := net.Dial("tcp", fraudCheckTCPSvc+":8080")
		defer conn.Close()
		if cerr != nil {
			log.Errorf("Error while connecting to Fraud Server")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		_, werr := conn.Write([]byte(*msg))
		if werr != nil {
			log.Errorf("Error while sending message to TCP server %v", werr)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		//Read message from TCP service until EOF
		tcpmsg, rerr := ioutil.ReadAll(conn)
		if rerr != nil {
			log.Errorf("Error while reading message to TCP server %v", rerr)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		//Parse response Fraud service
		fcheck, ferr := parseFTCPResponse(string(tcpmsg))
		if ferr != nil {
			log.Errorf("Error while reading message to TCP server %v", rerr)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		//Check if request is fraudulent
		if !fcheck.isOk {
			log.Infof("Fraudulent request received")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//Dispatch to next middleware check
		next.ServeHTTP(w, r)
	})
}

//make the TCP Fraud check message from request and Customer data
func makeFTCPMessage(r *http.Request, c *Customer) (*string, error) {
	message := string("smh_seg_id_version:000004|smh_source:")
	return &message, nil
}

//Pulls customer data from Customer.getCustomer GRCP service
func customerData(r *http.Request) (*Customer, error) {
	conn, err := grpc.Dial(customerDataSvc+":50051", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewCustomerClient(conn)
	customer, err := makeCustomerData(r, client)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

//Makes the customer data using customer data grpc service input for service is
//read from request header and body. Grpc service return datas cached in a in transient store.
func makeCustomerData(r *http.Request, client api.CustomerClient) (*Customer, error) {
	req := &api.CustomerRequest{Version: "v1", AccountNumber: "12345"}
	res, err := client.GetCustomer(context.Background(), req)

	if err != nil {
		return nil, err
	}
	c := &Customer{}
	c.CIF = res.CIF
	return c, nil
}

//Parses the TCP response from TCP Fraud check service
func parseFTCPResponse(msg string) (*fraudCheckRes, error) {
	log.Infof("parsing %v", msg)
	return &fraudCheckRes{isOk: true}, nil
}
