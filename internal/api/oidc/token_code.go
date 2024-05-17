package oidc

import (
	"context"
	"slices"
	"strings"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (s *Server) CodeExchange(ctx context.Context, r *op.ClientRequest[oidc.AccessTokenRequest]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		span.EndWithError(err)
		err = oidcError(err)
	}()

	client, ok := r.Client.(*Client)
	if !ok {
		return nil, zerrors.ThrowInternal(nil, "OIDC-Ae2ph", "Error.Internal")
	}

	plainCode, err := s.decryptCode(ctx, r.Data.Code)
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "OIDC-ahLi2", "Errors.User.Code.Invalid")
	}

	var (
		session *command.OIDCSession
		state   string
	)
	if strings.HasPrefix(plainCode, command.IDPrefixV2) {
		session, state, err = s.command.CreateOIDCSessionFromAuthRequest(
			setContextUserSystem(ctx),
			plainCode,
			codeExchangeComplianceChecker(client, r.Data),
			slices.Contains(client.GrantTypes(), oidc.GrantTypeRefreshToken),
		)
	} else {
		session, state, err = s.codeExchangeV1(ctx, client, r.Data, r.Data.Code)
	}
	if err != nil {
		return nil, err
	}
	return response(s.accessTokenResponseFromSession(ctx, client, session, state, client.client.ProjectID, client.client.ProjectRoleAssertion))
}

// codeExchangeV1 creates a v2 token from a v1 auth request.
func (s *Server) codeExchangeV1(ctx context.Context, client *Client, req *oidc.AccessTokenRequest, code string) (session *command.OIDCSession, state string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	authReq, err := s.getAuthRequestV1ByCode(ctx, code)
	if err != nil {
		return nil, "", err
	}

	if challenge := authReq.GetCodeChallenge(); challenge != nil || client.AuthMethod() == oidc.AuthMethodNone {
		if err = op.AuthorizeCodeChallenge(req.CodeVerifier, challenge); err != nil {
			return nil, "", err
		}
	}
	if req.RedirectURI != authReq.GetRedirectURI() {
		return nil, "", oidc.ErrInvalidGrant().WithDescription("redirect_uri does not correspond")
	}
	userAgentID, _, userOrgID, authTime, authMethodsReferences, preferredLanguage, reason, actor := getInfoFromRequest(authReq)

	scope := authReq.GetScopes()
	session, err = s.command.CreateOIDCSession(ctx,
		authReq.GetSubject(),
		userOrgID,
		client.client.ClientID,
		scope,
		authReq.GetAudience(),
		AMRToAuthMethodTypes(authMethodsReferences),
		authTime,
		authReq.GetNonce(),
		preferredLanguage,
		&domain.UserAgent{
			FingerprintID: &userAgentID,
		},
		reason,
		actor,
		slices.Contains(scope, oidc.ScopeOfflineAccess),
	)
	if err != nil {
		return nil, "", err
	}
	return session, authReq.GetState(), s.repo.DeleteAuthRequest(ctx, authReq.GetID())
}

// getAuthRequestV1ByCode finds the v1 auth request by code.
// code needs to be the encrypted version of the ID,
// this is required by the underlying repo.
func (s *Server) getAuthRequestV1ByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	authReq, err := s.repo.AuthRequestByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(authReq)
}

func codeExchangeComplianceChecker(client *Client, req *oidc.AccessTokenRequest) command.AuthRequestComplianceChecker {
	return func(ctx context.Context, authReq *command.AuthRequestWriteModel) error {
		if authReq.CodeChallenge != nil || client.AuthMethod() == oidc.AuthMethodNone {
			err := op.AuthorizeCodeChallenge(req.CodeVerifier, CodeChallengeToOIDC(authReq.CodeChallenge))
			if err != nil {
				return err
			}
		}
		if req.RedirectURI != authReq.RedirectURI {
			return oidc.ErrInvalidGrant().WithDescription("redirect_uri does not correspond")
		}
		if err := authReq.CheckAuthenticated(); err != nil {
			return err
		}
		return nil
	}
}
