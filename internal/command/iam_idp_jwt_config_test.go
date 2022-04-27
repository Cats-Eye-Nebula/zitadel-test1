package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
)

func TestCommandSide_ChangeDefaultIDPJWTConfig(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type (
		args struct {
			ctx    context.Context
			config *domain.JWTIDPConfig
		}
	)
	type res struct {
		want *domain.JWTIDPConfig
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid config, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				config: &domain.JWTIDPConfig{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "idp config not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.JWTIDPConfig{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "idp config removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewIDPConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeJWT,
								domain.IDPConfigStylingTypeGoogle,
								false,
							),
						),
						eventFromEventPusher(
							iam.NewIDPJWTConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"jwt-endpoint",
								"issuer",
								"keys-endpoint",
								"auth",
							),
						),
						eventFromEventPusher(
							iam.NewIDPConfigRemovedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"name",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.JWTIDPConfig{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewIDPConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeJWT,
								domain.IDPConfigStylingTypeGoogle,
								false,
							),
						),
						eventFromEventPusher(
							iam.NewIDPJWTConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"jwt-endpoint",
								"issuer",
								"keys-endpoint",
								"auth",
							),
						),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.JWTIDPConfig{
					IDPConfigID:  "config1",
					JWTEndpoint:  "jwt-endpoint",
					Issuer:       "issuer",
					KeysEndpoint: "keys-endpoint",
					HeaderName:   "auth",
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp config jwt add, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewIDPConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeJWT,
								domain.IDPConfigStylingTypeGoogle,
								false,
							),
						),
						eventFromEventPusher(
							iam.NewIDPJWTConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"jwt-endpoint",
								"issuer",
								"keys-endpoint",
								"auth",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultIDPJWTConfigChangedEvent(context.Background(),
									"config1",
									"jwt-endpoint-changed",
									"issuer-changed",
									"keys-endpoint-changed",
									"auth-changed",
								),
							),
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.JWTIDPConfig{
					IDPConfigID:  "config1",
					JWTEndpoint:  "jwt-endpoint-changed",
					Issuer:       "issuer-changed",
					KeysEndpoint: "keys-endpoint-changed",
					HeaderName:   "auth-changed",
				},
			},
			res: res{
				want: &domain.JWTIDPConfig{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "IAM",
						ResourceOwner: "IAM",
					},
					IDPConfigID:  "config1",
					JWTEndpoint:  "jwt-endpoint-changed",
					Issuer:       "issuer-changed",
					KeysEndpoint: "keys-endpoint-changed",
					HeaderName:   "auth-changed",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:            tt.fields.eventstore,
				idpConfigSecretCrypto: tt.fields.secretCrypto,
			}
			got, err := r.ChangeDefaultIDPJWTConfig(tt.args.ctx, tt.args.config)
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

func newDefaultIDPJWTConfigChangedEvent(ctx context.Context, configID, jwtEndpoint, issuer, keysEndpoint, headerName string) *iam.IDPJWTConfigChangedEvent {
	event, _ := iam.NewIDPJWTConfigChangedEvent(ctx,
		&iam.NewAggregate().Aggregate,
		configID,
		[]idpconfig.JWTConfigChanges{
			idpconfig.ChangeJWTEndpoint(jwtEndpoint),
			idpconfig.ChangeJWTIssuer(issuer),
			idpconfig.ChangeKeysEndpoint(keysEndpoint),
			idpconfig.ChangeHeaderName(headerName),
		},
	)
	return event
}
