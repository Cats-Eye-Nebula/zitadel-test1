package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"net/http"
)

const (
	tmplRegisterOption = "registeroption"
)

type registerOptionFormData struct {
	UsernamePassword bool `schema:"usernamepassword"`
}

type registerOptionData struct {
	baseData
}

func (l *Login) handleRegisterOption(w http.ResponseWriter, r *http.Request) {
	data := new(registerOptionFormData)
	authRequest, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authRequest, err)
		return
	}
	l.renderRegisterOption(w, r, authRequest, nil)
}

func (l *Login) renderRegisterOption(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	data := registerOptionData{
		baseData: l.getBaseData(r, authReq, "RegisterOption", errType, errMessage),
	}
	funcs := map[string]interface{}{
		"hasExternalLogin": func() bool {
			return authReq.LoginPolicy.AllowExternalIDP && authReq.AllowedExternalIDPs != nil && len(authReq.AllowedExternalIDPs) > 0
		},
		"idpProviderClass": func(stylingType iam_model.IDPStylingType) string {
			return stylingType.GetCSSClass()
		},
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplRegisterOption], data, funcs)
}

func (l *Login) handleRegisterOptionCheck(w http.ResponseWriter, r *http.Request) {
	data := new(registerOptionFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if data.UsernamePassword {
		l.handleRegister(w, r)
		return
	}
	l.handleRegisterOption(w, r)
}
