package handler

import (
	"net/http"
	"strconv"

	"github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

const (
	queryInitUserCode     = "code"
	queryInitUserUserID   = "userID"
	queryInitUserPassword = "passwordset"

	tmplInitUser     = "inituser"
	tmplInitUserDone = "inituserdone"
)

type initUserFormData struct {
	Code            string `schema:"code"`
	Password        string `schema:"password"`
	PasswordConfirm string `schema:"passwordconfirm"`
	UserID          string `schema:"userID"`
	PasswordSet     bool   `schema:"passwordSet"`
	Resend          bool   `schema:"resend"`
}

type initUserData struct {
	baseData
	Code        string
	UserID      string
	PasswordSet bool
}

func (l *Login) handleInitUser(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue(queryInitUserUserID)
	code := r.FormValue(queryInitUserCode)
	passwordSet, _ := strconv.ParseBool(r.FormValue(queryInitUserPassword))
	l.renderInitUser(w, r, nil, userID, code, passwordSet, nil)
}

func (l *Login) handleInitUserCheck(w http.ResponseWriter, r *http.Request) {
	data := new(initUserFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}

	if data.Resend {
		l.resendUserInit(w, r, authReq, data.UserID, data.PasswordSet)
		return
	}
	l.checkUserInitCode(w, r, authReq, data, nil)
}

func (l *Login) checkUserInitCode(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, data *initUserFormData, err error) {
	if data.Password != data.PasswordConfirm {
		err := caos_errs.ThrowInvalidArgument(nil, "VIEW-fsdfd", "Errors.User.Password.ConfirmationWrong")
		l.renderInitUser(w, r, authReq, data.UserID, data.Code, data.PasswordSet, err)
		return
	}
	userOrgID := login
	if authReq != nil {
		userOrgID = authReq.UserOrgID
	}
	err = l.authRepo.VerifyInitCode(setContext(r.Context(), userOrgID), data.UserID, data.Code, data.Password)
	if err != nil {
		l.renderInitUser(w, r, nil, data.UserID, "", data.PasswordSet, err)
		return
	}
	l.renderInitUserDone(w, r, nil)
}

func (l *Login) resendUserInit(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, userID string, showPassword bool) {
	userOrgID := login
	if authReq != nil {
		userOrgID = authReq.UserOrgID
	}
	err := l.authRepo.ResendInitVerificationMail(setContext(r.Context(), userOrgID), userID)
	l.renderInitUser(w, r, authReq, userID, "", showPassword, err)
}

func (l *Login) renderInitUser(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, userID, code string, passwordSet bool, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	if authReq != nil {
		userID = authReq.UserID
	}
	data := initUserData{
		baseData:    l.getBaseData(r, nil, "Init User", errType, errMessage),
		UserID:      userID,
		Code:        code,
		PasswordSet: passwordSet,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplInitUser], data, nil)
}

func (l *Login) renderInitUserDone(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest) {
	data := l.getUserData(r, authReq, "User Init Done", "", "")
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplInitUserDone], data, nil)
}
