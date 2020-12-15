package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/org_iam"
)

func (r *CommandSide) AddDefaultOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error) {
	addedPolicy := org_iam.NewWriteModel(policy.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy != nil {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-Lk0dS", "Errors.IAM.OrgIAMPolicy.AlreadyExists")
	}

	iamAgg := AggregateFromWriteModel(&addedPolicy.WriteModel.WriteModel).
		PushOrgIAMPolicyAddedEvent(ctx, policy.UserLoginMustBeDomain)

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToOrgIAMPolicy(addedPolicy), nil
}

func (r *CommandSide) ChangeDefaultOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error) {
	existingPolicy, err := r.defaultOrgIAMPolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	iamAgg := AggregateFromWriteModel(&existingPolicy.WriteModel.WriteModel).
		PushOrgIAMPolicyChangedFromExisting(ctx, existingPolicy, policy.UserLoginMustBeDomain)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToOrgIAMPolicy(existingPolicy), nil
}

func (r *CommandSide) defaultOrgIAMPolicyWriteModelByID(ctx context.Context, iamID string) (policy *org_iam.WriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := org_iam.NewWriteModel(iamID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
