package model

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
)

type OTP struct {
	es_models.ObjectRoot

	Secret *crypto.CryptoValue `json:"otpSecret,omitempty"`
	State  int32               `json:"-"`
}

func OTPFromModel(otp *model.OTP) *OTP {
	return &OTP{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  otp.ObjectRoot.AggregateID,
			Sequence:     otp.Sequence,
			ChangeDate:   otp.ChangeDate,
			CreationDate: otp.CreationDate,
		},
		Secret: otp.Secret,
		State:  int32(otp.State),
	}
}

func OTPToModel(otp *OTP) *model.OTP {
	return &model.OTP{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  otp.ObjectRoot.AggregateID,
			Sequence:     otp.Sequence,
			ChangeDate:   otp.ChangeDate,
			CreationDate: otp.CreationDate,
		},
		Secret: otp.Secret,
		State:  model.MfaState(otp.State),
	}
}

func (u *User) appendOtpAddedEvent(event *es_models.Event) error {
	u.OTP = new(OTP)
	u.OTP.setData(event)
	u.OTP.State = int32(model.MFASTATE_NOTREADY)
	return nil
}

func (u *User) appendOtpVerifiedEvent() error {
	u.OTP.State = int32(model.MFASTATE_READY)
	return nil
}

func (u *User) appendOtpRemovedEvent() error {
	u.OTP = nil
	return nil
}

func (o *OTP) setData(event *es_models.Event) error {
	o.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, o); err != nil {
		logging.Log("EVEN-d9soe").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
