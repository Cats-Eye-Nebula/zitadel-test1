package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type IAMPasswordLockoutPolicyWriteModel struct {
	PasswordLockoutPolicyWriteModel
}

func NewIAMPasswordLockoutPolicyWriteModel() *IAMPasswordLockoutPolicyWriteModel {
	return &IAMPasswordLockoutPolicyWriteModel{
		PasswordLockoutPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMPasswordLockoutPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.PasswordLockoutPolicyAddedEvent:
			wm.PasswordLockoutPolicyWriteModel.AppendEvents(&e.PasswordLockoutPolicyAddedEvent)
		case *iam.PasswordLockoutPolicyChangedEvent:
			wm.PasswordLockoutPolicyWriteModel.AppendEvents(&e.PasswordLockoutPolicyChangedEvent)
		}
	}
}

func (wm *IAMPasswordLockoutPolicyWriteModel) Reduce() error {
	return wm.PasswordLockoutPolicyWriteModel.Reduce()
}

func (wm *IAMPasswordLockoutPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.PasswordLockoutPolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

func (wm *IAMPasswordLockoutPolicyWriteModel) NewChangedEvent(ctx context.Context, maxAttempts uint64, showLockoutFailure bool) (*iam.PasswordLockoutPolicyChangedEvent, bool) {
	hasChanged := false
	changedEvent := iam.NewPasswordLockoutPolicyChangedEvent(ctx)
	if wm.MaxAttempts != maxAttempts {
		hasChanged = true
		changedEvent.MaxAttempts = &maxAttempts
	}
	if wm.ShowLockOutFailures != showLockoutFailure {
		hasChanged = true
		changedEvent.ShowLockOutFailures = &showLockoutFailure
	}
	return changedEvent, hasChanged
}
