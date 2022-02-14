package login

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"

	"github.com/caos/zitadel/internal/api/authz"
	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/api/http/middleware"
	_ "github.com/caos/zitadel/internal/api/ui/login/statik"
	auth_repository "github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/form"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/static"
)

type Login struct {
	endpoint      string
	router        http.Handler
	renderer      *Renderer
	parser        *form.Parser
	command       *command.Commands
	query         *query.Queries
	staticStorage static.Storage
	//staticCache         cache.Cache //TODO: enable when storage is implemented again
	authRepo            auth_repository.Repository
	baseURL             string
	zitadelURL          string
	oidcAuthCallbackURL string
	IDPConfigAesCrypto  crypto.EncryptionAlgorithm
	iamDomain           string
}

type Config struct {
	LanguageCookieName string
	CSRF               CSRF
	Cache              middleware.CacheConfig
	//StaticCache         cache_config.CacheConfig //TODO: enable when storage is implemented again
}

type CSRF struct {
	CookieName string
	Key        *crypto.KeyConfig
}

const (
	login                = "LOGIN"
	HandlerPrefix        = "/ui/login"
	DefaultLoggedOutPath = HandlerPrefix + EndpointLogoutDone
)

func CreateLogin(config Config, command *command.Commands, query *query.Queries, authRepo *eventsourcing.EsRepository, staticStorage static.Storage, systemDefaults systemdefaults.SystemDefaults, zitadelURL, domain, oidcAuthCallbackURL string, localDevMode bool, userAgentCookie mux.MiddlewareFunc) (*Login, error) {
	aesCrypto, err := crypto.NewAESCrypto(systemDefaults.IDPConfigVerificationKey)
	if err != nil {
		return nil, fmt.Errorf("error create new aes crypto: %w", err)
	}
	login := &Login{
		oidcAuthCallbackURL: oidcAuthCallbackURL,
		baseURL:             HandlerPrefix,
		zitadelURL:          zitadelURL,
		command:             command,
		query:               query,
		staticStorage:       staticStorage,
		authRepo:            authRepo,
		IDPConfigAesCrypto:  aesCrypto,
		iamDomain:           domain,
	}
	//TODO: enable when storage is implemented again
	//login.staticCache, err = config.StaticCache.Config.NewCache()
	//if err != nil {
	//	return nil, fmt.Errorf("unable to create storage cache: %w", err)
	//}

	statikFS, err := fs.NewWithNamespace("login")
	if err != nil {
		return nil, fmt.Errorf("unable to create filesystem: %w", err)
	}

	csrfInterceptor, err := createCSRFInterceptor(config.CSRF, localDevMode, login.csrfErrorHandler())
	if err != nil {
		return nil, fmt.Errorf("unable to create csrfInterceptor: %w", err)
	}
	cacheInterceptor, err := middleware.DefaultCacheInterceptor(EndpointResources, config.Cache.MaxAge, config.Cache.SharedMaxAge)
	if err != nil {
		return nil, fmt.Errorf("unable to create cacheInterceptor: %w", err)
	}
	security := middleware.SecurityHeaders(csp(), login.cspErrorHandler)
	login.router = CreateRouter(login, statikFS, csrfInterceptor, cacheInterceptor, security, userAgentCookie, middleware.TelemetryHandler(EndpointResources))
	login.renderer = CreateRenderer(HandlerPrefix, statikFS, staticStorage, config.LanguageCookieName, systemDefaults.DefaultLanguage)
	login.parser = form.NewParser()
	return login, nil
}

func csp() *middleware.CSP {
	csp := middleware.DefaultSCP
	csp.ObjectSrc = middleware.CSPSourceOptsSelf()
	csp.StyleSrc = csp.StyleSrc.AddNonce()
	csp.ScriptSrc = csp.ScriptSrc.AddNonce()
	return &csp
}

func createCSRFInterceptor(config CSRF, localDevMode bool, errorHandler http.Handler) (func(http.Handler) http.Handler, error) {
	csrfKey, err := crypto.LoadKey(config.Key, config.Key.EncryptionKeyID)
	if err != nil {
		return nil, err
	}
	path := "/"
	return csrf.Protect([]byte(csrfKey),
		csrf.Secure(!localDevMode),
		csrf.CookieName(http_utils.SetCookiePrefix(config.CookieName, "", path, !localDevMode)),
		csrf.Path(path),
		csrf.ErrorHandler(errorHandler),
	), nil
}

func (l *Login) Handler() http.Handler {
	return l.router
}

func (l *Login) getClaimedUserIDsOfOrgDomain(ctx context.Context, orgName string) ([]string, error) {
	loginName, err := query.NewUserPreferredLoginNameSearchQuery("@"+domain.NewIAMDomainName(orgName, l.iamDomain), query.TextEndsWithIgnoreCase)
	if err != nil {
		return nil, err
	}
	users, err := l.query.SearchUsers(ctx, &query.UserSearchQueries{Queries: []query.SearchQuery{loginName}})
	if err != nil {
		return nil, err
	}
	userIDs := make([]string, len(users.Users))
	for i, user := range users.Users {
		userIDs[i] = user.ID
	}
	return userIDs, nil
}

func setContext(ctx context.Context, resourceOwner string) context.Context {
	data := authz.CtxData{
		UserID: login,
		OrgID:  resourceOwner,
	}
	return authz.SetCtxData(ctx, data)
}
