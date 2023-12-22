package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rosariocannavo/go_auth/config"
	"github.com/rosariocannavo/go_auth/internal/circuit_breaker"
	"github.com/rosariocannavo/go_auth/internal/models"
	"github.com/rosariocannavo/go_auth/internal/nats"
)

// TODO: put initialization in a separated file
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
	capturedStatus := rrw.StatusCode

	message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "proxy_handler/handleResponse", capturedStatus, rrw.Body)
	nats.NatsConnection.PublishMessage(message)

	return capturedStatus
}

// wrap handler inside circuit breaker
func ProxyHandler(c *gin.Context) {

	remote, err := url.Parse(config.ProxyDestination)

	if err != nil {
		panic(err)
	}

	proxy := createReverseProxy(remote, c.Request.Header, c.Param("proxyPath")) //redirect the request to  the proxy

	_, errcb := circuit_breaker.CircuitBreaker.Execute(func() (interface{}, error) { //circuite breaker here

		status := handleResponse(proxy, c.Writer, c.Request)

		if status < 200 || status >= 300 {
			return nil, errors.New("server error")
		}

		return nil, nil

	})

	if errcb != nil {
		fmt.Println("circuit breaker error", errcb)
	}

}
