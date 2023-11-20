package org

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

const (
	orgEventTypePrefix = eventstore.EventType("org.")
)

const (
	AggregateType    = "org"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(id string) *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:          AggregateType,
			Version:       AggregateVersion,
			ID:            id,
			ResourceOwner: id,
		},
	}
}
