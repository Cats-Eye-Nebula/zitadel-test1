package eventsourcing

import (
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/key/repository/eventsourcing/model"
)

func KeyPairQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.KeyPairAggregate).
		LatestSequenceFilter(latestSequence)
}
