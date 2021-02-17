package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type IAMWriteModel struct {
	eventstore.WriteModel

	SetUpStarted domain.Step
	SetUpDone    domain.Step

	GlobalOrgID string
	ProjectID   string
}

func NewIAMWriteModel() *IAMWriteModel {
	return &IAMWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   domain.IAMID,
			ResourceOwner: domain.IAMID,
		},
	}
}

func (wm *IAMWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *iam.ProjectSetEvent:
			wm.ProjectID = e.ProjectID
		case *iam.GlobalOrgSetEvent:
			wm.GlobalOrgID = e.OrgID
		case *iam.SetupStepEvent:
			if e.Done {
				wm.SetUpDone = e.Step
			} else {
				wm.SetUpStarted = e.Step
			}
		}
	}
	return nil
}

func (wm *IAMWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			iam.ProjectSetEventType,
			iam.GlobalOrgSetEventType,
			iam.SetupStartedEventType,
			iam.SetupDoneEventType)
}

func IAMAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, iam.AggregateType, iam.AggregateVersion)
}
