package command

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
)

//TODO: make in login policy available
type Step8 struct {
	U2F bool
}

func (s *Step8) Step() domain.Step {
	return domain.Step8
}

func (s *Step8) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep8(ctx, s)
}

func (c *Commands) SetupStep8(ctx context.Context, step *Step8) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.Command, error) {
		secondFactorModel := NewIAMSecondFactorWriteModel(domain.SecondFactorTypeU2F)
		iamAgg := IAMAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
		if !step.U2F {
			return []eventstore.Command{}, nil
		}
		event, err := c.addSecondFactorToDefaultLoginPolicy(ctx, iamAgg, secondFactorModel, domain.SecondFactorTypeU2F)
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-BDhne").Info("added U2F to 2FA login policy")
		return []eventstore.Command{event}, nil
	}
	return c.setup(ctx, step, fn)
}
