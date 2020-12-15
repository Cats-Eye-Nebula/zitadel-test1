package org

import (
	"context"
	"github.com/caos/zitadel/internal/v2/business/command"
	"github.com/caos/zitadel/internal/v2/business/query"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

const (
	orgEventTypePrefix = eventstore.EventType("org.")
)

var (
	MemberAddedEventType   = orgEventTypePrefix + member.AddedEventType
	MemberChangedEventType = orgEventTypePrefix + member.ChangedEventType
	MemberRemovedEventType = orgEventTypePrefix + member.RemovedEventType
)

type MemberWriteModel struct {
	command.MemberWriteModel
}

// func NewMemberAggregate(userID string) *MemberAggregate {
// 	return &MemberAggregate{
// 		Aggregate: member.NewAggregate(
// 			eventstore.NewAggregate(userID, MemberAggregateType, "RO", AggregateVersion, 0),
// 		),
// 		// Aggregate: member.NewMemberAggregate(userID),
// 	}
// }

type MembersReadModel struct {
	query.MembersReadModel
}

func (rm *MembersReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *MemberAddedEvent:
			rm.MembersReadModel.AppendEvents(&e.MemberAddedEvent)
		case *MemberChangedEvent:
			rm.MembersReadModel.AppendEvents(&e.ChangedEvent)
		case *MemberRemovedEvent:
			rm.MembersReadModel.AppendEvents(&e.RemovedEvent)
		}
	}
}

type MemberReadModel query.MemberReadModel

func (rm *MemberReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *MemberAddedEvent:
			rm.ReadModel.AppendEvents(&e.MemberAddedEvent)
		case *MemberChangedEvent:
			rm.ReadModel.AppendEvents(&e.ChangedEvent)
		}
	}
}

type MemberAddedEvent struct {
	member.MemberAddedEvent
}

type MemberChangedEvent struct {
	member.ChangedEvent
}
type MemberRemovedEvent struct {
	member.RemovedEvent
}

func NewMemberAddedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *MemberAddedEvent {

	return &MemberAddedEvent{
		MemberAddedEvent: *member.NewMemberAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberAddedEventType,
			),
			userID,
			roles...,
		),
	}
}

func MemberChangedEventFromExisting(
	ctx context.Context,
	current *MemberWriteModel,
	roles ...string,
) (*MemberChangedEvent, error) {

	m, err := member.ChangeEventFromExisting(
		eventstore.NewBaseEventForPush(
			ctx,
			MemberChangedEventType,
		),
		&current.MemberWriteModel,
		roles...,
	)
	if err != nil {
		return nil, err
	}

	return &MemberChangedEvent{
		ChangedEvent: *m,
	}, nil
}

func NewMemberChangedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *MemberChangedEvent {

	return &MemberChangedEvent{
		ChangedEvent: *member.NewChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberChangedEventType,
			),
			userID,
			roles...,
		),
	}
}

func NewMemberRemovedEvent(
	ctx context.Context,
	userID string,
) *MemberRemovedEvent {

	return &MemberRemovedEvent{
		RemovedEvent: *member.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberRemovedEventType,
			),
			userID,
		),
	}
}
