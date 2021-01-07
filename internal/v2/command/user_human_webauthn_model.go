package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type HumanWebAuthNWriteModel struct {
	eventstore.WriteModel

	WebauthNTokenID string

	UserState domain.UserState
}

func NewHumanWebAuthNWriteModel(userID, wbAuthNTokenID string) *HumanWebAuthNWriteModel {
	return &HumanWebAuthNWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: userID,
		},
		WebauthNTokenID: wbAuthNTokenID,
	}
}

func (wm *HumanWebAuthNWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanWebAuthNAddedEvent:
			if wm.WebauthNTokenID == e.WebAuthNTokenID {
				wm.AppendEvents(e)
			}
		case *user.HumanWebAuthNRemovedEvent:
			if wm.WebauthNTokenID == e.WebAuthNTokenID {
				wm.AppendEvents(e)
			}
		case *user.UserRemovedEvent:
			wm.AppendEvents(e)
		}
	}
}

func (wm *HumanWebAuthNWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanWebAuthNAddedEvent:
			wm.WebauthNTokenID = e.WebAuthNTokenID
			wm.UserState = domain.UserStateActive
		case *user.HumanWebAuthNRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanWebAuthNWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID)
}
