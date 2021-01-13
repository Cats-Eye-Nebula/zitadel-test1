package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type OrgIdentityProviderWriteModel struct {
	IdentityProviderWriteModel
}

func NewOrgIdentityProviderWriteModel(orgID, idpConfigID string) *OrgIdentityProviderWriteModel {
	return &OrgIdentityProviderWriteModel{
		IdentityProviderWriteModel: IdentityProviderWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			IDPConfigID: idpConfigID,
		},
	}
}

func (wm *OrgIdentityProviderWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.IdentityProviderAddedEvent:
			if e.IDPConfigID != wm.IDPConfigID {
				continue
			}
			wm.IdentityProviderWriteModel.AppendEvents(&e.IdentityProviderAddedEvent)
		}
	}
}

func (wm *OrgIdentityProviderWriteModel) Reduce() error {
	return wm.IdentityProviderWriteModel.Reduce()
}

func (wm *OrgIdentityProviderWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}
