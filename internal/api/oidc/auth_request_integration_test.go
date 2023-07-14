//go:build integration

package oidc_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/command"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2alpha"
)

var (
	armPasskey  = []string{oidc_api.UserPresence, oidc_api.MFA}
	armPassword = []string{oidc_api.PWD}
)

func TestOPStorage_CreateAuthRequest(t *testing.T) {
	clientID := createClient(t)

	id := createAuthRequest(t, clientID, redirectURI)
	require.Contains(t, id, command.IDPrefixV2)
}

func TestOPStorage_CreateAccessToken_code(t *testing.T) {
	clientID := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI)
	sessionID, sessionToken, startTime, changeTime := Tester.CreatePasskeySession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.NoError(t, err)

	// test code exchange
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err := exchangeTokens(t, clientID, code)
	require.NoError(t, err)
	assertTokens(t, tokens, false)
	assertIDTokenClaims(t, tokens.IDTokenClaims, armPasskey, startTime, changeTime)

	// callback on a succeeded request must fail
	linkResp, err = Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.Error(t, err)

	// exchange with a used code must fail
	_, err = exchangeTokens(t, clientID, code)
	require.Error(t, err)
}

func TestOPStorage_CreateAccessToken_implicit(t *testing.T) {
	clientID := createImplicitClient(t)
	authRequestID := createAuthRequestImplicit(t, clientID, redirectURIImplicit)
	sessionID, sessionToken, startTime, changeTime := Tester.CreatePasskeySession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.NoError(t, err)

	// test implicit callback
	callback, err := url.Parse(linkResp.GetCallbackUrl())
	require.NoError(t, err)
	values, err := url.ParseQuery(callback.Fragment)
	require.NoError(t, err)
	accessToken := values.Get("access_token")
	idToken := values.Get("id_token")
	refreshToken := values.Get("refresh_token")
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, idToken)
	assert.Empty(t, refreshToken)
	assert.NotEmpty(t, values.Get("expires_in"))
	assert.Equal(t, oidc.BearerToken, values.Get("token_type"))
	assert.Equal(t, "state", values.Get("state"))

	// check id_token / claims
	provider, err := Tester.CreateRelyingParty(clientID, redirectURIImplicit)
	require.NoError(t, err)
	claims, err := rp.VerifyTokens[*oidc.IDTokenClaims](context.Background(), accessToken, idToken, provider.IDTokenVerifier())
	require.NoError(t, err)
	assertIDTokenClaims(t, claims, armPasskey, startTime, changeTime)

	// callback on a succeeded request must fail
	linkResp, err = Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.Error(t, err)
}

func TestOPStorage_CreateAccessAndRefreshTokens_code(t *testing.T) {
	clientID := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess)
	sessionID, sessionToken, startTime, changeTime := Tester.CreatePasskeySession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.NoError(t, err)

	// test code exchange (expect refresh token to be returned)
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err := exchangeTokens(t, clientID, code)
	require.NoError(t, err)
	assertTokens(t, tokens, true)
	assertIDTokenClaims(t, tokens.IDTokenClaims, armPasskey, startTime, changeTime)
}

func TestOPStorage_CreateAccessAndRefreshTokens_refresh(t *testing.T) {
	clientID := createClient(t)
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess)
	sessionID, sessionToken, startTime, changeTime := Tester.CreatePasskeySession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.NoError(t, err)

	// code exchange
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err := exchangeTokens(t, clientID, code)
	require.NoError(t, err)
	assertTokens(t, tokens, true)
	assertIDTokenClaims(t, tokens.IDTokenClaims, armPasskey, startTime, changeTime)

	// test actual refresh grant
	newTokens, err := refreshTokens(t, clientID, tokens.RefreshToken)
	require.NoError(t, err)
	idToken, _ := newTokens.Extra("id_token").(string)
	assert.NotEmpty(t, idToken)
	assert.NotEmpty(t, newTokens.AccessToken)
	assert.NotEmpty(t, newTokens.RefreshToken)
	claims, err := rp.VerifyTokens[*oidc.IDTokenClaims](context.Background(), newTokens.AccessToken, idToken, provider.IDTokenVerifier())
	require.NoError(t, err)
	// auth time must still be the initial
	assertIDTokenClaims(t, claims, armPasskey, startTime, changeTime)

	// refresh with an old refresh_token must fail
	_, err = rp.RefreshAccessToken(provider, tokens.RefreshToken, "", "")
	require.Error(t, err)
}

