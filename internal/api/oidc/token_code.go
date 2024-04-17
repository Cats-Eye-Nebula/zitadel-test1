package oidc

import (
	"context"
	"slices"
	"strings"
	"time"

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
		err = oidcError(err)
		span.EndWithError(err)
	}()

	client, ok := r.Client.(*Client)
	if !ok {
		// not supposed to happen, but just preventing a panic if it does.
		return nil, zerrors.ThrowInternal(nil, "OIDC-eShi5", "Error.Internal")
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
		session, state, err = s.command.CreateOIDCSessionFromCodeExchange(
			setContextUserSystem(ctx), plainCode, authRequestComplianceChecker(client, r.Data),
		)
	} else {
		session, state, err = s.codeExchangeV1(ctx, client, r.Data, plainCode)
	}
	if err != nil {
		return nil, err
	}

	accessToken, idToken, err := s.createTokensFromSession(ctx, client, session)
	if err != nil {
		return nil, err
	}
	return op.NewResponse(&oidc.AccessTokenResponse{
		AccessToken:  accessToken,
		TokenType:    oidc.BearerToken,
		RefreshToken: session.RefreshToken,
		ExpiresIn:    timeToOIDCExpiresIn(session.Expiration),
		IDToken:      idToken,
		State:        state,
	}), nil
}

func (s *Server) createTokensFromSession(ctx context.Context, client *Client, session *command.OIDCSession) (accessToken, idToken string, err error) {
	getUserInfoAndSigner := s.getUserInfoAndSignerOnce(ctx, session.UserID, client.client.ProjectID, client.client.ProjectRoleAssertion, session.Scopes)

	if client.AccessTokenType() == op.AccessTokenTypeJWT {
		accessToken, err = s.createJWT(ctx, client, session, getUserInfoAndSigner)
	} else {
		accessToken, err = op.CreateBearerToken(session.TokenID, session.UserID, s.opCrypto)
	}
	if err != nil {
		return "", "", err
	}

	if slices.Contains(session.Scopes, oidc.ScopeOpenID) {
		idToken, _, err = s.createIDToken(ctx, client, getUserInfoAndSigner, accessToken, session.Audience, nil, time.Time{}, nil) // TODO: authmethods, authTime
	}

	return accessToken, idToken, err
}

// codeExchangeV1 creates a v2 token from a v1 auth request.
func (s *Server) codeExchangeV1(ctx context.Context, client *Client, req *oidc.AccessTokenRequest, plainCode string) (session *command.OIDCSession, state string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	authReq, err := s.getAuthRequestV1ByCode(ctx, plainCode)
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
	userAgentID, _, userOrgID, authTime, authMethodsReferences, reason, actor := getInfoFromRequest(authReq)

	session, err = s.command.CreateOIDCSession(ctx,
		authReq.GetSubject(),
		userOrgID,
		client.client.ClientID,
		authReq.GetAudience(),
		authReq.GetScopes(),
		AMRToAuthMethodTypes(authMethodsReferences),
		authTime,
		&domain.UserAgent{
			FingerprintID: &userAgentID,
		},
		reason,
		actor,
	)
	return session, authReq.GetState(), err
}

func (s *Server) getAuthRequestV1ByCode(ctx context.Context, plainCode string) (op.AuthRequest, error) {
	authReq, err := s.repo.AuthRequestByCode(ctx, plainCode)
	if err != nil {
		return nil, err
	}
	return AuthRequestFromBusiness(authReq)
}

func authRequestComplianceChecker(client *Client, req *oidc.AccessTokenRequest) command.AuthRequestComplianceChecker {
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
