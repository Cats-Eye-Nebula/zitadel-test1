package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func TestCommandSide_AddOrgGenericOAuthIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GenericOAuthProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GenericOAuthProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewOAuthIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								nil,
								idp.Options{},
							)),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewOAuthIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								[]string{"user"},
								idp.Options{
									IsCreationAllowed: true,
									IsLinkingAllowed:  true,
									IsAutoCreation:    true,
									IsAutoUpdate:      true,
								},
							)),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
					Scopes:                []string{"user"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddOrgGenericOAuthProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGenericOAuthIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GenericOAuthProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GenericOAuthProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GenericOAuthProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewOAuthIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				want: &domain.ObjectDetails{},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewOAuthIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						eventPusherToEvents(
							func() eventstore.Command {
								t := true
								event, _ := org.NewOAuthIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
									"id1",
									"name",
									[]idp.OAuthIDPChanges{
										idp.ChangeOAuthName("new name"),
										idp.ChangeOAuthClientID("clientID2"),
										idp.ChangeOAuthClientSecret(&crypto.CryptoValue{
											CryptoType: crypto.TypeEncryption,
											Algorithm:  "enc",
											KeyID:      "id",
											Crypted:    []byte("newSecret"),
										}),
										idp.ChangeOAuthAuthorizationEndpoint("new auth"),
										idp.ChangeOAuthTokenEndpoint("new token"),
										idp.ChangeOAuthUserEndpoint("new user"),
										idp.ChangeOAuthScopes([]string{"openid", "profile"}),
										idp.ChangeOAuthOptions(idp.OptionChanges{
											IsCreationAllowed: &t,
											IsLinkingAllowed:  &t,
											IsAutoCreation:    &t,
											IsAutoUpdate:      &t,
										}),
									},
								)
								return event
							}(),
						),
						uniqueConstraintsFromEventConstraint(idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "org1")),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "org1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name:                  "new name",
					ClientID:              "clientID2",
					ClientSecret:          "newSecret",
					AuthorizationEndpoint: "new auth",
					TokenEndpoint:         "new token",
					UserEndpoint:          "new user",
					Scopes:                []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGenericOAuthProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgGenericOIDCIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GenericOIDCProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GenericOIDCProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOIDCProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOIDCProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOIDCProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewOIDCIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOIDCProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewOIDCIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"user"},
								idp.Options{
									IsCreationAllowed: true,
									IsLinkingAllowed:  true,
									IsAutoCreation:    true,
									IsAutoUpdate:      true,
								},
							)),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOIDCProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"user"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddOrgGenericOIDCProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGenericOIDCIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GenericOIDCProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GenericOIDCProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GenericOIDCProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOIDCProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOIDCProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOIDCProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewOIDCIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOIDCProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewOIDCIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						eventPusherToEvents(
							func() eventstore.Command {
								t := true
								event, _ := org.NewOIDCIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
									"id1",
									"name",
									[]idp.OIDCIDPChanges{
										idp.ChangeOIDCName("new name"),
										idp.ChangeOIDCIssuer("new issuer"),
										idp.ChangeOIDCClientID("clientID2"),
										idp.ChangeOIDCClientSecret(&crypto.CryptoValue{
											CryptoType: crypto.TypeEncryption,
											Algorithm:  "enc",
											KeyID:      "id",
											Crypted:    []byte("newSecret"),
										}),
										idp.ChangeOIDCScopes([]string{"openid", "profile"}),
										idp.ChangeOIDCOptions(idp.OptionChanges{
											IsCreationAllowed: &t,
											IsLinkingAllowed:  &t,
											IsAutoCreation:    &t,
											IsAutoUpdate:      &t,
										}),
									},
								)
								return event
							}(),
						),
						uniqueConstraintsFromEventConstraint(idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "org1")),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "org1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOIDCProvider{
					Name:         "new name",
					Issuer:       "new issuer",
					ClientID:     "clientID2",
					ClientSecret: "newSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGenericOIDCProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgJWTIDP(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      JWTProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      JWTProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: JWTProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid jwt endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: JWTProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid key endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid header name",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
					KeyEndpoint: "keys",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewJWTIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"jwt",
								"keys",
								"header",
								idp.Options{},
							)),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
					KeyEndpoint: "keys",
					HeaderName:  "header",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewJWTIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"jwt",
								"keys",
								"header",
								idp.Options{
									IsCreationAllowed: true,
									IsLinkingAllowed:  true,
									IsAutoCreation:    true,
									IsAutoUpdate:      true,
								},
							)),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
					KeyEndpoint: "keys",
					HeaderName:  "header",
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			id, got, err := c.AddOrgJWTProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgJWTIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      JWTProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      JWTProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      JWTProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: JWTProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid jwt endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: JWTProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid key endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid header name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
					KeyEndpoint: "keys",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
					KeyEndpoint: "keys",
					HeaderName:  "header",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewJWTIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"jwt",
								"keys",
								"header",
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
					KeyEndpoint: "keys",
					HeaderName:  "header",
				},
			},
			res: res{
				want: &domain.ObjectDetails{},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewJWTIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"jwt",
								"keys",
								"header",
								idp.Options{},
							)),
					),
					expectPush(
						eventPusherToEvents(
							func() eventstore.Command {
								t := true
								event, _ := org.NewJWTIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
									"id1",
									"name",
									[]idp.JWTIDPChanges{
										idp.ChangeJWTName("new name"),
										idp.ChangeJWTIssuer("new issuer"),
										idp.ChangeJWTEndpoint("new jwt"),
										idp.ChangeJWTKeysEndpoint("new keys"),
										idp.ChangeJWTHeaderName("new header"),
										idp.ChangeJWTOptions(idp.OptionChanges{
											IsCreationAllowed: &t,
											IsLinkingAllowed:  &t,
											IsAutoCreation:    &t,
											IsAutoUpdate:      &t,
										}),
									},
								)
								return event
							}(),
						),
						uniqueConstraintsFromEventConstraint(idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "org1")),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "org1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: JWTProvider{
					Name:        "new name",
					Issuer:      "new issuer",
					JWTEndpoint: "new jwt",
					KeyEndpoint: "new keys",
					HeaderName:  "new header",
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgJWTProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgAzureADIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      AzureADProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      AzureADProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid client id",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AzureADProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid client secret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewAzureADIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								"",
								false,
								idp.Options{},
							)),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AzureADProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewAzureADIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid"},
								"tenant",
								true,
								idp.Options{
									IsCreationAllowed: true,
									IsLinkingAllowed:  true,
									IsAutoCreation:    true,
									IsAutoUpdate:      true,
								},
							)),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AzureADProvider{
					Name:          "name",
					ClientID:      "clientID",
					ClientSecret:  "clientSecret",
					Scopes:        []string{"openid"},
					Tenant:        "tenant",
					EmailVerified: true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddOrgAzureADProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgAzureADIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      AzureADProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      AzureADProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      AzureADProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid client id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AzureADProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewAzureADIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								"",
								false,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewAzureADIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								"",
								false,
								idp.Options{},
							)),
					),
					expectPush(
						eventPusherToEvents(
							func() eventstore.Command {
								t := true
								event, _ := org.NewAzureADIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
									"id1",
									"name",
									[]idp.AzureADIDPChanges{
										idp.ChangeAzureADName("new name"),
										idp.ChangeAzureADClientID("new clientID"),
										idp.ChangeAzureADClientSecret(&crypto.CryptoValue{
											CryptoType: crypto.TypeEncryption,
											Algorithm:  "enc",
											KeyID:      "id",
											Crypted:    []byte("new clientSecret"),
										}),
										idp.ChangeAzureADScopes([]string{"openid", "profile"}),
										idp.ChangeAzureADTenant("new tenant"),
										idp.ChangeAzureADIsEmailVerified(true),
										idp.ChangeAzureADOptions(idp.OptionChanges{
											IsCreationAllowed: &t,
											IsLinkingAllowed:  &t,
											IsAutoCreation:    &t,
											IsAutoUpdate:      &t,
										}),
									},
								)
								return event
							}(),
						),
						uniqueConstraintsFromEventConstraint(idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "org1")),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "org1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AzureADProvider{
					Name:          "new name",
					ClientID:      "new clientID",
					ClientSecret:  "new clientSecret",
					Scopes:        []string{"openid", "profile"},
					Tenant:        "new tenant",
					EmailVerified: true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgAzureADProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgGitHubIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GitHubProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid client id",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitHubProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid client secret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewGitHubIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewGitHubIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid"},
								idp.Options{
									IsCreationAllowed: true,
									IsLinkingAllowed:  true,
									IsAutoCreation:    true,
									IsAutoUpdate:      true,
								},
							)),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddOrgGitHubProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGitHubIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GitHubProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitHubProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid client id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GitHubProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewGitHubIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewGitHubIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						eventPusherToEvents(
							func() eventstore.Command {
								t := true
								event, _ := org.NewGitHubIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
									"id1",
									[]idp.OAuthIDPChanges{
										idp.ChangeOAuthClientID("new clientID"),
										idp.ChangeOAuthClientSecret(&crypto.CryptoValue{
											CryptoType: crypto.TypeEncryption,
											Algorithm:  "enc",
											KeyID:      "id",
											Crypted:    []byte("new clientSecret"),
										}),
										idp.ChangeOAuthScopes([]string{"openid", "profile"}),
										idp.ChangeOAuthOptions(idp.OptionChanges{
											IsCreationAllowed: &t,
											IsLinkingAllowed:  &t,
											IsAutoCreation:    &t,
											IsAutoUpdate:      &t,
										}),
									},
								)
								return event
							}(),
						),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubProvider{
					ClientID:     "new clientID",
					ClientSecret: "new clientSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGitHubProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgGitHubEnterpriseIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GitHubEnterpriseProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitHubEnterpriseProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								nil,
								idp.Options{},
							)),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								[]string{"user"},
								idp.Options{
									IsCreationAllowed: true,
									IsLinkingAllowed:  true,
									IsAutoCreation:    true,
									IsAutoUpdate:      true,
								},
							)),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
					Scopes:                []string{"user"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddOrgGitHubEnterpriseProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGitHubEnterpriseIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GitHubEnterpriseProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitHubEnterpriseProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GitHubEnterpriseProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				want: &domain.ObjectDetails{},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						eventPusherToEvents(
							func() eventstore.Command {
								t := true
								event, _ := org.NewGitHubEnterpriseIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
									"id1",
									"name",
									[]idp.OAuthIDPChanges{
										idp.ChangeOAuthName("new name"),
										idp.ChangeOAuthClientID("clientID2"),
										idp.ChangeOAuthClientSecret(&crypto.CryptoValue{
											CryptoType: crypto.TypeEncryption,
											Algorithm:  "enc",
											KeyID:      "id",
											Crypted:    []byte("newSecret"),
										}),
										idp.ChangeOAuthAuthorizationEndpoint("new auth"),
										idp.ChangeOAuthTokenEndpoint("new token"),
										idp.ChangeOAuthUserEndpoint("new user"),
										idp.ChangeOAuthScopes([]string{"openid", "profile"}),
										idp.ChangeOAuthOptions(idp.OptionChanges{
											IsCreationAllowed: &t,
											IsLinkingAllowed:  &t,
											IsAutoCreation:    &t,
											IsAutoUpdate:      &t,
										}),
									},
								)
								return event
							}(),
						),
						uniqueConstraintsFromEventConstraint(idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "org1")),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "org1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "new name",
					ClientID:              "clientID2",
					ClientSecret:          "newSecret",
					AuthorizationEndpoint: "new auth",
					TokenEndpoint:         "new token",
					UserEndpoint:          "new user",
					Scopes:                []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGitHubEnterpriseProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgGitLabIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GitLabProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid clientID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitLabProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewGitLabIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewGitLabIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid"},
								idp.Options{
									IsCreationAllowed: true,
									IsLinkingAllowed:  true,
									IsAutoCreation:    true,
									IsAutoUpdate:      true,
								},
							)),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddOrgGitLabProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGitLabIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GitLabProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitLabProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GitLabProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewGitLabIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewGitLabIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						eventPusherToEvents(
							func() eventstore.Command {
								t := true
								event, _ := org.NewGitLabIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
									"id1",
									[]idp.GitLabIDPChanges{
										idp.ChangeGitLabClientID("clientID2"),
										idp.ChangeGitLabClientSecret(&crypto.CryptoValue{
											CryptoType: crypto.TypeEncryption,
											Algorithm:  "enc",
											KeyID:      "id",
											Crypted:    []byte("newSecret"),
										}),
										idp.ChangeGitLabScopes([]string{"openid", "profile"}),
										idp.ChangeGitLabOptions(idp.OptionChanges{
											IsCreationAllowed: &t,
											IsLinkingAllowed:  &t,
											IsAutoCreation:    &t,
											IsAutoUpdate:      &t,
										}),
									},
								)
								return event
							}(),
						),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabProvider{
					ClientID:     "clientID2",
					ClientSecret: "newSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGitLabProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgGitLabSelfHostedIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GitLabSelfHostedProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitLabSelfHostedProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabSelfHostedProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabSelfHostedProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabSelfHostedProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabSelfHostedProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid"},
								idp.Options{
									IsCreationAllowed: true,
									IsLinkingAllowed:  true,
									IsAutoCreation:    true,
									IsAutoUpdate:      true,
								},
							)),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabSelfHostedProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddOrgGitLabSelfHostedProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGitLabSelfHostedIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GitLabSelfHostedProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitLabSelfHostedProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GitLabSelfHostedProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabSelfHostedProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabSelfHostedProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabSelfHostedProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabSelfHostedProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						eventPusherToEvents(
							func() eventstore.Command {
								t := true
								event, _ := org.NewGitLabSelfHostedIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
									"id1",
									"name",
									[]idp.GitLabSelfHostedIDPChanges{
										idp.ChangeGitLabSelfHostedName("new name"),
										idp.ChangeGitLabSelfHostedIssuer("new issuer"),
										idp.ChangeGitLabSelfHostedClientID("clientID2"),
										idp.ChangeGitLabSelfHostedClientSecret(&crypto.CryptoValue{
											CryptoType: crypto.TypeEncryption,
											Algorithm:  "enc",
											KeyID:      "id",
											Crypted:    []byte("newSecret"),
										}),
										idp.ChangeGitLabSelfHostedScopes([]string{"openid", "profile"}),
										idp.ChangeGitLabSelfHostedOptions(idp.OptionChanges{
											IsCreationAllowed: &t,
											IsLinkingAllowed:  &t,
											IsAutoCreation:    &t,
											IsAutoUpdate:      &t,
										}),
									},
								)
								return event
							}(),
						),
						uniqueConstraintsFromEventConstraint(idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "org1")),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "org1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabSelfHostedProvider{
					Name:         "new name",
					Issuer:       "new issuer",
					ClientID:     "clientID2",
					ClientSecret: "newSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGitLabSelfHostedProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgGoogleIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GoogleProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid clientID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GoogleProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewGoogleIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewGoogleIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid"},
								idp.Options{
									IsCreationAllowed: true,
									IsLinkingAllowed:  true,
									IsAutoCreation:    true,
									IsAutoUpdate:      true,
								},
							)),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddOrgGoogleProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGoogleIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GoogleProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GoogleProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GoogleProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewGoogleIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewGoogleIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						eventPusherToEvents(
							func() eventstore.Command {
								t := true
								event, _ := org.NewGoogleIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
									"id1",
									[]idp.GoogleIDPChanges{
										idp.ChangeGoogleClientID("clientID2"),
										idp.ChangeGoogleClientSecret(&crypto.CryptoValue{
											CryptoType: crypto.TypeEncryption,
											Algorithm:  "enc",
											KeyID:      "id",
											Crypted:    []byte("newSecret"),
										}),
										idp.ChangeGoogleScopes([]string{"openid", "profile"}),
										idp.ChangeGoogleOptions(idp.OptionChanges{
											IsCreationAllowed: &t,
											IsLinkingAllowed:  &t,
											IsAutoCreation:    &t,
											IsAutoUpdate:      &t,
										}),
									},
								)
								return event
							}(),
						),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GoogleProvider{
					ClientID:     "clientID2",
					ClientSecret: "newSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGoogleProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}
