package command

import (
	"context"

	"github.com/caos/logging"
	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/domain"
)

type Step9 struct {
	Passwordless bool
}

func (s *Step9) Step() domain.Step {
	return domain.Step9
}

func (s *Step9) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep9(ctx, s)
}

func (c *Commands) SetupStep9(ctx context.Context, step *Step9) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.Command, error) {
		multiFactorModel := NewIAMMultiFactorWriteModel(domain.MultiFactorTypeU2FWithPIN)
		iamAgg := IAMAggregateFromWriteModel(&multiFactorModel.MultiFactorWriteModel.WriteModel)
		if !step.Passwordless {
			return []eventstore.Command{}, nil
		}
		passwordlessEvent, err := setPasswordlessAllowedInPolicy(ctx, c, iamAgg)
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-AEG2t").Info("allowed passwordless in login policy")
		multifactorEvent, err := c.addMultiFactorToDefaultLoginPolicy(ctx, iamAgg, multiFactorModel, domain.MultiFactorTypeU2FWithPIN)
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-ADfng").Info("added passwordless to MFA login policy")
		return []eventstore.Command{passwordlessEvent, multifactorEvent}, nil
	}
	return c.setup(ctx, step, fn)
}

func setPasswordlessAllowedInPolicy(ctx context.Context, c *Commands, iamAgg *eventstore.Aggregate) (eventstore.Command, error) {
	policy, err := c.getDefaultLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	policy.PasswordlessType = domain.PasswordlessTypeAllowed
	return c.changeDefaultLoginPolicy(ctx, iamAgg, NewIAMLoginPolicyWriteModel(), policy)
}
