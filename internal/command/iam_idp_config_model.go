package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/idpconfig"
)

type IAMIDPConfigWriteModel struct {
	IDPConfigWriteModel
}

func NewIAMIDPConfigWriteModel(configID string) *IAMIDPConfigWriteModel {
	return &IAMIDPConfigWriteModel{
		IDPConfigWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
			ConfigID: configID,
		},
	}
}

func (wm *IAMIDPConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			iam.IDPConfigAddedEventType,
			iam.IDPConfigChangedEventType,
			iam.IDPConfigDeactivatedEventType,
			iam.IDPConfigReactivatedEventType,
			iam.IDPConfigRemovedEventType,
			iam.IDPOIDCConfigAddedEventType,
			iam.IDPOIDCConfigChangedEventType).
		Builder()
}

func (wm *IAMIDPConfigWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.IDPConfigAddedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *iam.IDPConfigChangedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigChangedEvent)
		case *iam.IDPConfigDeactivatedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigDeactivatedEvent)
		case *iam.IDPConfigReactivatedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigReactivatedEvent)
		case *iam.IDPConfigRemovedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		case *iam.IDPOIDCConfigAddedEvent:
			if wm.ConfigID != e.IDPConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.OIDCConfigAddedEvent)
		case *iam.IDPOIDCConfigChangedEvent:
			if wm.ConfigID != e.IDPConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.OIDCConfigChangedEvent)
		}
	}
}

func (wm *IAMIDPConfigWriteModel) Reduce() error {
	return wm.IDPConfigWriteModel.Reduce()
}

func (wm *IAMIDPConfigWriteModel) AppendAndReduce(events ...eventstore.EventReader) error {
	wm.AppendEvents(events...)
	return wm.Reduce()
}

func (wm *IAMIDPConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	configID,
	name string,
	stylingType domain.IDPConfigStylingType,
) (*iam.IDPConfigChangedEvent, bool) {

	changes := make([]idpconfig.IDPConfigChanges, 0)
	oldName := ""
	if wm.Name != name {
		oldName = wm.Name
		changes = append(changes, idpconfig.ChangeName(name))
	}
	if stylingType.Valid() && wm.StylingType != stylingType {
		changes = append(changes, idpconfig.ChangeStyleType(stylingType))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changeEvent, err := iam.NewIDPConfigChangedEvent(ctx, aggregate, configID, oldName, changes)
	if err != nil {
		return nil, false
	}
	return changeEvent, true
}
