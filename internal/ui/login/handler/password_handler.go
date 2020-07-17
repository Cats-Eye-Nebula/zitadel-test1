package handler

import (
	"net/http"

	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	tmplPassword = "password"
)

type passwordFormData struct {
	Password string `schema:"password"`
}

func (l *Login) renderPassword(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	data := l.getUserData(r, authReq, "Password", errType, errMessage)
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplPassword], data, nil)
}

func (l *Login) handlePasswordCheck(w http.ResponseWriter, r *http.Request) {
	data := new(passwordFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	err = l.authRepo.VerifyPassword(setContext(r.Context(), authReq.UserOrgID), authReq.ID, authReq.UserID, data.Password, model.BrowserInfoFromRequest(r))
	if err != nil {
		l.renderPassword(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
