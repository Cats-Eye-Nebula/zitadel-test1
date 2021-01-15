package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type MachineWriteModel struct {
	eventstore.WriteModel

	UserName string

	Name        string
	Description string
	UserState   domain.UserState
}

func NewMachineWriteModel(userID, resourceOwner string) *MachineWriteModel {
	return &MachineWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *MachineWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
}

//TODO: Compute OTPState? initial/active
func (wm *MachineWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.MachineAddedEvent:
			wm.UserName = e.UserName
			wm.Name = e.Name
			wm.Description = e.Description
			wm.UserState = domain.UserStateActive
		case *user.UsernameChangedEvent:
			wm.UserName = e.UserName
		case *user.MachineChangedEvent:
			if e.Name != nil {
				wm.Name = *e.Name
			}
			if e.Description != nil {
				wm.Description = *e.Description
			}
		case *user.UserLockedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateLocked
			}
		case *user.UserUnlockedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateActive
			}
		case *user.UserDeactivatedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateInactive
			}
		case *user.UserReactivatedEvent:
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateActive
			}
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *MachineWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

func (wm *MachineWriteModel) NewChangedEvent(
	ctx context.Context,
	name,
	description string,
) (*user.MachineChangedEvent, bool) {
	hasChanged := false
	changedEvent := user.NewMachineChangedEvent(ctx)
	if wm.Name != name {
		hasChanged = true
		changedEvent.Name = &name
	}
	if wm.Description != description {
		hasChanged = true
		changedEvent.Description = &description
	}
	return changedEvent, hasChanged
}
