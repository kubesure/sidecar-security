package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

//FraudCheck message is sent to fraud checker to run fraud checks on the request.
type FraudCheck struct {
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
		//log.Infof("Authenticating request: %v", r)
		if r.Header.Get("user") != "foo" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
		log.Infof("Auth: Pass")
	})
}

//Final middleware forwards request to orgin
func Final(proxy *httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Passing call to origin")
		proxy.ServeHTTP(w, r)
	})
}

//FraudChecker checks if the request is fradulent.
func FraudChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Fraud checking request...")
		msg, msgerr := makeTCPMessage(r)
		if msgerr != nil {
			log.Errorf("Erorr while generating message %v", msgerr)
		}

		conn, cerr := net.Dial("tcp", ":8080")
		defer conn.Close()
		if cerr != nil {
			log.Errorf("Error while connecting to Fraud Server")
		}

		_, werr := conn.Write([]byte(*msg))
		if werr != nil {
			log.Errorf("Error while sending message to TCP server %v", werr)
		}

		reply := make([]byte, 1024)
		tcpmsg, rerr := conn.Read(reply)
		if rerr != nil {
			log.Errorf("Error while sending message to TCP server %v", rerr)
			return
		}
		log.Printf("EOF message received")
		log.Printf(string(tcpmsg))
		next.ServeHTTP(w, r)

		/*res := bufio.NewReader(conn)
		tcpmsg, rerr := res.ReadBytes(byte('\n'))
		if rerr == io.EOF {
			log.Printf("EOF message received")
			log.Printf(string(tcpmsg))
			next.ServeHTTP(w, r)
			return
		}*/
	})
}

func makeTCPMessage(r *http.Request) (*string, error) {
	message := string("smh_seg_id_version:000004|smh_source:")
	return &message, nil
}
