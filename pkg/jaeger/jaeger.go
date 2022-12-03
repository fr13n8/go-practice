package jaeger

import (
	"io"

	"github.com/fr13n8/go-practice/internal/config"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

// Init Jaeger
func InitJaeger(cfg *config.JaegerConfig, p string) (opentracing.Tracer, io.Closer, error) {
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: cfg.ServiceName + "-" + p,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           cfg.LogSpans,
			LocalAgentHostPort: cfg.Host + ":" + cfg.Port,
		},
	}

	return jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
}
