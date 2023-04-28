package model

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type Password struct {
	es_models.ObjectRoot

	Secret         *crypto.CryptoValue `json:"secret,omitempty"`
	ChangeRequired bool                `json:"changeRequired,omitempty"`
}

type PasswordCode struct {
	es_models.ObjectRoot

	Code             *crypto.CryptoValue `json:"code,omitempty"`
	Expiry           time.Duration       `json:"expiry,omitempty"`
	NotificationType int32               `json:"notificationType,omitempty"`
}

type PasswordChange struct {
	Password
	UserAgentID string `json:"userAgentID,omitempty"`
}

func (u *Human) appendUserPasswordChangedEvent(event eventstore.Event) error {
	u.Password = new(Password)
	err := u.Password.setData(event)
	if err != nil {
		return err
	}
	u.Password.ObjectRoot.CreationDate = event.CreationDate()
	return nil
}

func (u *Human) appendPasswordSetRequestedEvent(event eventstore.Event) error {
	u.PasswordCode = new(PasswordCode)
	return u.PasswordCode.SetData(event)
}

func (pw *Password) setData(event eventstore.Event) error {
	pw.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.DataAsBytes(), pw); err != nil {
		logging.Log("EVEN-dks93").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-sl9xlo2rsw", "could not unmarshal event")
	}
	return nil
}

func (c *PasswordCode) SetData(event eventstore.Event) error {
	c.ObjectRoot.AppendEvent(event)
	c.CreationDate = event.CreationDate()
	if err := json.Unmarshal(event.DataAsBytes(), c); err != nil {
		logging.Log("EVEN-lo0y2").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-q21dr", "could not unmarshal event")
	}
	return nil
}

func (pw *PasswordChange) SetData(event eventstore.Event) error {
	if err := json.Unmarshal(event.DataAsBytes(), pw); err != nil {
		logging.Log("EVEN-ADs31").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-BDd32", "could not unmarshal event")
	}
	pw.ObjectRoot.AppendEvent(event)
	return nil
}
