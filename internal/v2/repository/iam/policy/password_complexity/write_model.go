package password_complexity

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_complexity"
)

const (
	AggregateType = "iam"
)

type WriteModel struct {
	password_complexity.WriteModel
}

func NewWriteModel(iamID string) *WriteModel {
	return &WriteModel{
		password_complexity.WriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: iamID,
			},
		},
	}
}

func (wm *WriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *AddedEvent:
			wm.WriteModel.AppendEvents(&e.AddedEvent)
		case *ChangedEvent:
			wm.WriteModel.AppendEvents(&e.ChangedEvent)
		}
	}
}

func (wm *WriteModel) Reduce() error {
	return wm.WriteModel.Reduce()
}

func (wm *WriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID)
}
