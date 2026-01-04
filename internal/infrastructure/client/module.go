package client

import (
	"net/http"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	attributeapi "github.com/Sokol111/ecommerce-attribute-service-api/gen/httpapi"
	httpclient "github.com/Sokol111/ecommerce-commons/pkg/http/client"
)

// AttributeClientModule provides AttributeClient with its dependencies
func AttributeClientModule() fx.Option {
	return fx.Module("attribute-client",
		fx.Provide(
			fx.Private,
			httpclient.ProvideHTTPClient("attribute-service"),
		),
		fx.Provide(provideApiClient),
		fx.Provide(newAttributeClient),
	)
}

func provideApiClient(
	httpClient *http.Client,
	cfg httpclient.ClientConfig,
	tracerProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
) (*attributeapi.Client, error) {
	return attributeapi.NewClient(
		cfg.BaseURL,
		attributeapi.WithClient(httpClient),
		attributeapi.WithTracerProvider(tracerProvider),
		attributeapi.WithMeterProvider(meterProvider),
	)
}
