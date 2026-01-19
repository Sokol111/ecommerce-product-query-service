package client

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	catalogapi "github.com/Sokol111/ecommerce-catalog-service-api/gen/httpapi"
	httpclient "github.com/Sokol111/ecommerce-commons/pkg/http/client"
)

// serviceSecuritySource provides service-to-service authentication
type serviceSecuritySource struct{}

func (s *serviceSecuritySource) BearerAuth(ctx context.Context, operationName catalogapi.OperationName) (catalogapi.BearerAuth, error) {
	// For service-to-service calls, we pass through the token from context
	// In production, this should use a service account token
	return catalogapi.BearerAuth{Token: ""}, nil
}

// AttributeClientModule provides AttributeClient with its dependencies
func AttributeClientModule() fx.Option {
	return fx.Module("attribute-client",
		fx.Provide(
			fx.Private,
			httpclient.ProvideHTTPClient("catalog-service"),
		),
		fx.Provide(provideCatalogApiClient),
		fx.Provide(newAttributeClient),
	)
}

func provideCatalogApiClient(
	httpClient *http.Client,
	cfg httpclient.ClientConfig,
	tracerProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
) (*catalogapi.Client, error) {
	return catalogapi.NewClient(
		cfg.BaseURL,
		&serviceSecuritySource{},
		catalogapi.WithClient(httpClient),
		catalogapi.WithTracerProvider(tracerProvider),
		catalogapi.WithMeterProvider(meterProvider),
	)
}
