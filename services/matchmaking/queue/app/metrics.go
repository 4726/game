package app

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	inQueueTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "queue_in_queue_total",
		Help: "The total number of users in the queue",
	})
)
