package server

import (
	"context"

	grpc_api "github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"

	"github.com/caos/logging"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"golang.org/x/text/language"
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server/middleware"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

const (
	defaultGrpcPort = "80"
)

type Server interface {
	Gateway
	RegisterServer(*grpc.Server)
	AppName() string
	MethodPrefix() string
	AuthMethods() authz.MethodMapping
}

func CreateServer(verifier *authz.TokenVerifier, authConfig authz.Config, lang language.Tag) *grpc.Server {
	metricTypes := []metrics.MetricType{metrics.MetricTypeTotalCount, metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode}
	return grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.DefaultTracingServer(),
				middleware.MetricsHandler(metricTypes, grpc_api.Probes...),
				middleware.SentryHandler(),
				middleware.NoCacheInterceptor(),
				middleware.ErrorHandler(),
				middleware.AuthorizationInterceptor(verifier, authConfig),
				middleware.TranslationHandler(lang),
				middleware.ValidationHandler(),
				middleware.ServiceHandler(),
			),
		),
	)
}

func Serve(ctx context.Context, server *grpc.Server, port string) {
	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	go func() {
		listener := http.CreateListener(port)
		err := server.Serve(listener)
		logging.Log("SERVE-Ga3e94").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Panic("grpc server serve failed")
	}()
	logging.LogWithFields("SERVE-bZ44QM", "port", port).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Info("grpc server is listening")
}

func grpcPort(port string) string {
	if port == "" {
		return defaultGrpcPort
	}
	return port
}
