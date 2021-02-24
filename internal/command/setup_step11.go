package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
)

type Step11 struct {
	MigrateV1EventstoreToV2 bool
}

func (s *Step11) Step() domain.Step {
	return domain.Step11
}

func (s *Step11) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep11(ctx, s)
}

func (c *Commands) SetupStep11(ctx context.Context, step *Step11) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		var uniqueContraintMigrations []*domain.UniqueConstraintMigration
		if step.MigrateV1EventstoreToV2 {
			uniqueConstraints := NewUniqueConstraintReadModel(ctx, c)
			err := c.eventstore.FilterToQueryReducer(ctx, uniqueConstraints)
			if err != nil {
				return nil, err
			}
			uniqueContraintMigrations = uniqueConstraints.UniqueConstraints
		}
		logging.Log("SETUP-M9fsd").Info("migrate v1 eventstore to v2")
		return []eventstore.EventPusher{iam_repo.NewMigrateUniqueConstraintEvent(ctx, iamAgg, uniqueContraintMigrations)}, nil
	}
	return c.setup(ctx, step, fn)
}
