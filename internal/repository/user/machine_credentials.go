package user

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	machineSecretPrefix             = machineEventPrefix + "credentials."
	MachineSecretSetType            = machineSecretPrefix + "set"
	MachineSecretRemovedType        = machineSecretPrefix + "removed"
	MachineSecretCheckSucceededType = machineSecretPrefix + "check.succeeded"
	MachineSecretCheckFailedType    = machineSecretPrefix + "check.failed"
)

type MachineSecretSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
}

func (e *MachineSecretSetEvent) Data() interface{} {
	return e
}

func (e *MachineSecretSetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMachineSecretSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	clientSecret *crypto.CryptoValue,
) *MachineSecretSetEvent {
	return &MachineSecretSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineSecretSetType,
		),
		ClientSecret: clientSecret,
	}
}

func MachineSecretSetEventMapper(event *repository.Event) (eventstore.Event, error) {
	credentialsSet := &MachineSecretSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, credentialsSet)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-lrv2di", "unable to unmarshal machine added")
	}

	return credentialsSet, nil
}

type MachineSecretRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *MachineSecretRemovedEvent) Data() interface{} {
	return e
}

func (e *MachineSecretRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMachineSecretRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *MachineSecretRemovedEvent {
	return &MachineSecretRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineSecretRemovedType,
		),
	}
}

func MachineSecretRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	credentialsRemoved := &MachineSecretRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, credentialsRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-lrv2ei", "unable to unmarshal machine added")
	}

	return credentialsRemoved, nil
}

type MachineSecretCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *MachineSecretCheckSucceededEvent) Data() interface{} {
	return e
}

func (e *MachineSecretCheckSucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMachineSecretCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *MachineSecretCheckSucceededEvent {
	return &MachineSecretCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineSecretCheckSucceededType,
		),
	}
}

func MachineSecretCheckSucceededEventMapper(event *repository.Event) (eventstore.Event, error) {
	check := &MachineSecretCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, check)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-x9000ja", "unable to unmarshal machine added")
	}

	return check, nil
}

type MachineSecretCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *MachineSecretCheckFailedEvent) Data() interface{} {
	return e
}

func (e *MachineSecretCheckFailedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMachineSecretCheckFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *MachineSecretCheckFailedEvent {
	return &MachineSecretCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineSecretCheckFailedType,
		),
	}
}

func MachineSecretCheckFailedEventMapper(event *repository.Event) (eventstore.Event, error) {
	check := &MachineSecretCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, check)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-x9000ja", "unable to unmarshal machine added")
	}

	return check, nil
}
