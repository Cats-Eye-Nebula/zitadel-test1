package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *CommandSide) AddDefaultLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-5Mv0s", "Errors.IAM.LabelPolicyInvalid")
	}

	addedPolicy := NewIAMLabelPolicyWriteModel(policy.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.IsActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LabelPolicy.AlreadyExists")
	}

	iamAgg := AggregateFromWriteModel(&addedPolicy.LabelPolicyWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewLabelPolicyAddedEvent(ctx, policy.PrimaryColor, policy.SecondaryColor))

	err = r.eventstore.PushAggregate(ctx, addedPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLabelPolicy(addedPolicy), nil
}

func (r *CommandSide) ChangeDefaultLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-6M0od", "Errors.IAM.LabelPolicyInvalid")
	}

	existingPolicy, err := r.defaultLabelPolicyWriteModelByID(ctx, policy.AggregateID)
	if err != nil {
		return nil, err
	}

	if !existingPolicy.IsActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-0K9dq", "Errors.IAM.LabelPolicy.NotFound")
	}
	if !existingPolicy.HasChanged(policy.PrimaryColor, policy.SecondaryColor) {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}

	iamAgg := AggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	iamAgg.PushEvents(iam_repo.NewLabelPolicyChangedEvent(ctx, policy.PrimaryColor, policy.SecondaryColor))

	err = r.eventstore.PushAggregate(ctx, existingPolicy, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLabelPolicy(existingPolicy), nil
}

func (r *CommandSide) defaultLabelPolicyWriteModelByID(ctx context.Context, iamID string) (policy *IAMLabelPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMLabelPolicyWriteModel(iamID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
