package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type OrgLabelPolicyWriteModel struct {
	LabelPolicyWriteModel
}

func NewOrgLabelPolicyWriteModel(orgID string) *OrgLabelPolicyWriteModel {
	return &OrgLabelPolicyWriteModel{
		LabelPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgLabelPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.LabelPolicyAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyAddedEvent)
		case *org.LabelPolicyChangedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyChangedEvent)
		}
	}
}

func (wm *OrgLabelPolicyWriteModel) Reduce() error {
	return wm.LabelPolicyWriteModel.Reduce()
}

func (wm *OrgLabelPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.LabelPolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			org.LabelPolicyAddedEventType,
			org.LabelPolicyChangedEventType)
}

func (wm *OrgLabelPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	primaryColor,
	secondaryColor string,
) (*org.LabelPolicyChangedEvent, bool) {
	changes := make([]policy.LabelPolicyChanges, 0)
	if wm.PrimaryColor != primaryColor {
		changes = append(changes, policy.ChangePrimaryColor(primaryColor))
	}
	if wm.SecondaryColor != secondaryColor {
		changes = append(changes, policy.ChangeSecondaryColor(secondaryColor))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := org.NewLabelPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
