package oidc

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/op"
	"golang.org/x/text/language"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/http"
	authreq_model "github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	user_model "github.com/caos/zitadel/internal/user/model"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
)

const (
	ScopeProjectRolePrefix = "urn:zitadel:iam:org:project:role:"
	ClaimProjectRoles      = "urn:zitadel:iam:org:project:roles"
	ScopeUserMetaData      = "urn:zitadel:iam:user:metadata"
	ClaimUserMetaData      = ScopeUserMetaData
	ScopeResourceOwner     = "urn:zitadel:iam:user:resourceowner"
	ClaimResourceOwner     = ScopeResourceOwner + ":"

	oidcCtx = "oidc"
)

func (o *OPStorage) GetClientByClientID(ctx context.Context, id string) (_ op.Client, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	client, err := o.query.AppByOIDCClientID(ctx, id)
	if err != nil {
		return nil, err
	}
	if client.State != domain.AppStateActive {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-sdaGg", "client is not active")
	}
	projectIDQuery, err := query.NewProjectRoleProjectIDSearchQuery(client.ProjectID)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-mPxqP", "Errors.Internal")
	}
	projectRoles, err := o.query.SearchProjectRoles(context.TODO(), &query.ProjectRoleSearchQueries{Queries: []query.SearchQuery{projectIDQuery}})
	if err != nil {
		return nil, err
	}
	allowedScopes := make([]string, len(projectRoles.ProjectRoles))
	for i, role := range projectRoles.ProjectRoles {
		allowedScopes[i] = ScopeProjectRolePrefix + role.Key
	}
	return ClientFromBusiness(client, o.defaultLoginURL, o.defaultAccessTokenLifetime, o.defaultIdTokenLifetime, allowedScopes)
}

func (o *OPStorage) GetKeyByIDAndUserID(ctx context.Context, keyID, userID string) (_ *jose.JSONWebKey, err error) {
	return o.GetKeyByIDAndIssuer(ctx, keyID, userID)
}

func (o *OPStorage) GetKeyByIDAndIssuer(ctx context.Context, keyID, issuer string) (_ *jose.JSONWebKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	publicKeyData, err := o.query.GetAuthNKeyPublicKeyByIDAndIdentifier(ctx, keyID, issuer)
	if err != nil {
		return nil, err
	}
	publicKey, err := crypto.BytesToPublicKey(publicKeyData)
	if err != nil {
		return nil, err
	}
	return &jose.JSONWebKey{
		KeyID: keyID,
		Use:   "sig",
		Key:   publicKey,
	}, nil
}

func (o *OPStorage) ValidateJWTProfileScopes(ctx context.Context, subject string, scopes []string) ([]string, error) {
	user, err := o.repo.UserByID(ctx, subject)
	if err != nil {
		return nil, err
	}
	for i := len(scopes) - 1; i >= 0; i-- {
		scope := scopes[i]
		if strings.HasPrefix(scope, authreq_model.OrgDomainPrimaryScope) {
			var orgID string
			org, err := o.query.OrgByDomainGlobal(ctx, strings.TrimPrefix(scope, authreq_model.OrgDomainPrimaryScope))
			if err == nil {
				orgID = org.ID
			}
			if orgID != user.ResourceOwner {
				scopes[i] = scopes[len(scopes)-1]
				scopes[len(scopes)-1] = ""
				scopes = scopes[:len(scopes)-1]
			}
		}
	}
	return scopes, nil
}

func (o *OPStorage) AuthorizeClientIDSecret(ctx context.Context, id string, secret string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	ctx = authz.SetCtxData(ctx, authz.CtxData{
		UserID: oidcCtx,
		OrgID:  oidcCtx,
	})
	return o.repo.AuthorizeClientIDSecret(ctx, id, secret)
}

func (o *OPStorage) SetUserinfoFromToken(ctx context.Context, userInfo oidc.UserInfoSetter, tokenID, subject, origin string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	token, err := o.repo.TokenByID(ctx, subject, tokenID)
	if err != nil {
		return errors.ThrowPermissionDenied(nil, "OIDC-Dsfb2", "token is not valid or has expired")
	}
	if token.ApplicationID != "" {
		app, err := o.query.AppByOIDCClientID(ctx, token.ApplicationID)
		if err != nil {
			return err
		}
		if origin != "" && !http.IsOriginAllowed(app.OIDCConfig.AllowedOrigins, origin) {
			return errors.ThrowPermissionDenied(nil, "OIDC-da1f3", "origin is not allowed")
		}
	}
	return o.SetUserinfoFromScopes(ctx, userInfo, token.UserID, token.ApplicationID, token.Scopes)
}

