package command

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	zitadel_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/restrictions"
)

func TestSetRestrictions(t *testing.T) {
	type fields func(*testing.T) (*eventstore.Eventstore, id.Generator)
	type args struct {
		ctx             context.Context
		setRestrictions *SetRestrictions
		defaultLanguage language.Tag
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
			name: "set new restrictions",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								restrictions.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&restrictions.NewAggregate("restrictions1", "instance1", "instance1").Aggregate,
										restrictions.SetEventType,
									),
									restrictions.ChangePublicOrgRegistrations(true),
								),
							),
						),
					),
					id_mock.NewIDGeneratorExpectIDs(t, "restrictions1")
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setRestrictions: &SetRestrictions{
					DisallowPublicOrgRegistration: gu.Ptr(true),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
			},
		},
		{
			name: "change restrictions",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								restrictions.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&restrictions.NewAggregate("restrictions1", "instance1", "instance1").Aggregate,
										restrictions.SetEventType,
									),
									restrictions.ChangePublicOrgRegistrations(true),
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								restrictions.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&restrictions.NewAggregate("restrictions1", "instance1", "instance1").Aggregate,
										restrictions.SetEventType,
									),
									restrictions.ChangePublicOrgRegistrations(false),
								),
							),
						),
					),
					nil
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setRestrictions: &SetRestrictions{
					DisallowPublicOrgRegistration: gu.Ptr(false),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
			},
		},
		{
			name: "set restrictions idempotency",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								restrictions.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&restrictions.NewAggregate("restrictions1", "instance1", "instance1").Aggregate,
										restrictions.SetEventType,
									),
									restrictions.ChangePublicOrgRegistrations(true),
								),
							),
						),
					),
					nil
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setRestrictions: &SetRestrictions{
					DisallowPublicOrgRegistration: gu.Ptr(true),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
			},
		},
		{
			name: "no restrictions defined",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							restrictions.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&restrictions.NewAggregate("restrictions1", "instance1", "instance1").Aggregate,
									restrictions.SetEventType,
								),
								restrictions.ChangePublicOrgRegistrations(true),
							),
						),
					),
				), nil
			},
			args: args{
				ctx:             authz.WithInstanceID(context.Background(), "instance1"),
				setRestrictions: &SetRestrictions{},
			},
			res: res{
				err: zitadel_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "unsupported language restricted",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							restrictions.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&restrictions.NewAggregate("restrictions1", "instance1", "instance1").Aggregate,
									restrictions.SetEventType,
								),
								restrictions.ChangeAllowedLanguages(SupportedLanguages),
							),
						),
					),
				), nil
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setRestrictions: &SetRestrictions{
					AllowedLanguages: []language.Tag{AllowedLanguage, UnsupportedLanguage},
				},
				defaultLanguage: AllowedLanguage,
			},
			res: res{
				err: zitadel_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "default language not allowed",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							restrictions.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&restrictions.NewAggregate("restrictions1", "instance1", "instance1").Aggregate,
									restrictions.SetEventType,
								),
								restrictions.ChangeAllowedLanguages(OnlyAllowedLanguages),
							),
						),
					),
				), nil
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setRestrictions: &SetRestrictions{
					AllowedLanguages: []language.Tag{DisallowedLanguage},
				},
				defaultLanguage: AllowedLanguage,
			},
			res: res{
				err: zitadel_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := new(Commands)
			r.eventstore, r.idGenerator = tt.fields(t)
			got, err := r.SetInstanceRestrictions(tt.args.ctx, tt.args.setRestrictions, tt.args.defaultLanguage)
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
