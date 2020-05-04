package model

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"time"
)

type Phone struct {
	es_models.ObjectRoot

	PhoneNumber     string `json:"phone,omitempty"`
	IsPhoneVerified bool   `json:"-"`
}

type PhoneCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue `json:"code,omitempty"`
	Expiry time.Duration       `json:"expiry,omitempty"`
}

func (p *Phone) Changes(changed *Phone) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.PhoneNumber != "" && p.PhoneNumber != changed.PhoneNumber {
		changes["phone"] = changed.PhoneNumber
	}
	return changes
}

func PhoneFromModel(phone *model.Phone) *Phone {
	return &Phone{
		ObjectRoot:      phone.ObjectRoot,
		PhoneNumber:     phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func PhoneToModel(phone *Phone) *model.Phone {
	return &model.Phone{
		ObjectRoot:      phone.ObjectRoot,
		PhoneNumber:     phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func PhoneCodeToModel(code *PhoneCode) *model.PhoneCode {
	return &model.PhoneCode{
		ObjectRoot: code.ObjectRoot,
		Expiry:     code.Expiry,
		Code:       code.Code,
	}
}

func (u *User) appendUserPhoneChangedEvent(event *es_models.Event) error {
	u.Phone = new(Phone)
	u.IsPhoneVerified = false
	return u.Phone.setData(event)
}

func (u *User) appendUserPhoneCodeAddedEvent(event *es_models.Event) error {
	u.PhoneCode = new(PhoneCode)
	return u.PhoneCode.setData(event)
}

func (u *User) appendUserPhoneVerifiedEvent() {
	u.IsPhoneVerified = true
}

func (p *Phone) setData(event *es_models.Event) error {
	p.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, p); err != nil {
		logging.Log("EVEN-lco9s").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lre56", "could not unmarshal event")
	}
	return nil
}

func (a *PhoneCode) setData(event *es_models.Event) error {
	a.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-sk8ws").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-7hdj3", "could not unmarshal event")
	}
	return nil
}
