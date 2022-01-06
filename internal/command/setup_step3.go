package command

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
)

type Step3 struct {
	DefaultPasswordAgePolicy domain.PasswordAgePolicy
}

func (s *Step3) Step() domain.Step {
	return domain.Step3
}

func (s *Step3) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep3(ctx, s)
}

func (c *Commands) SetupStep3(ctx context.Context, step *Step3) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.Command, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		event, err := c.addDefaultPasswordAgePolicy(ctx, iamAgg, NewIAMPasswordAgePolicyWriteModel(), &domain.PasswordAgePolicy{
			MaxAgeDays:     step.DefaultPasswordAgePolicy.MaxAgeDays,
			ExpireWarnDays: step.DefaultPasswordAgePolicy.ExpireWarnDays,
		})
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-DBqgq").Info("default password age policy set up")
		return []eventstore.Command{event}, nil
	}
	return c.setup(ctx, step, fn)
}
