package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/rosariocannavo/go_auth/config"
	"github.com/rosariocannavo/go_auth/internal/circuit_breaker"
	"github.com/rosariocannavo/go_auth/internal/models"
)

func createReverseProxy(remote *url.URL, headers http.Header, proxyPath string) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(remote)

	proxy.Director = func(req *http.Request) {
		req.Header = headers
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = proxyPath
	}

	return proxy
}

// real handler
func handleResponse(proxy *httputil.ReverseProxy, w http.ResponseWriter, r *http.Request) int {
	rrw := models.NewResponseRecorderWriter(w)
	proxy.ServeHTTP(rrw, r)
	capturedResponse := rrw.Body.String()
	capturedStatus := rrw.StatusCode
	fmt.Println("captuter response", capturedResponse)
	return capturedStatus
}

// wrap handler inside circuit breaker
func ProxyHandler(c *gin.Context) {

	remote, err := url.Parse(config.ProxyDestination)

	if err != nil {
		panic(err)
	}

	proxy := createReverseProxy(remote, c.Request.Header, c.Param("proxyPath"))

	_, errcb := circuit_breaker.CircuitBreaker.Execute(func() (interface{}, error) { //circuite breaker here

		status := handleResponse(proxy, c.Writer, c.Request)
		fmt.Println("captured status ", status)

		if status < 200 || status >= 300 {
			return nil, errors.New("server error")
		}

		return nil, nil

	})

	if errcb != nil {
		fmt.Println("circuit breaker error", errcb)
	}

}
