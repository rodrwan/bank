package tracer

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"
)

// NewTracer ...
func NewTracer() (opentracing.Tracer, io.Closer, error) {
	// load config from environment variables
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, nil, err
	}

	// create tracer from config
	return cfg.NewTracer(
		config.Metrics(jprom.New()),
	)
}
