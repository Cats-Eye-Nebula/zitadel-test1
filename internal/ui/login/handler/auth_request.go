package handler

import (
	"net/http"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	queryAuthRequestID = "authRequestID"
)

func (l *Login) getAuthRequest(r *http.Request) (*model.AuthRequest, error) {
	authRequestID := r.FormValue(queryAuthRequestID)
	if authRequestID == "" {
		return nil, nil
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	return l.authRepo.AuthRequestByID(r.Context(), authRequestID, userAgentID)
}

func (l *Login) getAuthRequestAndParseData(r *http.Request, data interface{}) (*model.AuthRequest, error) {
	authReq, err := l.getAuthRequest(r)
	if err != nil {
		return nil, err
	}
	err = l.parser.Parse(r, data)
	return authReq, err
}

func (l *Login) getParseData(r *http.Request, data interface{}) error {
	return l.parser.Parse(r, data)
}
