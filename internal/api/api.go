package api

import (
	"context"
	"net/http"

	"github.com/caos/logging"
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/api/authz"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server"
	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/api/oidc"
	authz_es "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type Config struct {
	GRPC grpc_util.Config
	OIDC oidc.OPHandlerConfig
}

type API struct {
	grpcServer     *grpc.Server
	gatewayHandler *server.GatewayHandler
	verifier       *authz.TokenVerifier
	serverPort     string
	health         health
}
type health interface {
	Health(ctx context.Context) error
	IamByID(ctx context.Context) (*iam_model.IAM, error)
	VerifierClientID(ctx context.Context, appName string) (string, error)
}

func Create(config Config, authZ authz.Config, authZRepo *authz_es.EsRepository, sd systemdefaults.SystemDefaults) *API {
	api := &API{
		serverPort: config.GRPC.ServerPort,
	}
	api.verifier = authz.Start(authZRepo)
	api.health = authZRepo
	api.grpcServer = server.CreateServer(api.verifier, authZ, sd.DefaultLanguage)
	api.gatewayHandler = server.CreateGatewayHandler(config.GRPC)
	api.RegisterHandler("", api.healthHandler())

	return api
}

func (a *API) RegisterServer(ctx context.Context, server server.Server) {
	server.RegisterServer(a.grpcServer)
	a.gatewayHandler.RegisterGateway(ctx, server)
	a.verifier.RegisterServer(server.AppName(), server.MethodPrefix(), server.AuthMethods())
}

func (a *API) RegisterHandler(prefix string, handler http.Handler) {
	a.gatewayHandler.RegisterHandler(prefix, handler)
}

func (a *API) Start(ctx context.Context) {
	server.Serve(ctx, a.grpcServer, a.serverPort)
	a.gatewayHandler.Serve(ctx)
}

func (a *API) healthHandler() http.Handler {
	checks := []ValidationFunction{
		func(ctx context.Context) error {
			if err := a.health.Health(ctx); err != nil {
				return errors.ThrowInternal(err, "API-F24h2", "DB CONNECTION ERROR")
			}
			return nil
		},
		func(ctx context.Context) error {
			iam, err := a.health.IamByID(ctx)
			if err != nil && !errors.IsNotFound(err) {
				return errors.ThrowPreconditionFailed(err, "API-dsgT2", "IAM SETUP CHECK FAILED")
			}
			if iam == nil || iam.SetUpStarted < iam_model.StepCount-1 {
				return errors.ThrowPreconditionFailed(nil, "API-HBfs3", "IAM NOT SET UP")
			}
			if iam.SetUpDone < iam_model.StepCount-1 {
				return errors.ThrowPreconditionFailed(nil, "API-DASs2", "IAM SETUP RUNNING")
			}
			return nil
		},
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/healthz", handleHealth)
	handler.HandleFunc("/ready", handleReadiness(checks))
	handler.HandleFunc("/clientID", a.handleClientID)
	handler.HandleFunc("/projectID", a.handleProjectID)

	return handler
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	logging.Log("API-Hfss2").OnError(err).Error("error writing ok for health")
}

func handleReadiness(checks []ValidationFunction) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := validate(r.Context(), checks)
		if err == nil {
			http_util.MarshalJSON(w, "ok", nil, http.StatusOK)
			return
		}
		http_util.MarshalJSON(w, nil, err, http.StatusPreconditionFailed)
	}
}

func (a *API) handleClientID(w http.ResponseWriter, r *http.Request) {
	id, err := a.health.VerifierClientID(r.Context(), "Zitadel Console")
	if err != nil {
		http_util.MarshalJSON(w, nil, err, http.StatusPreconditionFailed)
		return
	}
	http_util.MarshalJSON(w, id, nil, http.StatusOK)
}

func (a *API) handleProjectID(w http.ResponseWriter, r *http.Request) {
	iam, err := a.health.IamByID(r.Context())
	if err != nil {
		http_util.MarshalJSON(w, nil, err, http.StatusPreconditionFailed)
		return
	}
	http_util.MarshalJSON(w, iam.IAMProjectID, nil, http.StatusOK)
}

type ValidationFunction func(ctx context.Context) error

func validate(ctx context.Context, validations []ValidationFunction) error {
	for _, validation := range validations {
		if err := validation(ctx); err != nil {
			logging.Log("API-vf823").WithError(err).Error("validation failed")
			return err
		}
	}
	return nil
}
