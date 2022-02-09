package command

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMCustomMessageTextReadModel struct {
	CustomMessageTextReadModel
}

func NewIAMCustomMessageTextWriteModel(messageTextType string, lang language.Tag) *IAMCustomMessageTextReadModel {
	return &IAMCustomMessageTextReadModel{
		CustomMessageTextReadModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
			MessageTextType: messageTextType,
			Language:        lang,
		},
	}
}

func (wm *IAMCustomMessageTextReadModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.CustomTextSetEvent:
			wm.CustomMessageTextReadModel.AppendEvents(&e.CustomTextSetEvent)
		case *iam.CustomTextRemovedEvent:
			wm.CustomMessageTextReadModel.AppendEvents(&e.CustomTextRemovedEvent)
		case *iam.CustomTextTemplateRemovedEvent:
			wm.CustomMessageTextReadModel.AppendEvents(&e.CustomTextTemplateRemovedEvent)
		}
	}
}

func (wm *IAMCustomMessageTextReadModel) Reduce() error {
	return wm.CustomMessageTextReadModel.Reduce()
}

func (wm *IAMCustomMessageTextReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.CustomMessageTextReadModel.AggregateID).
		EventTypes(iam.CustomTextSetEventType, iam.CustomTextRemovedEventType, iam.CustomTextTemplateRemovedEventType).
		Builder()
}
