package oidc

import (
	"context"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
)

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		name         string
		issuer       string
		clientID     string
		clientSecret string
		redirectURI  string
		httpMock     func(issuer string)
	}
	tests := []struct {
		name   string
		fields fields
		want   idp.Session
	}{
		{
			name: "successful auth",
			fields: fields{
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
			},
			want: &Session{AuthURL: "https://issuer.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			tt.fields.httpMock(tt.fields.issuer)
			a := assert.New(t)

			provider, err := New(tt.fields.name, tt.fields.issuer, tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI)
			a.NoError(err)

			session, err := provider.BeginAuth(context.Background(), "testState")
			a.NoError(err)

			a.Equal(tt.want.GetAuthURL(), session.GetAuthURL())
		})
	}
}

func TestProvider_Options(t *testing.T) {
	type fields struct {
		name         string
		issuer       string
		clientID     string
		clientSecret string
		redirectURI  string
		httpMock     func(issuer string)
		opts         []ProviderOpts
	}
	type want struct {
		name            string
		linkingAllowed  bool
		creationAllowed bool
		autoCreation    bool
		autoUpdate      bool
		pkce            bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "default",
			fields: fields{
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				opts:         nil,
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
			},
			want: want{
				name:            "oidc",
				linkingAllowed:  false,
				creationAllowed: false,
				autoCreation:    false,
				autoUpdate:      false,
				pkce:            false,
			},
		},
		{
			name: "all true",
			fields: fields{
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				opts: []ProviderOpts{
					WithLinkingAllowed(),
					WithCreationAllowed(),
					WithAutoCreation(),
					WithAutoUpdate(),
					WithRelyingPartyOption(rp.WithPKCE(nil)),
				},
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
			},
			want: want{
				name:            "oidc",
				linkingAllowed:  true,
				creationAllowed: true,
				autoCreation:    true,
				autoUpdate:      true,
				pkce:            true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			gock.EnableNetworking()
			tt.fields.httpMock(tt.fields.issuer)
			a := assert.New(t)

			provider, err := New(tt.fields.name, tt.fields.issuer, tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.opts...)
			a.NoError(err)

			a.Equal(tt.want.name, provider.Name())
			a.Equal(tt.want.linkingAllowed, provider.IsLinkingAllowed())
			a.Equal(tt.want.creationAllowed, provider.IsCreationAllowed())
			a.Equal(tt.want.autoCreation, provider.IsAutoCreation())
			a.Equal(tt.want.autoUpdate, provider.IsAutoUpdate())
		})
	}
}