func (o *OPStorage) SetUserinfoFromScopes(ctx context.Context, userInfo oidc.UserInfoSetter, userID, applicationID string, scopes []string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	user, err := o.repo.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	roles := make([]string, 0)
	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeOpenID:
			userInfo.SetSubject(user.ID)
		case oidc.ScopeEmail:
			if user.HumanView == nil {
				continue
			}
			userInfo.SetEmail(user.Email, user.IsEmailVerified)
		case oidc.ScopeProfile:
			userInfo.SetPreferredUsername(user.PreferredLoginName)
			userInfo.SetUpdatedAt(user.ChangeDate)
			if user.HumanView != nil {
				userInfo.SetName(user.DisplayName)
				userInfo.SetFamilyName(user.LastName)
				userInfo.SetGivenName(user.FirstName)
				userInfo.SetNickname(user.NickName)
				userInfo.SetGender(oidc.Gender(getGender(user.Gender)))
				locale, _ := language.Parse(user.PreferredLanguage)
				userInfo.SetLocale(locale)
				userInfo.SetPicture(user.AvatarURL)
			} else {
				userInfo.SetName(user.MachineView.Name)
			}
		case oidc.ScopePhone:
			if user.HumanView == nil {
				continue
			}
			userInfo.SetPhone(user.Phone, user.IsPhoneVerified)
		case oidc.ScopeAddress:
			if user.HumanView == nil {
				continue
			}
			if user.StreetAddress == "" && user.Locality == "" && user.Region == "" && user.PostalCode == "" && user.Country == "" {
				continue
			}
			userInfo.SetAddress(oidc.NewUserInfoAddress(user.StreetAddress, user.Locality, user.Region, user.PostalCode, user.Country, ""))
		case ScopeUserMetaData:
			userMetaData, err := o.assertUserMetaData(ctx, userID)
			if err != nil {
				return err
			}
			if len(userMetaData) > 0 {
				userInfo.AppendClaims(ClaimUserMetaData, userMetaData)
			}
		case ScopeResourceOwner:
			resourceOwnerClaims, err := o.assertUserResourceOwner(ctx, userID)
			if err != nil {
				return err
			}
			for claim, value := range resourceOwnerClaims {
				userInfo.AppendClaims(claim, value)
			}

		default:
			if strings.HasPrefix(scope, ScopeProjectRolePrefix) {
				roles = append(roles, strings.TrimPrefix(scope, ScopeProjectRolePrefix))
			}
			if strings.HasPrefix(scope, authreq_model.OrgDomainPrimaryScope) {
				userInfo.AppendClaims(authreq_model.OrgDomainPrimaryClaim, strings.TrimPrefix(scope, authreq_model.OrgDomainPrimaryScope))
			}
		}
	}
	if len(roles) == 0 || applicationID == "" {
		return nil
	}
	projectRoles, err := o.assertRoles(ctx, userID, applicationID, roles)
	if err != nil {
		return err
	}
	if len(projectRoles) > 0 {
		userInfo.AppendClaims(ClaimProjectRoles, projectRoles)
	}
	return nil
}

func (o *OPStorage) SetIntrospectionFromToken(ctx context.Context, introspection oidc.IntrospectionResponse, tokenID, subject, clientID string) error {
	token, err := o.repo.TokenByID(ctx, subject, tokenID)
	if err != nil {
		return errors.ThrowPermissionDenied(nil, "OIDC-Dsfb2", "token is not valid or has expired")
	}
	projectID, err := o.query.ProjectIDFromClientID(ctx, clientID)
	if err != nil {
		return errors.ThrowPermissionDenied(nil, "OIDC-Adfg5", "client not found")
	}
	for _, aud := range token.Audience {
		if aud == clientID || aud == projectID {
			err := o.SetUserinfoFromScopes(ctx, introspection, token.UserID, clientID, token.Scopes)
			if err != nil {
				return err
			}
			introspection.SetScopes(token.Scopes)
			introspection.SetClientID(token.ApplicationID)
			return nil
		}
	}
	return errors.ThrowPermissionDenied(nil, "OIDC-sdg3G", "token is not valid for this client")
}

