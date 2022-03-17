package command

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
)

type Step5 struct {
	DefaultOrgIAMPolicy domain.OrgIAMPolicy
}

func (s *Step5) Step() domain.Step {
	return domain.Step5
}

func (s *Step5) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep5(ctx, s)
}

func (c *Commands) SetupStep5(ctx context.Context, step *Step5) error {
	fn := func(iam *InstanceWriteModel) ([]eventstore.Command, error) {
		iamAgg := InstanceAggregateFromWriteModel(&iam.WriteModel)
		event, err := c.addDefaultOrgIAMPolicy(ctx, iamAgg, NewInstanceOrgIAMPolicyWriteModel(), &domain.OrgIAMPolicy{
			UserLoginMustBeDomain: step.DefaultOrgIAMPolicy.UserLoginMustBeDomain,
		})
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-ADgd2").Info("default org iam policy set up")
		return []eventstore.Command{event}, nil
	}
	return c.setup(ctx, step, fn)
}
