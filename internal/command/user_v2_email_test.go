package command

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestCommands_ChangeUserEmail(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		userID        string
		resourceOwner string
		email         string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Email
		wantErr error
	}{
		{
			name: "missing email",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyEmailCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "",
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty"),
		},
		{
			name: "not changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyEmailCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email@test.ch",
			},
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Uch5e", "Errors.User.Email.NotChanged"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := c.ChangeUserEmail(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.email, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestCommands_ChangeUserEmailURLTemplate(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		userID        string
		resourceOwner string
		email         string
		urlTmpl       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Email
		wantErr error
	}{
		{
			name: "invalid template",
			fields: fields{
				eventstore: eventstoreExpect(t),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email-changed@test.ch",
				urlTmpl:       "{{",
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "USERv2-ooD8p", "Errors.User.V2.Email.InvalidURLTemplate"),
		},
		{
			name: "not changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyEmailCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email@test.ch",
			},
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Uch5e", "Errors.User.Email.NotChanged"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := c.ChangeUserEmailURLTemplate(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.email, crypto.CreateMockEncryptionAlg(gomock.NewController(t)), tt.args.urlTmpl)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestCommands_ChangeUserEmailReturnCode(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		userID        string
		resourceOwner string
		email         string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Email
		wantErr error
	}{
		{
			name: "missing email",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("inst1").Aggregate,
								domain.SecretGeneratorTypeVerifyEmailCode,
								12, time.Minute, true, true, true, true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "",
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := c.ChangeUserEmailReturnCode(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.email, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestCommands_ChangeUserEmailVerified(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		userID        string
		resourceOwner string
		email         string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Email
		wantErr error
	}{
		{
			name: "missing email",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "",
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty"),
		},
		{
			name: "email changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanEmailChangedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"email-changed@test.ch",
								),
							),
							eventFromEventPusher(
								user.NewHumanEmailVerifiedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
								),
							),
						},
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email-changed@test.ch",
			},
			want: &domain.Email{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				EmailAddress:    "email-changed@test.ch",
				IsEmailVerified: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := c.ChangeUserEmailVerified(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.email)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestCommands_changeUserEmail(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		userID        string
		resourceOwner string
		email         string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Email
		wantErr error
	}{
		{
			name: "missing email",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "",
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty"),
		},
		{
			name: "not changed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email@test.ch",
			},
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Uch5e", "Errors.User.Email.NotChanged"),
		},
		{
			name: "email changed, with code",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanEmailChangedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"email-changed@test.ch",
								),
							),
							eventFromEventPusher(
								user.NewHumanEmailCodeAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("a"),
									},
									time.Hour*1,
								),
							),
						},
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				email:         "email-changed@test.ch",
			},
			want: &domain.Email{
				ObjectRoot: models.ObjectRoot{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				EmailAddress:    "email-changed@test.ch",
				IsEmailVerified: false,
				PlainCode:       gu.Ptr("a"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := c.changeUserEmail(context.Background(), tt.args.userID, tt.args.resourceOwner, tt.args.email, GetMockSecretGenerator(t), true, nil)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestCommands_NewUserEmailEvents(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		userID        string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "missing userID",
			fields: fields{
				eventstore: eventstoreExpect(t),
			},
			args: args{
				userID:        "",
				resourceOwner: "org1",
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "COMMAND-0Gzs3", "Errors.User.Email.IDMissing"),
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t, expectFilter()),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: caos_errs.ThrowNotFound(nil, "COMMAND-ieJ2e", "Errors.User.Email.NotFound"),
		},
		{
			name: "user not initialized",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanInitialCodeAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								nil, time.Hour*1,
							),
						),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-uz0Uu", "Errors.User.NotInitialised"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			_, err := c.NewUserEmailEvents(context.Background(), tt.args.userID, tt.args.resourceOwner)
			require.ErrorIs(t, err, tt.wantErr)
			// success cases are tested through the Change commands.
		})
	}
}
