package command

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type Step2 struct {
	DefaultPasswordComplexityPolicy iam_model.PasswordComplexityPolicy
}

func (r *CommandSide) SetupStep2(ctx context.Context, iamID string, step Step2) error {

	//_, err = r.setup(ctx, iam, domain.Step1, iam_repo.NewSetupStepDoneEvent(ctx, domain.Step1))
	//return err
	return nil
}

//
//func (r *CommandSide) addDefaultPasswordComplexityPolicy(ctx context.Context, iam *IAMWriteModel, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
//	if !policy.IsValid() {
//		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-5Mv0s", "Errors.IAM.LoginPolicyInvalid")
//	}
//
//	addedPolicy := NewIAMLoginPolicyWriteModel(policy.AggregateID)
//	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
//	if err != nil {
//		return nil, err
//	}
//	if addedPolicy.IsActive {
//		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.AlreadyExists")
//	}
//
//	//iamAgg.PushEvents(iam_repo.NewLoginPolicyAddedEvent(ctx, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIdp, policy.ForceMFA, domain.PasswordlessType(policy.PasswordlessType)))
//	return nil, nil
//}
