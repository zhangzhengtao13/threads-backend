package tracer

import (
	"io"
	"time"

	trace "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func NewJaegerTracer(serviceName, agentHostPort string) (trace.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  agentHostPort,
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}
	trace.SetGlobalTracer(tracer)
	return tracer, closer, nil
}
