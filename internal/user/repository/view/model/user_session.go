package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	req_model "github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	UserSessionKeyUserAgentID   = "user_agent_id"
	UserSessionKeyUserID        = "user_id"
	UserSessionKeyState         = "state"
	UserSessionKeyResourceOwner = "resource_owner"
)

type UserSessionView struct {
	CreationDate                time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate                  time.Time `json:"-" gorm:"column:change_date"`
	ResourceOwner               string    `json:"-" gorm:"column:resource_owner"`
	State                       int32     `json:"-" gorm:"column:state"`
	UserAgentID                 string    `json:"userAgentID" gorm:"column:user_agent_id;primary_key"`
	UserID                      string    `json:"userID" gorm:"column:user_id;primary_key"`
	UserName                    string    `json:"-" gorm:"column:user_name"`
	LoginName                   string    `json:"-" gorm:"column:login_name"`
	DisplayName                 string    `json:"-" gorm:"column:user_display_name"`
	PasswordVerification        time.Time `json:"-" gorm:"column:password_verification"`
	MfaSoftwareVerification     time.Time `json:"-" gorm:"column:mfa_software_verification"`
	MfaSoftwareVerificationType int32     `json:"-" gorm:"column:mfa_software_verification_type"`
	MfaHardwareVerification     time.Time `json:"-" gorm:"column:mfa_hardware_verification"`
	MfaHardwareVerificationType int32     `json:"-" gorm:"column:mfa_hardware_verification_type"`
	Sequence                    uint64    `json:"-" gorm:"column:sequence"`
}

func UserSessionFromEvent(event *models.Event) (*UserSessionView, error) {
	v := new(UserSessionView)
	if err := json.Unmarshal(event.Data, v); err != nil {
		logging.Log("EVEN-lso9e").WithError(err).Error("could not unmarshal event data")
		return nil, caos_errs.ThrowInternal(nil, "MODEL-sd325", "could not unmarshal data")
	}
	return v, nil
}

func UserSessionToModel(userSession *UserSessionView) *model.UserSessionView {
	return &model.UserSessionView{
		ChangeDate:                  userSession.ChangeDate,
		CreationDate:                userSession.CreationDate,
		ResourceOwner:               userSession.ResourceOwner,
		State:                       req_model.UserSessionState(userSession.State),
		UserAgentID:                 userSession.UserAgentID,
		UserID:                      userSession.UserID,
		UserName:                    userSession.UserName,
		LoginName:                   userSession.LoginName,
		DisplayName:                 userSession.DisplayName,
		PasswordVerification:        userSession.PasswordVerification,
		MfaSoftwareVerification:     userSession.MfaSoftwareVerification,
		MfaSoftwareVerificationType: req_model.MfaType(userSession.MfaSoftwareVerificationType),
		MfaHardwareVerification:     userSession.MfaHardwareVerification,
		MfaHardwareVerificationType: req_model.MfaType(userSession.MfaHardwareVerificationType),
		Sequence:                    userSession.Sequence,
	}
}

func UserSessionsToModel(userSessions []*UserSessionView) []*model.UserSessionView {
	result := make([]*model.UserSessionView, len(userSessions))
	for i, s := range userSessions {
		result[i] = UserSessionToModel(s)
	}
	return result
}

func (v *UserSessionView) AppendEvent(event *models.Event) {
	v.Sequence = event.Sequence
	v.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.UserPasswordCheckSucceeded,
		es_model.HumanPasswordCheckSucceeded:
		v.PasswordVerification = event.CreationDate
		v.State = int32(req_model.UserSessionStateActive)
	case es_model.UserPasswordCheckFailed,
		es_model.UserPasswordChanged,
		es_model.HumanPasswordCheckFailed,
		es_model.HumanPasswordChanged:
		v.PasswordVerification = time.Time{}
	case es_model.MFAOTPCheckSucceeded,
		es_model.HumanMFAOTPCheckSucceeded:
		v.MfaSoftwareVerification = event.CreationDate
		v.MfaSoftwareVerificationType = int32(req_model.MfaTypeOTP)
		v.State = int32(req_model.UserSessionStateActive)
	case es_model.MFAOTPCheckFailed,
		es_model.MFAOTPRemoved,
		es_model.HumanMFAOTPCheckFailed,
		es_model.HumanMFAOTPRemoved:
		v.MfaSoftwareVerification = time.Time{}
	case es_model.SignedOut,
		es_model.HumanSignedOut,
		es_model.UserLocked,
		es_model.UserDeactivated:
		v.PasswordVerification = time.Time{}
		v.MfaSoftwareVerification = time.Time{}
		v.State = int32(req_model.UserSessionStateTerminated)
	}
}
