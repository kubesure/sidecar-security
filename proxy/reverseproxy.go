package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/kubesure/sidecar-security/middleware"
	"github.com/kubesure/sidecar-security/routing"
)

//SetupProxy configures reverse proxies
func SetupProxy(router *httprouter.Router) {
	origin, _ := url.Parse("http://localhost:9000/")
	path := "/*catchall"

	proxy := httputil.NewSingleHostReverseProxy(origin)
	proxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = origin.Scheme
		req.URL.Host = origin.Host
	}

	for _, config := range routing.Configurations {
		router.Handler(config.Method, path, middleware.Logger(middleware.Auth(middleware.Final(proxy))))
	}
}
