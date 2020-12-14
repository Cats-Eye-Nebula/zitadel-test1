package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

var (
	MemberAddedEventType   = IamEventTypePrefix + member.AddedEventType
	MemberChangedEventType = IamEventTypePrefix + member.ChangedEventType
	MemberRemovedEventType = IamEventTypePrefix + member.RemovedEventType
)

type MemberReadModel struct {
	member.ReadModel

	userID string
	iamID  string
}

func NewMemberReadModel(iamID, userID string) *MemberReadModel {
	return &MemberReadModel{
		iamID:  iamID,
		userID: userID,
	}
}

func (rm *MemberReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *MemberAddedEvent:
			rm.ReadModel.AppendEvents(&e.AddedEvent)
		case *MemberChangedEvent:
			rm.ReadModel.AppendEvents(&e.ChangedEvent)
		case *member.AddedEvent, *member.ChangedEvent, *MemberRemovedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

func (rm *MemberReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(rm.iamID).
		EventData(map[string]interface{}{
			"userId": rm.userID,
		})
}

type MemberWriteModel struct {
	member.WriteModel
}

func NewMemberWriteModel(iamID, userID string) *MemberWriteModel {
	return &MemberWriteModel{
		member.WriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: iamID,
			},
			UserID: userID,
		},
	}
}

func (wm *MemberWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *MemberAddedEvent:
			if e.UserID != wm.WriteModel.UserID {
				continue
			}
			wm.WriteModel.AppendEvents(&e.AddedEvent)
		case *MemberChangedEvent:
			if e.UserID != wm.WriteModel.UserID {
				continue
			}
			wm.WriteModel.AppendEvents(&e.ChangedEvent)
		case *MemberRemovedEvent:
			if e.UserID != wm.WriteModel.UserID {
				continue
			}
			wm.WriteModel.AppendEvents(&e.RemovedEvent)
		}
	}
}

func (wm *MemberWriteModel) Reduce() error {
	return wm.WriteModel.Reduce()
}

func (wm *MemberWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID)
}

type MemberAddedEvent struct {
	member.AddedEvent
}

func NewMemberAddedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *MemberAddedEvent {

	return &MemberAddedEvent{
		AddedEvent: *member.NewAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberAddedEventType,
			),
			userID,
			roles...,
		),
	}
}

func MemberAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MemberAddedEvent{AddedEvent: *e.(*member.AddedEvent)}, nil
}

type MemberChangedEvent struct {
	member.ChangedEvent
}

func MemberChangedEventFromExisting(
	ctx context.Context,
	current *MemberWriteModel,
	roles ...string,
) (*MemberChangedEvent, error) {

	event, err := member.ChangeEventFromExisting(
		eventstore.NewBaseEventForPush(
			ctx,
			MemberChangedEventType,
		),
		&current.WriteModel,
		roles...,
	)
	if err != nil {
		return nil, err
	}

	return &MemberChangedEvent{
		ChangedEvent: *event,
	}, nil
}

func MemberChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.ChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MemberChangedEvent{ChangedEvent: *e.(*member.ChangedEvent)}, nil
}

type MemberRemovedEvent struct {
	member.RemovedEvent
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

func MemberRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MemberRemovedEvent{RemovedEvent: *e.(*member.RemovedEvent)}, nil
}
