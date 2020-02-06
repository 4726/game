package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HTTP struct {
	s *http.Server
}

func NewHTTP() *HTTP {
	return &HTTP{nil}
}

func (h *HTTP) Run(port int) error {
	h.s = &http.Server{Addr: fmt.Sprintf(":%v", port)}
	h.s.Handler = promhttp.Handler()
	return h.s.ListenAndServe()
}

func (h *HTTP) Close() {
	if h.s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		h.s.Shutdown(ctx)
	}
}
