package command

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/domain"
)

type Step6 struct {
	DefaultLabelPolicy domain.LabelPolicy
}

func (s *Step6) Step() domain.Step {
	return domain.Step6
}

func (s *Step6) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep6(ctx, s)
}

func (c *Commands) SetupStep6(ctx context.Context, step *Step6) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.Command, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		event, err := c.addDefaultLabelPolicy(ctx, iamAgg, NewIAMLabelPolicyWriteModel(), &step.DefaultLabelPolicy)
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-ADgd2").Info("default label policy set up")
		return []eventstore.Command{event}, nil
	}
	return c.setup(ctx, step, fn)
}