func TestOPStorage_GetRefreshTokenInfo(t *testing.T) {

}

func TestOPStorage_RevokeToken_access_token(t *testing.T) {
	clientID := createClient(t)
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess)
	sessionID, sessionToken, startTime, changeTime := Tester.CreatePasskeySession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.NoError(t, err)

	// code exchange
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err := exchangeTokens(t, clientID, code)
	require.NoError(t, err)
	assertTokens(t, tokens, true)
	assertIDTokenClaims(t, tokens.IDTokenClaims, armPasskey, startTime, changeTime)

	err = rp.RevokeToken(provider, tokens.AccessToken, "")
	require.Error(t, err)

	// userinfo must fail
	_, err = rp.Userinfo(tokens.AccessToken, tokens.TokenType, tokens.IDTokenClaims.Subject, provider)
	require.Error(t, err)

	// refresh grant must still work
	_, err = refreshTokens(t, clientID, tokens.RefreshToken)
	require.NoError(t, err)
}

func TestOPStorage_RevokeToken_refresh_access(t *testing.T) {
	clientID := createClient(t)
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess)
	sessionID, sessionToken, startTime, changeTime := Tester.CreatePasskeySession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.NoError(t, err)

	// code exchange
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err := exchangeTokens(t, clientID, code)
	require.NoError(t, err)
	assertTokens(t, tokens, true)
	assertIDTokenClaims(t, tokens.IDTokenClaims, armPasskey, startTime, changeTime)

	// revoke refresh token -> invalidates also access token
	err = rp.RevokeToken(provider, tokens.RefreshToken, "")
	require.Error(t, err)

	// userinfo must fail
	_, err = rp.Userinfo(tokens.AccessToken, tokens.TokenType, tokens.IDTokenClaims.Subject, provider)
	require.Error(t, err)

	// refresh must fail
	_, err = refreshTokens(t, clientID, tokens.RefreshToken)
	require.Error(t, err)
}

func TestOPStorage_TerminateSession(t *testing.T) {

}

func exchangeTokens(t testing.TB, clientID, code string) (*oidc.Tokens[*oidc.IDTokenClaims], error) {
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)

	codeVerifier := "codeVerifier"
	return rp.CodeExchange[*oidc.IDTokenClaims](context.Background(), code, provider, rp.WithCodeVerifier(codeVerifier))
}

func refreshTokens(t testing.TB, clientID, refreshToken string) (*oauth2.Token, error) {
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)

	return rp.RefreshAccessToken(provider, refreshToken, "", "")
}

func assertCodeResponse(t *testing.T, callback string) string {
	callbackURL, err := url.Parse(callback)
	require.NoError(t, err)
	code := callbackURL.Query().Get("code")
	require.NotEmpty(t, code)
	assert.Equal(t, "state", callbackURL.Query().Get("state"))
	return code
}

func assertTokens(t *testing.T, tokens *oidc.Tokens[*oidc.IDTokenClaims], requireRefreshToken bool) {
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.IDToken)
	if requireRefreshToken {
		assert.NotEmpty(t, tokens.RefreshToken)
	} else {
		assert.Empty(t, tokens.RefreshToken)
	}
}

func assertIDTokenClaims(t *testing.T, claims *oidc.IDTokenClaims, arm []string, sessionStart, sessionChange time.Time) {
	assert.Equal(t, User.GetUserId(), claims.Subject)
	assert.Equal(t, arm, claims.AuthenticationMethodsReferences)
	assertOIDCTimeRange(t, claims.AuthTime, sessionStart, sessionChange)
}
