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
	UserKeyUserID        = "id"
	UserKeyUserName      = "user_name"
	UserKeyFirstName     = "first_name"
	UserKeyLastName      = "last_name"
	UserKeyNickName      = "nick_name"
	UserKeyDisplayName   = "display_name"
	UserKeyEmail         = "email"
	UserKeyState         = "user_state"
	UserKeyResourceOwner = "resource_owner"
)

type UserView struct {
	ID                     string    `json:"-" gorm:"column:id;primary_key"`
	CreationDate           time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate             time.Time `json:"-" gorm:"column:change_date"`
	ResourceOwner          string    `json:"-" gorm:"column:resource_owner"`
	State                  int32     `json:"-" gorm:"column:user_state"`
	PasswordSet            bool      `json:"-" gorm:"column:password_set"`
	PasswordChangeRequired bool      `json:"-" gorm:"column:password_change_required"`
	PasswordChanged        time.Time `json:"-" gorm:"column:password_change"`
	LastLogin              time.Time `json:"-" gorm:"column:last_login"`
	UserName               string    `json:"userName" gorm:"column:user_name"`
	FirstName              string    `json:"firstName" gorm:"column:first_name"`
	LastName               string    `json:"lastName" gorm:"column:last_name"`
	NickName               string    `json:"nickName" gorm:"column:nick_name"`
	DisplayName            string    `json:"displayName" gorm:"column:display_name"`
	PreferredLanguage      string    `json:"preferredLanguage" gorm:"column:preferred_language"`
	Gender                 int32     `json:"gender" gorm:"column:gender"`
	Email                  string    `json:"email" gorm:"column:email"`
	IsEmailVerified        bool      `json:"-" gorm:"column:is_email_verified"`
	Phone                  string    `json:"phone" gorm:"column:phone"`
	IsPhoneVerified        bool      `json:"-" gorm:"column:is_phone_verified"`
	Country                string    `json:"country" gorm:"column:country"`
	Locality               string    `json:"locality" gorm:"column:locality"`
	PostalCode             string    `json:"postalCode" gorm:"column:postal_code"`
	Region                 string    `json:"region" gorm:"column:region"`
	StreetAddress          string    `json:"streetAddress" gorm:"column:street_address"`
	OTPState               int32     `json:"-" gorm:"column:otp_state"`
	MfaMaxSetUp            int32     `json:"-" gorm:"column:mfa_max_set_up"`
	MfaInitSkipped         time.Time `json:"-" gorm:"column:mfa_init_skipped"`
	Sequence               uint64    `json:"-" gorm:"column:sequence"`
}

func UserFromModel(user *model.UserView) *UserView {
	return &UserView{
		ID:                     user.ID,
		ChangeDate:             user.ChangeDate,
		CreationDate:           user.CreationDate,
		ResourceOwner:          user.ResourceOwner,
		State:                  int32(user.State),
		PasswordSet:            user.PasswordSet,
		PasswordChangeRequired: user.PasswordChangeRequired,
		PasswordChanged:        user.PasswordChanged,
		LastLogin:              user.LastLogin,
		UserName:               user.UserName,
		FirstName:              user.FirstName,
		LastName:               user.LastName,
		NickName:               user.NickName,
		DisplayName:            user.DisplayName,
		PreferredLanguage:      user.PreferredLanguage,
		Gender:                 int32(user.Gender),
		Email:                  user.Email,
		IsEmailVerified:        user.IsEmailVerified,
		Phone:                  user.Phone,
		IsPhoneVerified:        user.IsPhoneVerified,
		Country:                user.Country,
		Locality:               user.Locality,
		PostalCode:             user.PostalCode,
		Region:                 user.Region,
		StreetAddress:          user.StreetAddress,
		OTPState:               int32(user.OTPState),
		MfaMaxSetUp:            int32(user.MfaMaxSetUp),
		MfaInitSkipped:         user.MfaInitSkipped,
		Sequence:               user.Sequence,
	}
}

