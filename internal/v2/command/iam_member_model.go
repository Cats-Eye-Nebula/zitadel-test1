package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type IAMMemberWriteModel struct {
	MemberWriteModel
}

func NewIAMMemberWriteModel(userID string) *IAMMemberWriteModel {
	return &IAMMemberWriteModel{
		MemberWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
			UserID: userID,
		},
	}
}

func (wm *IAMMemberWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.MemberAddedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberAddedEvent)
		case *iam.MemberChangedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberChangedEvent)
		case *iam.MemberRemovedEvent:
			if e.UserID != wm.MemberWriteModel.UserID {
				continue
			}
			wm.MemberWriteModel.AppendEvents(&e.MemberRemovedEvent)
		}
	}
}

func (wm *IAMMemberWriteModel) Reduce() error {
	return wm.MemberWriteModel.Reduce()
}

func (wm *IAMMemberWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.MemberWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			iam.MemberAddedEventType,
			iam.MemberChangedEventType,
			iam.MemberRemovedEventType)
}
