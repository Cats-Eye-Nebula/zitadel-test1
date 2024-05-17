//go:build integration

package oidc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/profile"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/domain"
)

func TestServer_JWTProfile(t *testing.T) {
	userID, keyData, err := Tester.CreateOIDCJWTProfileClient(CTX)
	require.NoError(t, err)

	type claims struct {
		resourceOwnerID            any
		resourceOwnerName          any
		resourceOwnerPrimaryDomain any
		orgDomain                  any
	}
	tests := []struct {
		name       string
		keyData    []byte
		scope      []string
		wantClaims claims
		wantErr    bool
	}{
		{
			name:    "success",
			keyData: keyData,
			scope:   []string{oidc.ScopeOpenID},
		},
		{
			name:    "org id and domain scope",
			keyData: keyData,
			scope: []string{
				oidc.ScopeOpenID,
				domain.OrgIDScope + Tester.Organisation.ID,
				domain.OrgDomainPrimaryScope + Tester.Organisation.Domain,
			},
			wantClaims: claims{
				resourceOwnerID:            Tester.Organisation.ID,
				resourceOwnerName:          Tester.Organisation.Name,
				resourceOwnerPrimaryDomain: Tester.Organisation.Domain,
				orgDomain:                  Tester.Organisation.Domain,
			},
		},
		{
			name:    "invalid org domain filtered",
			keyData: keyData,
			scope: []string{
				oidc.ScopeOpenID,
				domain.OrgDomainPrimaryScope + Tester.Organisation.Domain,
				domain.OrgDomainPrimaryScope + "foo"},
			wantClaims: claims{
				orgDomain: Tester.Organisation.Domain,
			},
		},
		{
			name:    "invalid org id filtered",
			keyData: keyData,
			scope: []string{oidc.ScopeOpenID,
				domain.OrgIDScope + Tester.Organisation.ID,
				domain.OrgIDScope + "foo",
			},
			wantClaims: claims{
				resourceOwnerID:            Tester.Organisation.ID,
				resourceOwnerName:          Tester.Organisation.Name,
				resourceOwnerPrimaryDomain: Tester.Organisation.Domain,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenSource, err := profile.NewJWTProfileTokenSourceFromKeyFileData(CTX, Tester.OIDCIssuer(), tt.keyData, tt.scope)
			require.NoError(t, err)

			tokens, err := tokenSource.TokenCtx(CTX)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, tokens)

			provider, err := rp.NewRelyingPartyOIDC(CTX, Tester.OIDCIssuer(), "", "", redirectURI, tt.scope)
			require.NoError(t, err)
			userinfo, err := rp.Userinfo[*oidc.UserInfo](CTX, tokens.AccessToken, oidc.BearerToken, userID, provider)
			require.NoError(t, err)
			assert.Equal(t, tt.wantClaims.resourceOwnerID, userinfo.Claims[oidc_api.ClaimResourceOwnerID])
			assert.Equal(t, tt.wantClaims.resourceOwnerName, userinfo.Claims[oidc_api.ClaimResourceOwnerName])
			assert.Equal(t, tt.wantClaims.resourceOwnerPrimaryDomain, userinfo.Claims[oidc_api.ClaimResourceOwnerPrimaryDomain])
			assert.Equal(t, tt.wantClaims.orgDomain, userinfo.Claims[domain.OrgDomainPrimaryClaim])
		})
	}
}