func UserToModel(user *UserView) *model.UserView {
	return &model.UserView{
		ID:                     user.ID,
		ChangeDate:             user.ChangeDate,
		CreationDate:           user.CreationDate,
		ResourceOwner:          user.ResourceOwner,
		State:                  model.UserState(user.State),
		PasswordSet:            user.PasswordSet,
		PasswordChangeRequired: user.PasswordChangeRequired,
		PasswordChanged:        user.PasswordChanged,
		LastLogin:              user.LastLogin,
		UserName:               user.UserName,
		FirstName:              user.FirstName,
		LastName:               user.LastName,
		NickName:               user.NickName,
		DisplayName:            user.DisplayName,
		PreferredLanguage:      user.PreferredLanguage,
		Gender:                 model.Gender(user.Gender),
		Email:                  user.Email,
		IsEmailVerified:        user.IsEmailVerified,
		Phone:                  user.Phone,
		IsPhoneVerified:        user.IsPhoneVerified,
		Country:                user.Country,
		Locality:               user.Locality,
		PostalCode:             user.PostalCode,
		Region:                 user.Region,
		StreetAddress:          user.StreetAddress,
		OTPState:               model.MfaState(user.OTPState),
		MfaMaxSetUp:            req_model.MfaLevel(user.MfaMaxSetUp),
		MfaInitSkipped:         user.MfaInitSkipped,
		Sequence:               user.Sequence,
	}
}

func UsersToModel(users []*UserView) []*model.UserView {
	result := make([]*model.UserView, len(users))
	for i, p := range users {
		result[i] = UserToModel(p)
	}
	return result
}

func (u *UserView) AppendEvent(event *models.Event) (err error) {
	u.ChangeDate = event.CreationDate
	u.Sequence = event.Sequence
	switch event.Type {
	case es_model.UserAdded,
		es_model.UserRegistered:
		u.CreationDate = event.CreationDate
		u.setRootData(event)
		err = u.setData(event)
		if err != nil {
			return err
		}
		err = u.setPasswordData(event)
	case es_model.UserPasswordChanged:
		err = u.setPasswordData(event)
	case es_model.UserProfileChanged,
		es_model.UserAddressChanged:
		err = u.setData(event)
	case es_model.UserEmailChanged:
		u.IsEmailVerified = false
		err = u.setData(event)
	case es_model.UserEmailVerified:
		u.IsEmailVerified = true
	case es_model.UserPhoneChanged:
		u.IsPhoneVerified = false
		err = u.setData(event)
	case es_model.UserPhoneVerified:
		u.IsPhoneVerified = true
	case es_model.UserDeactivated:
		u.State = int32(model.USERSTATE_INACTIVE)
	case es_model.UserReactivated,
		es_model.UserUnlocked:
		u.State = int32(model.USERSTATE_ACTIVE)
	case es_model.UserLocked:
		u.State = int32(model.USERSTATE_LOCKED)
	case es_model.MfaOtpAdded:
		u.OTPState = int32(model.MFASTATE_NOTREADY)
	case es_model.MfaOtpVerified:
		u.OTPState = int32(model.MFASTATE_READY)
		u.MfaInitSkipped = time.Time{}
	case es_model.MfaOtpRemoved:
		u.OTPState = int32(model.MFASTATE_UNSPECIFIED)
	case es_model.MfaInitSkipped:
		u.MfaInitSkipped = event.CreationDate
	}
	u.ComputeObject()
	return err
}

func (u *UserView) setRootData(event *models.Event) {
	u.ID = event.AggregateID
	u.ResourceOwner = event.ResourceOwner
}

func (u *UserView) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("MODEL-lso9e").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-8iows", "could not unmarshal data")
	}
	return nil
}

func (u *UserView) setPasswordData(event *models.Event) error {
	password := new(es_model.Password)
	if err := json.Unmarshal(event.Data, password); err != nil {
		logging.Log("MODEL-sdw4r").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-6jhsw", "could not unmarshal data")
	}
	u.PasswordSet = password.Secret != nil
	u.PasswordChangeRequired = password.ChangeRequired
	return nil
}

func (u *UserView) ComputeObject() {
	if u.State == int32(model.USERSTATE_UNSPECIFIED) || u.State == int32(model.USERSTATE_INITIAL) {
		if u.IsEmailVerified {
			u.State = int32(model.USERSTATE_ACTIVE)
		} else {
			u.State = int32(model.USERSTATE_INITIAL)
		}
	}
	if u.OTPState == int32(model.MFASTATE_READY) {
		u.MfaMaxSetUp = int32(req_model.MfaLevelSoftware)
	}
}
