package generator

import (
	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	logLinesTotal *prometheus.CounterVec
	logBytesTotal *prometheus.CounterVec
}

func (m *metrics) Collect(ch chan<- prometheus.Metric) {
	m.logLinesTotal.Collect(ch)
	m.logBytesTotal.Collect(ch)
}

func (m *metrics) Describe(ch chan<- *prometheus.Desc) {
	m.logLinesTotal.Describe(ch)
	m.logBytesTotal.Describe(ch)
}
