package org

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/repository/policy"
)

var (
	LockoutPolicyAddedEventType   = orgEventTypePrefix + policy.LockoutPolicyAddedEventType
	LockoutPolicyChangedEventType = orgEventTypePrefix + policy.LockoutPolicyChangedEventType
	LockoutPolicyRemovedEventType = orgEventTypePrefix + policy.LockoutPolicyRemovedEventType
)

type LockoutPolicyAddedEvent struct {
	policy.LockoutPolicyAddedEvent
}

func NewLockoutPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	maxAttempts uint64,
	showLockoutFailure bool,
) *LockoutPolicyAddedEvent {
	return &LockoutPolicyAddedEvent{
		LockoutPolicyAddedEvent: *policy.NewLockoutPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LockoutPolicyAddedEventType),
			maxAttempts,
			showLockoutFailure),
	}
}

func LockoutPolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LockoutPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LockoutPolicyAddedEvent{LockoutPolicyAddedEvent: *e.(*policy.LockoutPolicyAddedEvent)}, nil
}

type LockoutPolicyChangedEvent struct {
	policy.LockoutPolicyChangedEvent
}

func NewLockoutPolicyChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []policy.LockoutPolicyChanges,
) (*LockoutPolicyChangedEvent, error) {
	changedEvent, err := policy.NewLockoutPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			LockoutPolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &LockoutPolicyChangedEvent{LockoutPolicyChangedEvent: *changedEvent}, nil
}

func LockoutPolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LockoutPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LockoutPolicyChangedEvent{LockoutPolicyChangedEvent: *e.(*policy.LockoutPolicyChangedEvent)}, nil
}

type LockoutPolicyRemovedEvent struct {
	policy.LockoutPolicyRemovedEvent
}

func NewLockoutPolicyRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *LockoutPolicyRemovedEvent {
	return &LockoutPolicyRemovedEvent{
		LockoutPolicyRemovedEvent: *policy.NewLockoutPolicyRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LockoutPolicyRemovedEventType),
		),
	}
}

func LockoutPolicyRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.LockoutPolicyRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LockoutPolicyRemovedEvent{LockoutPolicyRemovedEvent: *e.(*policy.LockoutPolicyRemovedEvent)}, nil
}
