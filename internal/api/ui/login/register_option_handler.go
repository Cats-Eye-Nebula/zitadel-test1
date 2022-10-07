package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
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

func (l *Login) renderRegisterOption(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	allowed := registrationAllowed(authReq)
	externalAllowed := externalRegistrationAllowed(authReq)
	if err == nil {
		// if only external allowed with a single idp then use that
		if !allowed && externalAllowed && len(authReq.AllowedExternalIDPs) == 1 {
			l.handleExternalRegisterByConfigID(w, r, authReq, authReq.AllowedExternalIDPs[0].IDPConfigID)
			return
		}
		// if only direct registration is allowed, show the form
		if allowed && !externalAllowed {
			l.renderRegister(w, r, authReq, nil, nil)
			return
		}
	}
	data := registerOptionData{
		baseData: l.getBaseData(r, authReq, "RegisterOption", errID, errMessage),
	}
	funcs := map[string]interface{}{
		"hasRegistration": func() bool {
			return allowed
		},
		"hasExternalLogin": func() bool {
			return externalAllowed
		},
	}
	translator := l.getTranslator(r.Context(), authReq)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplRegisterOption], data, funcs)
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

func registrationAllowed(authReq *domain.AuthRequest) bool {
	return authReq != nil && authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowRegister
}

func externalRegistrationAllowed(authReq *domain.AuthRequest) bool {
	return authReq != nil && authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowExternalIDP && authReq.AllowedExternalIDPs != nil && len(authReq.AllowedExternalIDPs) > 0
}
