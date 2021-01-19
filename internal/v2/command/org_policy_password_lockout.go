package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

func (r *CommandSide) AddPasswordLockoutPolicy(ctx context.Context, policy *domain.PasswordLockoutPolicy) (*domain.PasswordLockoutPolicy, error) {
	addedPolicy := NewOrgPasswordLockoutPolicyWriteModel(policy.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "ORG-0olDf", "Errors.ORG.PasswordLockoutPolicy.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.WriteModel)
	orgAgg.PushEvents(org.NewPasswordLockoutPolicyAddedEvent(ctx, policy.MaxAttempts, policy.ShowLockOutFailures))

	err = r.eventstore.PushAggregate(ctx, addedPolicy, orgAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordLockoutPolicy(&addedPolicy.PasswordLockoutPolicyWriteModel), nil
}

func (r *CommandSide) ChangePasswordLockoutPolicy(ctx context.Context, policy *domain.PasswordLockoutPolicy) (*domain.PasswordLockoutPolicy, error) {
	existingPolicy := NewOrgPasswordLockoutPolicyWriteModel(policy.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-ADfs1", "Errors.Org.PasswordLockoutPolicy.NotFound")
	}

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.MaxAttempts, policy.ShowLockOutFailures)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-4M9vs", "Errors.Org.PasswordLockoutPolicy.NotChanged")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.PasswordLockoutPolicyWriteModel.WriteModel)
	orgAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, orgAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToPasswordLockoutPolicy(&existingPolicy.PasswordLockoutPolicyWriteModel), nil
}

func (r *CommandSide) RemovePasswordLockoutPolicy(ctx context.Context, orgID string) error {
	existingPolicy := NewOrgPasswordLockoutPolicyWriteModel(orgID)
	err := r.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return caos_errs.ThrowNotFound(nil, "ORG-D4zuz", "Errors.Org.PasswordLockoutPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	orgAgg.PushEvents(org.NewPasswordLockoutPolicyRemovedEvent(ctx))
	return r.eventstore.PushAggregate(ctx, existingPolicy, orgAgg)
}
