package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/user"
)

type HumanPhoneWriteModel struct {
	eventstore.WriteModel

	Phone           string
	IsPhoneVerified bool

	Code             *crypto.CryptoValue
	CodeCreationDate time.Time
	CodeExpiry       time.Duration

	State domain.PhoneState
}

func NewHumanPhoneWriteModel(userID, resourceOwner string) *HumanPhoneWriteModel {
	return &HumanPhoneWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanPhoneWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			if e.PhoneNumber != "" {
				wm.Phone = e.PhoneNumber
			}
			wm.State = domain.PhoneStateActive
		case *user.HumanRegisteredEvent:
			if e.PhoneNumber != "" {
				wm.Phone = e.PhoneNumber
				wm.State = domain.PhoneStateActive
			}
		case *user.HumanPhoneChangedEvent:
			wm.Phone = e.PhoneNumber
			wm.IsPhoneVerified = false
			wm.State = domain.PhoneStateActive
			wm.Code = nil
		case *user.HumanPhoneVerifiedEvent:
			wm.IsPhoneVerified = true
			wm.Code = nil
		case *user.HumanPhoneCodeAddedEvent:
			wm.Code = e.Code
			wm.CodeCreationDate = e.CreationDate()
			wm.CodeExpiry = e.Expiry
		case *user.HumanPhoneRemovedEvent:
			wm.State = domain.PhoneStateRemoved
		case *user.UserRemovedEvent:
			wm.State = domain.PhoneStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanPhoneWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(user.HumanAddedType,
			user.HumanRegisteredType,
			user.HumanPhoneChangedType,
			user.HumanPhoneVerifiedType,
			user.HumanPhoneCodeAddedType,
			user.HumanPhoneRemovedType,
			user.UserRemovedType)
}

func (wm *HumanPhoneWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	phone string,
) (*user.HumanPhoneChangedEvent, bool) {
	hasChanged := false
	changedEvent := user.NewHumanPhoneChangedEvent(ctx, aggregate)
	if wm.Phone != phone {
		hasChanged = true
		changedEvent.PhoneNumber = phone
	}
	return changedEvent, hasChanged
}
