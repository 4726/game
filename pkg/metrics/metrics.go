package metrics

//maybe implement error logging

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HTTP is an http server
type HTTP struct {
	s *http.Server
}

// NewHTTP returns a new http
func NewHTTP() *HTTP {
	return &HTTP{nil}
}

// Run starts the http server on the specified port
func (h *HTTP) Run(port int) error {
	h.s = &http.Server{Addr: fmt.Sprintf(":%v", port)}
	h.s.Handler = promhttp.Handler()
	return h.s.ListenAndServe()
}

// Close gracefully stops the server
func (h *HTTP) Close() {
	if h.s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		h.s.Shutdown(ctx)
	}
}
