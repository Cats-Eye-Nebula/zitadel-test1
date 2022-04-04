package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/policy"
)

type InstanceLoginPolicyWriteModel struct {
	LoginPolicyWriteModel
}

func NewInstanceLoginPolicyWriteModel(ctx context.Context) *InstanceLoginPolicyWriteModel {
	return &InstanceLoginPolicyWriteModel{
		LoginPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
			},
		},
	}
}

func (wm *InstanceLoginPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.LoginPolicyAddedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyAddedEvent)
		case *instance.LoginPolicyChangedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyChangedEvent)
		}
	}
}

func (wm *InstanceLoginPolicyWriteModel) IsValid() bool {
	return wm.AggregateID != ""
}

func (wm *InstanceLoginPolicyWriteModel) Reduce() error {
	return wm.LoginPolicyWriteModel.Reduce()
}

func (wm *InstanceLoginPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.LoginPolicyWriteModel.AggregateID).
		EventTypes(
			instance.LoginPolicyAddedEventType,
			instance.LoginPolicyChangedEventType).
		Builder()
}

func (wm *InstanceLoginPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	allowUsernamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA,
	hidePasswordReset bool,
	passwordlessType domain.PasswordlessType,
	passwordCheckLifetime,
	externalLoginCheckLifetime,
	mfaInitSkipLifetime,
	secondFactorCheckLifetime,
	multiFactorCheckLifetime time.Duration,
) (*instance.LoginPolicyChangedEvent, bool) {

	changes := make([]policy.LoginPolicyChanges, 0)
	if wm.AllowUserNamePassword != allowUsernamePassword {
		changes = append(changes, policy.ChangeAllowUserNamePassword(allowUsernamePassword))
	}
	if wm.AllowRegister != allowRegister {
		changes = append(changes, policy.ChangeAllowRegister(allowRegister))
	}
	if wm.AllowExternalIDP != allowExternalIDP {
		changes = append(changes, policy.ChangeAllowExternalIDP(allowExternalIDP))
	}
	if wm.ForceMFA != forceMFA {
		changes = append(changes, policy.ChangeForceMFA(forceMFA))
	}
	if passwordlessType.Valid() && wm.PasswordlessType != passwordlessType {
		changes = append(changes, policy.ChangePasswordlessType(passwordlessType))
	}
	if wm.HidePasswordReset != hidePasswordReset {
		changes = append(changes, policy.ChangeHidePasswordReset(hidePasswordReset))
	}
	if wm.PasswordCheckLifetime != passwordCheckLifetime {
		changes = append(changes, policy.ChangePasswordCheckLifetime(passwordCheckLifetime))
	}
	if wm.ExternalLoginCheckLifetime != externalLoginCheckLifetime {
		changes = append(changes, policy.ChangeExternalLoginCheckLifetime(externalLoginCheckLifetime))
	}
	if wm.MFAInitSkipLifetime != mfaInitSkipLifetime {
		changes = append(changes, policy.ChangeMFAInitSkipLifetime(mfaInitSkipLifetime))
	}
	if wm.SecondFactorCheckLifetime != secondFactorCheckLifetime {
		changes = append(changes, policy.ChangeSecondFactorCheckLifetime(secondFactorCheckLifetime))
	}
	if wm.MultiFactorCheckLifetime != multiFactorCheckLifetime {
		changes = append(changes, policy.ChangeMultiFactorCheckLifetime(multiFactorCheckLifetime))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := instance.NewLoginPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
