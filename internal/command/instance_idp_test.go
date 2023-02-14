package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestCommandSide_AddInstanceGenericOAuthIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GenericOAuthProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
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
			id, got, err := c.AddInstanceGenericOAuthProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceGenericOAuthIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GenericOAuthProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{},
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GenericOAuthProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewOAuthIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "instance1")),
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "instance1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGenericOAuthProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceGenericOIDCIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GenericOIDCProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
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
			id, got, err := c.AddInstanceGenericOIDCProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceGenericOIDCIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GenericOIDCProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{},
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GenericOIDCProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewOIDCIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "instance1")),
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "instance1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGenericOIDCProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceJWTIDP(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx      context.Context
		provider JWTProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: JWTProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewJWTIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
									"id1",
									"name",
									"issuer",
									"jwt",
									"keys",
									"header",
									idp.Options{},
								)),
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewJWTIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			id, got, err := c.AddInstanceJWTProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceJWTIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider JWTProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: JWTProvider{},
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: JWTProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewJWTIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewJWTIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewJWTIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "instance1")),
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "instance1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceJWTProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceAzureADIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider AzureADProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewAzureADIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewAzureADIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
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
			id, got, err := c.AddInstanceAzureADProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceAzureADIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider AzureADProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{},
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: AzureADProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewAzureADIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewAzureADIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewAzureADIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "instance1")),
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "instance1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceAzureADProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceGitHubIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GitHubProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitHubIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitHubIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
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
			id, got, err := c.AddInstanceGitHubProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceGitHubIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GitHubProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubProvider{},
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GitHubProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewGitHubIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewGitHubIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewGitHubIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGitHubProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceGitHubEnterpriseIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GitHubEnterpriseProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
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
			id, got, err := c.AddInstanceGitHubEnterpriseProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceGitHubEnterpriseIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GitHubEnterpriseProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{},
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GitHubEnterpriseProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewGitHubEnterpriseIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "instance1")),
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "instance1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGitHubEnterpriseProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceGitLabIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GitLabProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitLabIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitLabIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
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
			id, got, err := c.AddInstanceGitLabProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceGitLabIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GitLabProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabProvider{},
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GitLabProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewGitLabIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewGitLabIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewGitLabIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGitLabProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceGitLabSelfHostedIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GitLabSelfHostedProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
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
			id, got, err := c.AddInstanceGitLabSelfHostedProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceGitLabSelfHostedIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GitLabSelfHostedProvider
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{},
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
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GitLabSelfHostedProvider{},
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
							instance.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewGitLabSelfHostedIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "instance1")),
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "instance1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGitLabSelfHostedProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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
