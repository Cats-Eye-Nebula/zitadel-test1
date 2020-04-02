package eventstore

import (
	"github.com/caos/zitadel/internal/eventstore/internal/repository/sql"
	"github.com/caos/zitadel/internal/eventstore/models"
)

type Config struct {
	Repository  sql.Config
	ServiceName string
}

func Start(conf Config) (Eventstore, error) {
	repo, err := sql.Start(conf.Repository)
	if err != nil {
		return nil, err
	}

	return &eventstore{
		repo:             repo,
		aggregateCreator: models.NewAggregateCreator(conf.ServiceName),
	}, nil
}
