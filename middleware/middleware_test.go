package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	api "github.com/kubesure/sidecar-security/api"
)

const mockBody string = `{"accountNumber": 1234}`

//TestCustomer data by making just mock http request
func TestCustomerData(t *testing.T) {
	r := httptest.NewRequest("GET", "/", strings.NewReader(mockBody))
	r.Header.Add("foo", "bar")
	c, err := customerData(r)
	if err != nil {
		t.Errorf("Some error occured")
	}

	if c.CIF != 1234 {
		t.Errorf("excpted 1234 got %v", c.CIF)
	}
	t.Logf("Customer CIF: %v", c.CIF)
}

//Tests makeCustomer Data by creating mock request and mock grpc Customer.getCustomer service response
func TestGetGRCPCustomerData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockclient := api.NewMockCustomerClient(ctrl)
	mockclient.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).Return(&api.CustomerResponse{CIF: 77777}, nil)
	testResponse(t, mockclient)
}

func testResponse(t *testing.T, client api.CustomerClient) {
	res, err := makeCustomerData(makeMockProxyRequest(), client)
	if err != nil {
		t.Errorf("Some error occured")
	}

	if res.CIF != 77777 {
		t.Errorf("excpted 77777 got %v", res.CIF)
	}
	t.Logf("Customer CIF: %v", res.CIF)
}

func makeMockProxyRequest() *http.Request {
	r := httptest.NewRequest("GET", "/", strings.NewReader(mockBody))
	r.Header.Add("foo", "bar")
	return r
}

//reference code

/*func TestGetGRCPCustomerData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockclient := api.NewMockCustomerClient(ctrl)
	//cr := api.CustomerRequest{AccountNumber: "12345"}
	mockclient.EXPECT().GetCustomer(gomock.Any(), gomock.Any()).Return(&api.CustomerResponse{CIF: 77777}, nil)
	testResponse(t, mockclient)
}

func testResponse(t *testing.T, client api.CustomerClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	req := api.CustomerRequest{AccountNumber: "12345"}
	res, err := client.GetCustomer(ctx, &req)
	if err != nil {
		t.Errorf("Some error occured")
	}

	if res.CIF != 77777 {
		t.Errorf("excpted 77777 got %v", res.CIF)
	}
	t.Logf("Customer CIF: %v", res.CIF)
}*/
