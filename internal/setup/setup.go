package setup

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

func Execute(ctx context.Context, setUpConfig IAMSetUp, iamID string, commands *command.Commands) error {
	logging.Log("SETUP-JAK2q").Info("starting setup")

	iam, err := commands.GetIAM(ctx)
	if err != nil && !caos_errs.IsNotFound(err) {
		return err
	}
	if iam != nil && (iam.SetUpDone == domain.StepCount-1 || iam.SetUpStarted != iam.SetUpDone) {
		logging.Log("SETUP-VA2k1").Info("all steps done")
		return nil
	}

	if iam == nil {
		iam = &domain.IAM{ObjectRoot: models.ObjectRoot{AggregateID: iamID}}
	}

	steps, err := setUpConfig.Steps(iam.SetUpDone)
	if err != nil || len(steps) == 0 {
		return err
	}

	err = commands.ExecuteSetupSteps(ctx, steps)
	if err != nil {
		return err
	}

	logging.Log("SETUP-ds31h").Info("setup done")
	return nil
}