func (o *OPStorage) GetPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (claims map[string]interface{}, err error) {
	roles := make([]string, 0)
	for _, scope := range scopes {
		switch scope {
		case ScopeUserMetaData:
			userMetaData, err := o.assertUserMetaData(ctx, userID)
			if err != nil {
				return nil, err
			}
			if len(userMetaData) > 0 {
				claims = appendClaim(claims, ClaimUserMetaData, userMetaData)
			}
		case ScopeResourceOwner:
			resourceOwnerClaims, err := o.assertUserResourceOwner(ctx, userID)
			if err != nil {
				return nil, err
			}
			for claim, value := range resourceOwnerClaims {
				claims = appendClaim(claims, claim, value)
			}
		}
		if strings.HasPrefix(scope, ScopeProjectRolePrefix) {
			roles = append(roles, strings.TrimPrefix(scope, ScopeProjectRolePrefix))
		} else if strings.HasPrefix(scope, authreq_model.OrgDomainPrimaryScope) {
			claims = appendClaim(claims, authreq_model.OrgDomainPrimaryClaim, strings.TrimPrefix(scope, authreq_model.OrgDomainPrimaryScope))
		}
	}
	if len(roles) == 0 || clientID == "" {
		return claims, nil
	}
	projectRoles, err := o.assertRoles(ctx, userID, clientID, roles)
	if err != nil {
		return nil, err
	}
	if len(projectRoles) > 0 {
		claims = appendClaim(claims, ClaimProjectRoles, projectRoles)
	}
	return claims, err
}

func (o *OPStorage) assertRoles(ctx context.Context, userID, applicationID string, requestedRoles []string) (map[string]map[string]string, error) {
	projectID, err := o.query.ProjectIDFromClientID(ctx, applicationID)
	if err != nil {
		return nil, err
	}
	grants, err := o.repo.UserGrantsByProjectAndUserID(projectID, userID)
	if err != nil {
		return nil, err
	}
	projectRoles := make(map[string]map[string]string)
	for _, requestedRole := range requestedRoles {
		for _, grant := range grants {
			checkGrantedRoles(projectRoles, grant, requestedRole)
		}
	}
	return projectRoles, nil
}

func (o *OPStorage) assertUserMetaData(ctx context.Context, userID string) (map[string]string, error) {
	metaData, err := o.repo.SearchUserMetadata(ctx, userID)
	if err != nil {
		return nil, err
	}

	userMetaData := make(map[string]string)
	for _, md := range metaData.Result {
		userMetaData[md.Key] = base64.RawURLEncoding.EncodeToString(md.Value)
	}
	return userMetaData, nil
}

func (o *OPStorage) assertUserResourceOwner(ctx context.Context, userID string) (map[string]string, error) {
	user, err := o.repo.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	resourceOwner, err := o.query.OrgByID(ctx, user.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		ClaimResourceOwner + "id":             resourceOwner.ID,
		ClaimResourceOwner + "name":           resourceOwner.Name,
		ClaimResourceOwner + "primary_domain": resourceOwner.Domain,
	}, nil
}

func checkGrantedRoles(roles map[string]map[string]string, grant *grant_model.UserGrantView, requestedRole string) {
	for _, grantedRole := range grant.RoleKeys {
		if requestedRole == grantedRole {
			appendRole(roles, grantedRole, grant.ResourceOwner, grant.OrgPrimaryDomain)
		}
	}
}

func appendRole(roles map[string]map[string]string, role, orgID, orgPrimaryDomain string) {
	if roles[role] == nil {
		roles[role] = make(map[string]string, 0)
	}
	roles[role][orgID] = orgPrimaryDomain
}

func getGender(gender user_model.Gender) string {
	switch gender {
	case user_model.GenderFemale:
		return "female"
	case user_model.GenderMale:
		return "male"
	case user_model.GenderDiverse:
		return "diverse"
	}
	return ""
}

func appendClaim(claims map[string]interface{}, claim string, value interface{}) map[string]interface{} {
	if claims == nil {
		claims = make(map[string]interface{})
	}
	claims[claim] = value
	return claims
}
