package tracer

import (
	"go-social-network/configs"

	"github.com/hertz-contrib/obs-opentelemetry/provider"
)

func InitTracer(serviceName string) provider.OtelProvider {
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(serviceName),
		provider.WithExportEndpoint(configs.Data().Jaeger.Addr),
		provider.WithInsecure(),
	)
	return p
}
