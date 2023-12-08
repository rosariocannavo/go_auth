package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rosariocannavo/go_auth/internal/models"
	"github.com/sony/gobreaker"
)

var cb *gobreaker.CircuitBreaker

// TODO: separate cb and proxy
func init() {
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "My Circuit Breaker",
		MaxRequests: 0,               // Maximum number of consecutive failures before tripping the circuit
		Interval:    5 * time.Second, // Duration to wait before allowing another request
		Timeout:     2 * time.Second, // Timeout for a single request attempt
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			// return counts.Requests >= 3 && failureRatio >= 0.6 // Trips the circuit if failure rate exceeds 60%
			return counts.ConsecutiveFailures > 3
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("Circuit breaker '%s' changed from '%s' to '%s'\n", name, from, to)
		},
	})
}

func handler(proxy *httputil.ReverseProxy, w http.ResponseWriter, r *http.Request) int {
	rrw := models.NewResponseRecorderWriter(w)
	proxy.ServeHTTP(rrw, r)
	capturedResponse := rrw.Body.String()
	capturedStatus := rrw.StatusCode
	fmt.Println("caputer response", capturedResponse)
	return capturedStatus
}

func ProxyHandler(c *gin.Context) {
	remote, err := url.Parse("http://localhost:8081") //original server back the proxy TODO: MAKE THIS AS A VARIABLE
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("proxyPath")

	}

	_, errcb := cb.Execute(func() (interface{}, error) {
		status := handler(proxy, c.Writer, c.Request)
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
