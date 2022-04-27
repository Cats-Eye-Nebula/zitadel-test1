package iam

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/features"
)

var (
	FeaturesSetEventType = iamEventTypePrefix + features.FeaturesSetEventType
)

type FeaturesSetEvent struct {
	features.FeaturesSetEvent
}

func NewFeaturesSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []features.FeaturesChanges,
) (*FeaturesSetEvent, error) {
	changedEvent, err := features.NewFeaturesSetEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			FeaturesSetEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &FeaturesSetEvent{FeaturesSetEvent: *changedEvent}, nil
}

func FeaturesSetEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := features.FeaturesSetEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &FeaturesSetEvent{FeaturesSetEvent: *e.(*features.FeaturesSetEvent)}, nil
}
