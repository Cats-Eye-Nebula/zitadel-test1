package login

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/sirupsen/logrus"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
)

const (
	tmplDeviceAuthUserCode = "device-usercode"
	tmplDeviceAuthConfirm  = "device-confirm"
)

func (l *Login) renderDeviceAuthUserCode(w io.Writer, r *http.Request, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}

	data := l.getBaseData(r, &domain.AuthRequest{}, "DeviceAuth.Title", "DeviceAuth.Description", errID, errMessage)
	err = l.renderer.Templates[tmplDeviceAuthUserCode].Execute(w, data)
	if err != nil {
		logrus.Error(err)
	}
}

func (l *Login) renderDeviceAuthConfirm(w http.ResponseWriter, username, clientID string, scopes []string) {
	data := &struct {
		Username string
		ClientID string
		Scopes   []string
	}{
		Username: username,
		ClientID: clientID,
		Scopes:   scopes,
	}

	err := l.renderer.Templates[tmplDeviceAuthConfirm].Execute(w, data)
	if err != nil {
		logrus.Error(err)
	}
}

// handleDeviceUserCode serves the Device Authorization user code submission form.
// The "user_code" may be submitted by URL (GET) or form (POST).
// When a "user_code" is received and found through query,
// handleDeviceAuthUserCode will create a new AuthRequest in the repository.
// The user is then redirected to the /login endpoint to complete authentication.
//
// The agent ID from the context is set to the authentication request
// to ensure the complete login flow is completed from the same browser.
func (l *Login) handleDeviceAuthUserCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		l.renderDeviceAuthUserCode(w, r, err)
		return
	}
	userCode := r.Form.Get("user_code")
	if userCode == "" {
		if prompt, _ := url.QueryUnescape(r.Form.Get("prompt")); prompt != "" {
			err = errors.New(prompt)
		}
		l.renderDeviceAuthUserCode(w, r, err)
		return
	}
	deviceAuth, err := l.query.DeviceAuthByUserCode(ctx, userCode)
	if err != nil {
		l.renderDeviceAuthUserCode(w, r, err)
		return
	}
	userAgentID, ok := middleware.UserAgentIDFromCtx(ctx)
	if !ok {
		l.renderDeviceAuthUserCode(w, r, errors.New("internal error: agent ID missing"))
		return
	}
	authRequest, err := l.authRepo.CreateAuthRequest(ctx, &domain.AuthRequest{
		CreationDate:  time.Now(),
		AgentID:       userAgentID,
		ApplicationID: deviceAuth.ClientID,
		InstanceID:    authz.GetInstance(ctx).InstanceID(),
		Request: &domain.AuthRequestDevice{
			ID:         deviceAuth.AggregateID,
			DeviceCode: deviceAuth.DeviceCode,
			UserCode:   deviceAuth.UserCode,
			Scopes:     deviceAuth.Scopes,
		},
	})
	if err != nil {
		l.renderDeviceAuthUserCode(w, r, err)
		return
	}

	http.Redirect(w, r, l.renderer.pathPrefix+EndpointLogin+"?authRequestID="+authRequest.ID, http.StatusFound)
}

// redirectDeviceAuthStart redirects the user to the start point of
// the device authorization flow. A prompt can be set to inform the user
// of the reason why they are redirected back.
func (l *Login) redirectDeviceAuthStart(w http.ResponseWriter, r *http.Request, prompt string) {
	values := make(url.Values)
	values.Set("prompt", url.QueryEscape(prompt))

	url := url.URL{
		Path:     l.renderer.pathPrefix + EndpointDeviceAuth,
		RawQuery: values.Encode(),
	}
	http.Redirect(w, r, url.String(), http.StatusSeeOther)
}

type deviceConfirmRequest struct {
	AuthRequestID string `schema:"authRequestID"`
	Action        string `schema:"action"`
}

// handleDeviceAuthConfirm is the handler where the user is redirected after login.
// The authRequest is checked if the login was indeed completed.
// When the action of "allowed" or "denied", the device authorization is updated accordingly.
// Else the user is presented with a page where they can choose / submit either action.
func (l *Login) handleDeviceAuthConfirm(w http.ResponseWriter, r *http.Request) {
	req := new(deviceConfirmRequest)
	if err := l.getParseData(r, req); err != nil {
		l.redirectDeviceAuthStart(w, r, err.Error())
		return

	}
	agentID, ok := middleware.UserAgentIDFromCtx(r.Context())
	if !ok {
		l.redirectDeviceAuthStart(w, r, "internal error: agent ID missing")
		return
	}
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), req.AuthRequestID, agentID)
	if err != nil {
		l.redirectDeviceAuthStart(w, r, err.Error())
		return
	}
	if !authReq.Done() {
		l.redirectDeviceAuthStart(w, r, "authentication not completed")
		return
	}
	authDev, ok := authReq.Request.(*domain.AuthRequestDevice)
	if !ok {
		l.redirectDeviceAuthStart(w, r, fmt.Sprintf("wrong auth request type: %T", authReq.Request))
		return
	}

	switch req.Action {
	case "allowed":
		_, err = l.command.ApproveDeviceAuth(r.Context(), authDev.ID, authReq.UserID)
	case "denied":
		_, err = l.command.DenyDeviceAuth(r.Context(), authDev.ID)
	default:
		l.renderDeviceAuthConfirm(w, authReq.UserName, authReq.ApplicationID, authDev.Scopes)
		return
	}
	if err != nil {
		l.redirectDeviceAuthStart(w, r, err.Error())
		return
	}

	fmt.Fprintf(w, "Device authorization %s. You can now return to the device", req.Action)
}

func (l *Login) deviceAuthCallbackURL(authRequestID string) string {
	return l.renderer.pathPrefix + EndpointDeviceAuthConfirm + "?authRequestID=" + authRequestID
}

// RedirectDeviceAuthToPrefix allows users to use https://domain.com/device without the /ui/login prefix
// and redirects them to the prefixed endpoint.
// https://www.rfc-editor.org/rfc/rfc8628#section-3.2 recommends the URL to be as short as possible.
func RedirectDeviceAuthToPrefix(w http.ResponseWriter, r *http.Request) {
	target := gu.PtrCopy(r.URL)
	target.Path = HandlerPrefix + EndpointDeviceAuth
	http.Redirect(w, r, target.String(), http.StatusFound)
}
