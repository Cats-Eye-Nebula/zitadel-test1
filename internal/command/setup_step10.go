package command

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type Step10 struct {
	DefaultMailTemplate domain.MailTemplate
}

func (s *Step10) Step() domain.Step {
	return domain.Step10
}

func (s *Step10) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep10(ctx, s)
}

func (c *Commands) SetupStep10(ctx context.Context, step *Step10) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.Command, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		mailTemplateEvent, err := c.addDefaultMailTemplate(ctx, iamAgg, NewIAMMailTemplateWriteModel(), &step.DefaultMailTemplate)
		if err != nil {
			return nil, err
		}
		events := []eventstore.Command{
			mailTemplateEvent,
		}
		logging.Log("SETUP-3N9fs").Info("default mail template/text set up")
		return events, nil
	}
	return c.setup(ctx, step, fn)
}
