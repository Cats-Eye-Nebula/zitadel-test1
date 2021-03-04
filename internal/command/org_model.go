package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

type OrgWriteModel struct {
	eventstore.WriteModel

	Name          string
	State         domain.OrgState
	PrimaryDomain string
}

func NewOrgWriteModel(orgID string) *OrgWriteModel {
	return &OrgWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   orgID,
			ResourceOwner: orgID,
		},
	}
}

func (wm *OrgWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *org.OrgAddedEvent:
			wm.Name = e.Name
			wm.State = domain.OrgStateActive
		case *org.OrgDeactivatedEvent:
			wm.State = domain.OrgStateInactive
		case *org.OrgChangedEvent:
			wm.Name = e.Name
		case *org.DomainPrimarySetEvent:
			wm.PrimaryDomain = e.Domain
		}
	}
	return nil
}

func (wm *OrgWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			org.OrgAddedEventType,
			org.OrgChangedEventType,
			org.OrgDomainPrimarySetEventType)
}

func OrgAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, org.AggregateType, org.AggregateVersion)
}
