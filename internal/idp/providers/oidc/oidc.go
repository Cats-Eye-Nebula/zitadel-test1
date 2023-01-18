package oidc

import (
	"context"

	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Provider = (*Provider)(nil)

// Provider is the [idp.Provider] implementation for a generic OIDC provider
type Provider struct {
	rp.RelyingParty
	options           []rp.Option
	name              string
	isLinkingAllowed  bool
	isCreationAllowed bool
	isAutoCreation    bool
	isAutoUpdate      bool
}

type ProviderOpts func(provider *Provider)

// WithLinkingAllowed allows end users to link the federated user to an existing one
func WithLinkingAllowed() ProviderOpts {
	return func(p *Provider) {
		p.isLinkingAllowed = true
	}
}

// WithCreationAllowed allows end users to create a new user using the federated information
func WithCreationAllowed() ProviderOpts {
	return func(p *Provider) {
		p.isCreationAllowed = true
	}
}

// WithAutoCreation enables that federated users are automatically created if not already existing
func WithAutoCreation() ProviderOpts {
	return func(p *Provider) {
		p.isAutoCreation = true
	}
}

// WithAutoUpdate enables that information retrieved from the provider is automatically used to update
// the existing user on each authentication
func WithAutoUpdate() ProviderOpts {
	return func(p *Provider) {
		p.isAutoUpdate = true
	}
}

// WithRelyingPartyOption allows to set an additional [rp.Option] like [rp.WithPKCE]
func WithRelyingPartyOption(option rp.Option) ProviderOpts {
	return func(p *Provider) {
		p.options = append(p.options, option)
	}
}

// New creates a generic OIDC provider
func New(name, issuer, clientID, clientSecret, redirectURI string, options ...ProviderOpts) (*Provider, error) {
	provider := &Provider{
		name: name,
	}
	for _, option := range options {
		option(provider)
	}
	relyingParty, err := rp.NewRelyingPartyOIDC(issuer, clientID, clientSecret, redirectURI, []string{oidc.ScopeOpenID}, provider.options...)
	if err != nil {
		return nil, err
	}
	provider.RelyingParty = relyingParty
	return provider, nil
}

// Name implements the [idp.Provider] interface
func (p *Provider) Name() string {
	return p.name
}

// BeginAuth implements the [idp.Provider] interface
// it will create a [Session] with an OIDC authorization request as AuthURL
func (p *Provider) BeginAuth(ctx context.Context, state string, _ ...any) (idp.Session, error) {
	url := rp.AuthURL(state, p.RelyingParty)
	return &Session{AuthURL: url, Provider: p}, nil
}

// IsLinkingAllowed implements the [idp.Provider] interface
func (p *Provider) IsLinkingAllowed() bool {
	return p.isLinkingAllowed
}

// IsCreationAllowed implements the [idp.Provider] interface
func (p *Provider) IsCreationAllowed() bool {
	return p.isCreationAllowed
}

// IsAutoCreation implements the [idp.Provider] interface
func (p *Provider) IsAutoCreation() bool {
	return p.isAutoCreation
}

// IsAutoUpdate implements the [idp.Provider] interface
func (p *Provider) IsAutoUpdate() bool {
	return p.isAutoUpdate
}
