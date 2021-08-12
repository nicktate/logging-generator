package generator

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultKillPeriod represents the default kill period.
	DefaultKillPeriod = time.Minute
	// DefaultMessageSize represents the default message size.
	DefaultMessageSize = 1 << (10 * 1)
	// DefaultMessageTotal represents the default message total.
	DefaultMessageTotal = 30
	// DefaultMessagePeriod represents the default message period.
	DefaultMessagePeriod = time.Second
)

// Generator represents a generator to run
type Generator interface {
	Run(context.Context) error

	Collect(ch chan<- prometheus.Metric)
	Describe(ch chan<- *prometheus.Desc)
}

// Generator represents a logging generator.
type generator struct {
	writer        io.Writer
	killPeriod    time.Duration
	messageSize   int64
	messageTotal  int64
	messagePeriod time.Duration

	*metrics
}

type writer struct {
	l logrus.FieldLogger
}

func (w *writer) Write(in []byte) (int, error) {
	w.l.Info(in)
	return 0, nil
}

func (w *writer) WriteString(in string) (int, error) {
	w.l.Info(in)
	return 0, nil
}

// Option represents options for the generator
type Option func(*generator)

// WithLog provides the logger for the generator.
func WithLog(l logrus.FieldLogger) func(*generator) {
	return func(g *generator) {
		g.writer = &writer{
			l,
		}
	}
}

// WithKillPeriod sets the kill period for the generator.
func WithKillPeriod(p time.Duration) func(*generator) {
	return func(g *generator) {
		g.killPeriod = p
	}
}

// WithMessageSize sets the kill period for the generator.
func WithMessageSize(p int64) func(*generator) {
	return func(g *generator) {
		g.messageSize = p
	}
}

// WithMessageTotal sets the kill period for the generator.
func WithMessageTotal(p int64) func(*generator) {
	return func(g *generator) {
		g.messageTotal = p
	}
}

// WithMessagePeriod sets the kill period for the generator.
func WithMessagePeriod(p time.Duration) func(*generator) {
	return func(g *generator) {
		g.messagePeriod = p
	}
}

// NewGenerator returns a new logging generator with the provided options.
func NewGenerator(options ...Option) Generator {
	g := &generator{
		os.Stdout,
		DefaultKillPeriod,
		DefaultMessageSize,
		DefaultMessageTotal,
		DefaultMessagePeriod,
		&metrics{
			logLinesTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "logging_generator",
				Name:      "lines_total",
				Help:      "Total number of lines emitted.",
			}, []string{}),
			logBytesTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "logging_generator",
				Name:      "bytes_total",
				Help:      "Total number of bytes emitted.",
			}, []string{}),
		},
	}

	for _, o := range options {
		o(g)
	}

	return g
}

// Run starts the logging generator.
func (g *generator) Run(ctx context.Context) error {
	var total int64
	kill := time.NewTimer(g.killPeriod)
	defer kill.Stop()
	period := time.NewTicker(g.messagePeriod)
	defer period.Stop()
	for ctx.Err() == nil && total < g.messageTotal {
		select {
		case <-kill.C:
			return nil
		case <-period.C:
			// continue
		}

		err := writePassages(g.messageSize, g.writer)
		if err != nil {
			return err
		}

		total++
		g.metrics.logLinesTotal.With(prometheus.Labels{}).Inc()
		g.metrics.logBytesTotal.With(prometheus.Labels{}).Add(float64(g.messageSize))
	}

	select {
	case <-ctx.Done():
	case <-kill.C:
	}

	return ctx.Err()
}
