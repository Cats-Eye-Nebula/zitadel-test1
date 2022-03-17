package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
)

type Step4 struct {
	DefaultPasswordLockoutPolicy domain.LockoutPolicy
}

func (s *Step4) Step() domain.Step {
	return domain.Step4
}

func (s *Step4) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep4(ctx, s)
}

//This step should not be executed when a new instance is setup, because its not used anymore
//SetupStep4 is no op in favour of step 18.
//Password lockout policy is replaced by lockout policy
func (c *Commands) SetupStep4(ctx context.Context, step *Step4) error {
	fn := func(iam *InstanceWriteModel) ([]eventstore.Command, error) {
		return nil, nil
	}
	return c.setup(ctx, step, fn)
}
