package client

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	catalogapi "github.com/Sokol111/ecommerce-catalog-service-api/gen/httpapi"
	httpclient "github.com/Sokol111/ecommerce-commons/pkg/http/client"
	"github.com/Sokol111/ecommerce-commons/pkg/security/token"
)

// serviceSecuritySource provides service-to-service authentication.
type serviceSecuritySource struct {
	serviceToken string
}

func newServiceSecuritySource(cfg token.Config) *serviceSecuritySource {
	return &serviceSecuritySource{serviceToken: cfg.ServiceToken}
}

func (s *serviceSecuritySource) BearerAuth(_ context.Context, _ catalogapi.OperationName) (catalogapi.BearerAuth, error) {
	return catalogapi.BearerAuth{Token: s.serviceToken}, nil
}

// AttributeClientModule provides AttributeClient with its dependencies
func AttributeClientModule() fx.Option {
	return fx.Module("attribute-client",
		fx.Provide(
			fx.Private,
			httpclient.ProvideHTTPClient("catalog-service"),
			newServiceSecuritySource,
		),
		fx.Provide(provideCatalogApiClient),
		fx.Provide(newAttributeClient),
	)
}

func provideCatalogApiClient(
	httpClient *http.Client,
	cfg httpclient.ClientConfig,
	securitySource *serviceSecuritySource,
	tracerProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
) (*catalogapi.Client, error) {
	return catalogapi.NewClient(
		cfg.BaseURL,
		securitySource,
		catalogapi.WithClient(httpClient),
		catalogapi.WithTracerProvider(tracerProvider),
		catalogapi.WithMeterProvider(meterProvider),
	)
}
