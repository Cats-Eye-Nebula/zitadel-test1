package oidc

import (
	"context"
	"time"

	"github.com/caos/logging"
	"github.com/caos/oidc/pkg/op"
	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"

	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/telemetry/metrics"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type OPHandlerConfig struct {
	OPConfig              *op.Config
	StorageConfig         StorageConfig
	UserAgentCookieConfig *middleware.UserAgentCookieConfig
	Cache                 *middleware.CacheConfig
	Endpoints             *EndpointConfig
}

type StorageConfig struct {
	DefaultLoginURL                   string
	SigningKeyAlgorithm               string
	DefaultAccessTokenLifetime        types.Duration
	DefaultIdTokenLifetime            types.Duration
	DefaultRefreshTokenIdleExpiration types.Duration
	DefaultRefreshTokenExpiration     types.Duration
}

type EndpointConfig struct {
	Auth          *Endpoint
	Token         *Endpoint
	Introspection *Endpoint
	Userinfo      *Endpoint
	Revocation    *Endpoint
	EndSession    *Endpoint
	Keys          *Endpoint
}

type Endpoint struct {
	Path string
	URL  string
}

type OPStorage struct {
	repo                              repository.Repository
	command                           *command.Commands
	query                             *query.Queries
	defaultLoginURL                   string
	defaultAccessTokenLifetime        time.Duration
	defaultIdTokenLifetime            time.Duration
	signingKeyAlgorithm               string
	defaultRefreshTokenIdleExpiration time.Duration
	defaultRefreshTokenExpiration     time.Duration
}

func NewProvider(ctx context.Context, config OPHandlerConfig, command *command.Commands, query *query.Queries, repo repository.Repository, keyConfig *crypto.KeyConfig, localDevMode bool) op.OpenIDProvider {
	cookieHandler, err := middleware.NewUserAgentHandler(config.UserAgentCookieConfig, id.SonyFlakeGenerator, localDevMode)
	logging.Log("OIDC-sd4fd").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Panic("cannot user agent handler")
	tokenKey, err := crypto.LoadKey(keyConfig, keyConfig.EncryptionKeyID)
	logging.Log("OIDC-ADvbv").OnError(err).Panic("cannot load OP crypto key")
	cryptoKey := []byte(tokenKey)
	if len(cryptoKey) != 32 {
		logging.Log("OIDC-Dsfds").Panic("OP crypto key must be exactly 32 bytes")
	}
	copy(config.OPConfig.CryptoKey[:], cryptoKey)
	config.OPConfig.CodeMethodS256 = true
	config.OPConfig.AuthMethodPost = true
	config.OPConfig.AuthMethodPrivateKeyJWT = true
	config.OPConfig.GrantTypeRefreshToken = true
	supportedLanguages, err := getSupportedLanguages()
	logging.Log("OIDC-GBd3t").OnError(err).Panic("cannot get supported languages")
	config.OPConfig.SupportedUILocales = supportedLanguages
	metricTypes := []metrics.MetricType{metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode, metrics.MetricTypeTotalCount}
	provider, err := op.NewOpenIDProvider(
		ctx,
		config.OPConfig,
		newStorage(config.StorageConfig, command, query, repo),
		op.WithHttpInterceptors(
			middleware.MetricsHandler(metricTypes),
			middleware.TelemetryHandler(),
			middleware.NoCacheInterceptor,
			cookieHandler,
			http_utils.CopyHeadersToContext,
		),
		op.WithCustomAuthEndpoint(op.NewEndpointWithURL(config.Endpoints.Auth.Path, config.Endpoints.Auth.URL)),
		op.WithCustomTokenEndpoint(op.NewEndpointWithURL(config.Endpoints.Token.Path, config.Endpoints.Token.URL)),
		op.WithCustomIntrospectionEndpoint(op.NewEndpointWithURL(config.Endpoints.Introspection.Path, config.Endpoints.Introspection.URL)),
		op.WithCustomUserinfoEndpoint(op.NewEndpointWithURL(config.Endpoints.Userinfo.Path, config.Endpoints.Userinfo.URL)),
		op.WithCustomRevocationEndpoint(op.NewEndpointWithURL(config.Endpoints.Revocation.Path, config.Endpoints.Revocation.URL)),
		op.WithCustomEndSessionEndpoint(op.NewEndpointWithURL(config.Endpoints.EndSession.Path, config.Endpoints.EndSession.URL)),
		op.WithCustomKeysEndpoint(op.NewEndpointWithURL(config.Endpoints.Keys.Path, config.Endpoints.Keys.URL)),
	)
	logging.Log("OIDC-asf13").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Panic("cannot create provider")
	return provider
}

func newStorage(config StorageConfig, command *command.Commands, query *query.Queries, repo repository.Repository) *OPStorage {
	return &OPStorage{
		repo:                              repo,
		command:                           command,
		query:                             query,
		defaultLoginURL:                   config.DefaultLoginURL,
		signingKeyAlgorithm:               config.SigningKeyAlgorithm,
		defaultAccessTokenLifetime:        config.DefaultAccessTokenLifetime.Duration,
		defaultIdTokenLifetime:            config.DefaultIdTokenLifetime.Duration,
		defaultRefreshTokenIdleExpiration: config.DefaultRefreshTokenIdleExpiration.Duration,
		defaultRefreshTokenExpiration:     config.DefaultRefreshTokenExpiration.Duration,
	}
}

func (o *OPStorage) Health(ctx context.Context) error {
	return o.repo.Health(ctx)
}

func getSupportedLanguages() ([]language.Tag, error) {
	statikLoginFS, err := fs.NewWithNamespace("login")
	if err != nil {
		return nil, err
	}
	return i18n.SupportedLanguages(statikLoginFS)
}
