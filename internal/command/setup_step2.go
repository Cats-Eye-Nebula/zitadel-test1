package command

import (
	"context"

	"github.com/caos/logging"

	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/domain"
)

type Step2 struct {
	DefaultPasswordComplexityPolicy domain.PasswordComplexityPolicy
}

func (s *Step2) Step() domain.Step {
	return domain.Step2
}

func (s *Step2) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep2(ctx, s)
}

func (c *Commands) SetupStep2(ctx context.Context, step *Step2) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.Command, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		event, err := c.addDefaultPasswordComplexityPolicy(ctx, iamAgg, NewIAMPasswordComplexityPolicyWriteModel(), &domain.PasswordComplexityPolicy{
			MinLength:    step.DefaultPasswordComplexityPolicy.MinLength,
			HasLowercase: step.DefaultPasswordComplexityPolicy.HasLowercase,
			HasUppercase: step.DefaultPasswordComplexityPolicy.HasUppercase,
			HasNumber:    step.DefaultPasswordComplexityPolicy.HasNumber,
			HasSymbol:    step.DefaultPasswordComplexityPolicy.HasSymbol,
		})
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-ADgd2").Info("default password complexity policy set up")
		return []eventstore.Command{event}, nil
	}
	return c.setup(ctx, step, fn)
}
