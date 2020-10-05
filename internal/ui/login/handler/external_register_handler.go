package handler

import (
	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/rp"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
	caos_errors "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"net/http"
	"strings"
)

func (l *Login) handleExternalRegister(w http.ResponseWriter, r *http.Request) {
	data := new(externalIDPData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if authReq == nil {
		http.Redirect(w, r, l.zitadelURL, http.StatusFound)
		return
	}
	idpConfig, err := l.getIDPConfigByID(r, data.IDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.authRepo.SelectExternalIDP(r.Context(), authReq.ID, idpConfig.IDPConfigID, userAgentID)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	if !idpConfig.IsOIDC {
		l.renderError(w, r, authReq, caos_errors.ThrowInternal(nil, "LOGIN-Rio9s", "Errors.User.ExternalIDP.IDPTypeNotImplemented"))
		return
	}
	l.handleOIDCAuthorize(w, r, authReq, idpConfig, EndpointExternalRegisterCallback)
}

func (l *Login) handleExternalRegisterCallback(w http.ResponseWriter, r *http.Request) {
	data := new(externalIDPCallbackData)
	err := l.getParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), data.State, userAgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	provider := l.getRPConfig(w, r, authReq, idpConfig, EndpointExternalRegisterCallback)
	tokens, err := rp.CodeExchange(r.Context(), data.Code, provider)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	l.handleExternalUserRegister(w, r, authReq, idpConfig, userAgentID, tokens)
}

func (l *Login) handleExternalUserRegister(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, idpConfig *iam_model.IDPConfigView, userAgentID string, tokens *oidc.Tokens) {
	iam, err := l.authRepo.GetIAM(r.Context())
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	resourceOwner := iam.GlobalOrgID
	member := &org_model.OrgMember{
		ObjectRoot: models.ObjectRoot{AggregateID: iam.GlobalOrgID},
		Roles:      []string{orgProjectCreatorRole},
	}

	if authReq.GetScopeOrgPrimaryDomain() != "" {
		primaryDomain := authReq.GetScopeOrgPrimaryDomain()
		org, err := l.authRepo.GetOrgByPrimaryDomain(primaryDomain)
		if err != nil {
			l.renderRegisterOption(w, r, authReq, err)
			return
		}
		if org.ID != iam.GlobalOrgID {
			member = nil
			resourceOwner = org.ID
		}
	}
	orgIamPolicy, err := l.getOrgIamPolicy(r, resourceOwner)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	user, externalIDP := l.mapTokenToLoginUserAndExternalIDP(orgIamPolicy, tokens, idpConfig)
	_, err = l.authRepo.RegisterExternalUser(setContext(r.Context(), resourceOwner), user, externalIDP, member, resourceOwner)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) mapTokenToLoginUserAndExternalIDP(orgIamPolicy *iam_model.OrgIAMPolicyView, tokens *oidc.Tokens, idpConfig *iam_model.IDPConfigView) (*usr_model.User, *usr_model.ExternalIDP) {
	username := tokens.IDTokenClaims.PreferredUsername
	switch idpConfig.OIDCUsernameMapping {
	case iam_model.OIDCMappingFieldEmail:
		if tokens.IDTokenClaims.EmailVerified && tokens.IDTokenClaims.Email != "" {
			username = tokens.IDTokenClaims.Email
		}
	}

	if orgIamPolicy.UserLoginMustBeDomain {
		splittedUsername := strings.Split(username, "@")
		if len(splittedUsername) > 1 {
			username = splittedUsername[0]
		}
	}

	user := &usr_model.User{
		UserName: username,
		Human: &usr_model.Human{
			Profile: &usr_model.Profile{
				FirstName:         tokens.IDTokenClaims.GivenName,
				LastName:          tokens.IDTokenClaims.FamilyName,
				PreferredLanguage: tokens.IDTokenClaims.Locale,
				NickName:          tokens.IDTokenClaims.Nickname,
			},
			Email: &usr_model.Email{
				EmailAddress:    tokens.IDTokenClaims.Email,
				IsEmailVerified: tokens.IDTokenClaims.EmailVerified,
			},
		},
	}
	if tokens.IDTokenClaims.PhoneNumber != "" {
		user.Phone = &usr_model.Phone{
			PhoneNumber:     tokens.IDTokenClaims.PhoneNumber,
			IsPhoneVerified: tokens.IDTokenClaims.PhoneNumberVerified,
		}
	}

	displayName := tokens.IDTokenClaims.PreferredUsername
	switch idpConfig.OIDCIDPDisplayNameMapping {
	case iam_model.OIDCMappingFieldEmail:
		if tokens.IDTokenClaims.EmailVerified && tokens.IDTokenClaims.Email != "" {
			displayName = tokens.IDTokenClaims.Email
		}
	}

	externalIDP := &usr_model.ExternalIDP{
		IDPConfigID: idpConfig.IDPConfigID,
		UserID:      tokens.IDTokenClaims.Subject,
		DisplayName: displayName,
	}
	return user, externalIDP
}
