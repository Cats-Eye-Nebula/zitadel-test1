package login

import (
	"net/http"

	"github.com/caos/zitadel/internal/api/authz"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/domain"
)

const (
	tmplLinkUsersDone = "linkusersdone"
)

func (l *Login) linkUsers(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	instanceID := authz.GetInstance(r.Context()).InstanceID()
	err = l.authRepo.LinkExternalUsers(setContext(r.Context(), authReq.UserOrgID), authReq.ID, userAgentID, instanceID, domain.BrowserInfoFromRequest(r))
	l.renderLinkUsersDone(w, r, authReq, err)
}

func (l *Login) renderLinkUsersDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errType, errMessage string
	data := l.getUserData(r, authReq, "Linking Users Done", errType, errMessage)
	l.renderer.RenderTemplate(w, r, l.getTranslator(authReq), l.renderer.Templates[tmplLinkUsersDone], data, nil)
}
