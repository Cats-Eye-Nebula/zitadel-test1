package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/iam/model"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
)

type IamRepo struct {
	IAMID     string
	IAMEvents *iam_event.IAMEventstore
}

func (repo *IamRepo) Health(ctx context.Context) error {
	return repo.IAMEvents.Health(ctx)
}

func (repo *IamRepo) IamByID(ctx context.Context) (*model.IAM, error) {
	return repo.IAMEvents.IAMByID(ctx, repo.IAMID)
}
