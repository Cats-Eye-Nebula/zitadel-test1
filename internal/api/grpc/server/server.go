package server

import (
	grpc_api "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/telemetry/metrics"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
)

type Server interface {
	Gateway
	RegisterServer(*grpc.Server)
	AppName() string
	MethodPrefix() string
	AuthMethods() authz.MethodMapping
}

func CreateServer(verifier *authz.TokenVerifier, authConfig authz.Config, queries *query.Queries, hostHeaderName string) *grpc.Server {
	metricTypes := []metrics.MetricType{metrics.MetricTypeTotalCount, metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode}
	return grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.DefaultTracingServer(),
				middleware.MetricsHandler(metricTypes, grpc_api.Probes...),
				middleware.SentryHandler(),
				middleware.NoCacheInterceptor(),
				middleware.ErrorHandler(),
				middleware.InstanceInterceptor(queries, hostHeaderName),
				middleware.AuthorizationInterceptor(verifier, authConfig),
				middleware.TranslationHandler(queries),
				middleware.ValidationHandler(),
				middleware.ServiceHandler(),
			),
		),
	)
}

func CreateServerWithoutAuth(queries *query.Queries) *grpc.Server {
	metricTypes := []metrics.MetricType{metrics.MetricTypeTotalCount, metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode}
	return grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.DefaultTracingServer(),
				middleware.MetricsHandler(metricTypes, grpc_api.Probes...),
				middleware.SentryHandler(),
				middleware.NoCacheInterceptor(),
				middleware.ErrorHandler(),
				middleware.TranslationHandler(queries),
				middleware.ValidationHandler(),
				middleware.ServiceHandler(),
			),
		),
	)
}
