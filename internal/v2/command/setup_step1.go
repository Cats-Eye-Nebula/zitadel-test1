package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/models"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

const (
	OIDCResponseTypeCode           = "CODE"
	OIDCResponseTypeIDToken        = "ID_TOKEN"
	OIDCResponseTypeToken          = "ID_TOKEN TOKEN"
	OIDCGrantTypeAuthorizationCode = "AUTHORIZATION_CODE"
	OIDCGrantTypeImplicit          = "IMPLICIT"
	OIDCGrantTypeRefreshToken      = "REFRESH_TOKEN"
	OIDCApplicationTypeNative      = "NATIVE"
	OIDCApplicationTypeUserAgent   = "USER_AGENT"
	OIDCApplicationTypeWeb         = "WEB"
	OIDCAuthMethodTypeNone         = "NONE"
	OIDCAuthMethodTypeBasic        = "BASIC"
	OIDCAuthMethodTypePost         = "POST"
)

type Step1 struct {
	GlobalOrg          string
	IAMProject         string
	DefaultLoginPolicy LoginPolicy
	Orgs               []Org
}

func (s *Step1) Step() domain.Step {
	return domain.Step1
}

func (s *Step1) execute(ctx context.Context, commandSide *CommandSide) error {
	return commandSide.SetupStep1(ctx, s)
}

type LoginPolicy struct {
	AllowRegister         bool
	AllowUsernamePassword bool
	AllowExternalIdp      bool
}

type User struct {
	FirstName string
	LastName  string
	UserName  string
	Email     string
	Password  string
}

type Org struct {
	Name         string
	Domain       string
	OrgIamPolicy bool
	Owner        User
	Projects     []Project
}

type Project struct {
	Name     string
	Users    []User
	Members  []string
	OIDCApps []OIDCApp
}

type OIDCApp struct {
	Name                   string
	RedirectUris           []string
	ResponseTypes          []string
	GrantTypes             []string
	ApplicationType        string
	AuthMethodType         string
	PostLogoutRedirectUris []string
	DevMode                bool
}

func (r *CommandSide) SetupStep1(ctx context.Context, step1 *Step1) error {
	var events []eventstore.EventPusher
	iamWriteModel := NewIAMWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&iamWriteModel.WriteModel)
	//create default login policy
	loginPolicyEvent, err := r.addDefaultLoginPolicy(ctx, iamAgg, NewIAMLoginPolicyWriteModel(),
		&domain.LoginPolicy{
			AllowUsernamePassword: step1.DefaultLoginPolicy.AllowUsernamePassword,
			AllowRegister:         step1.DefaultLoginPolicy.AllowRegister,
			AllowExternalIDP:      step1.DefaultLoginPolicy.AllowExternalIdp,
		})
	if err != nil {
		return err
	}
	events = append(events, loginPolicyEvent)
	logging.Log("SETUP-sd2hj").Info("default login policy set up")
	//create orgs
	for _, organisation := range step1.Orgs {
		orgAgg, humanWriteModel, _, setUpOrgEvents, err := r.setUpOrg(ctx,
			&domain.Org{
				Name:    organisation.Name,
				Domains: []*domain.OrgDomain{{Domain: organisation.Domain}},
			},
			&domain.Human{
				Username: organisation.Owner.UserName,
				Profile: &domain.Profile{
					FirstName: organisation.Owner.FirstName,
					LastName:  organisation.Owner.LastName,
				},
				Password: &domain.Password{
					SecretString: organisation.Owner.Password,
				},
				Email: &domain.Email{
					EmailAddress:    organisation.Owner.Email,
					IsEmailVerified: true,
				},
			})
		if err != nil {
			return err
		}
		events = append(events, setUpOrgEvents...)
		logging.LogWithFields("SETUP-Gdsfg", "id", orgAgg.ID, "name", organisation.Name).Info("org set up")

		if organisation.OrgIamPolicy {
			orgIAMPolicyEvent, err := r.addOrgIAMPolicy(ctx, orgAgg, NewORGOrgIAMPolicyWriteModel(orgAgg.ID), &domain.OrgIAMPolicy{UserLoginMustBeDomain: false})
			if err != nil {
				return err
			}
			events = append(events, orgIAMPolicyEvent)
		}
		if organisation.Name == step1.GlobalOrg {
			globalOrgEvent, err := r.setGlobalOrg(ctx, iamAgg, iamWriteModel, orgAgg.ID)
			if err != nil {
				return err
			}
			events = append(events, globalOrgEvent)
			logging.Log("SETUP-BDn52").Info("global org set")
		}
		//projects
		for _, proj := range organisation.Projects {
			project := &domain.Project{Name: proj.Name}
			projectEvents, projectWriteModel, err := r.addProject(ctx, project, orgAgg.ID, humanWriteModel.AggregateID)
			if err != nil {
				return err
			}
			events = append(events, projectEvents...)
			if project.Name == step1.IAMProject {
				iamProjectEvent, err := r.setIAMProject(ctx, iamAgg, iamWriteModel, projectWriteModel.AggregateID)
				if err != nil {
					return err
				}
				events = append(events, iamProjectEvent)
				logging.Log("SETUP-Bdfs1").Info("IAM project set")
				iamEvent, err := r.addIAMMember(ctx, iamAgg, NewIAMMemberWriteModel(humanWriteModel.AggregateID), domain.NewMember(iamAgg.ID, humanWriteModel.AggregateID, domain.RoleIAMOwner))
				if err != nil {
					return err
				}
				events = append(events, iamEvent)
				logging.Log("SETUP-BSf2h").Info("IAM owner set")
			}
			//create applications
			for _, app := range proj.OIDCApps {
				applicationEvents, err := setUpApplication(ctx, r, projectWriteModel, project, app, orgAgg.ID)
				if err != nil {
					return err
				}
				events = append(events, applicationEvents...)
			}
		}
	}

	events = append(events, iam_repo.NewSetupStepDoneEvent(ctx, iamAgg, domain.Step1))

	_, err = r.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-Gr2hh", "Setup Step1 failed")
	}
	return nil
}

func setUpApplication(ctx context.Context, r *CommandSide, projectWriteModel *ProjectWriteModel, project *domain.Project, oidcApp OIDCApp, resourceOwner string) ([]eventstore.EventPusher, error) {
	app := &domain.OIDCApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectWriteModel.AggregateID,
		},
		AppName:         oidcApp.Name,
		RedirectUris:    oidcApp.RedirectUris,
		ResponseTypes:   getOIDCResponseTypes(oidcApp.ResponseTypes),
		GrantTypes:      getOIDCGrantTypes(oidcApp.GrantTypes),
		ApplicationType: getOIDCApplicationType(oidcApp.ApplicationType),
		AuthMethodType:  getOIDCAuthMethod(oidcApp.AuthMethodType),
		DevMode:         oidcApp.DevMode,
	}

	projectAgg := ProjectAggregateFromWriteModel(&projectWriteModel.WriteModel)
	events, _, err := r.addOIDCApplication(ctx, projectAgg, project, app, resourceOwner)
	if err != nil {
		return nil, err
	}
	logging.LogWithFields("SETUP-Edgw4", "name", app.AppName, "clientID", app.ClientID).Info("application set up")
	return events, nil
}

func getOIDCResponseTypes(responseTypes []string) []domain.OIDCResponseType {
	types := make([]domain.OIDCResponseType, len(responseTypes))
	for i, t := range responseTypes {
		types[i] = getOIDCResponseType(t)
	}
	return types
}

func getOIDCResponseType(responseType string) domain.OIDCResponseType {
	switch responseType {
	case OIDCResponseTypeCode:
		return domain.OIDCResponseTypeCode
	case OIDCResponseTypeIDToken:
		return domain.OIDCResponseTypeIDToken
	case OIDCResponseTypeToken:
		return domain.OIDCResponseTypeIDTokenToken
	}
	return domain.OIDCResponseTypeCode
}

func getOIDCGrantTypes(grantTypes []string) []domain.OIDCGrantType {
	types := make([]domain.OIDCGrantType, len(grantTypes))
	for i, t := range grantTypes {
		types[i] = getOIDCGrantType(t)
	}
	return types
}

func getOIDCGrantType(grantTypes string) domain.OIDCGrantType {
	switch grantTypes {
	case OIDCGrantTypeAuthorizationCode:
		return domain.OIDCGrantTypeAuthorizationCode
	case OIDCGrantTypeImplicit:
		return domain.OIDCGrantTypeImplicit
	case OIDCGrantTypeRefreshToken:
		return domain.OIDCGrantTypeRefreshToken
	}
	return domain.OIDCGrantTypeAuthorizationCode
}

func getOIDCApplicationType(appType string) domain.OIDCApplicationType {
	switch appType {
	case OIDCApplicationTypeNative:
		return domain.OIDCApplicationTypeNative
	case OIDCApplicationTypeUserAgent:
		return domain.OIDCApplicationTypeUserAgent
	case OIDCApplicationTypeWeb:
		return domain.OIDCApplicationTypeWeb
	}
	return domain.OIDCApplicationTypeWeb
}

func getOIDCAuthMethod(authMethod string) domain.OIDCAuthMethodType {
	switch authMethod {
	case OIDCAuthMethodTypeNone:
		return domain.OIDCAuthMethodTypeNone
	case OIDCAuthMethodTypeBasic:
		return domain.OIDCAuthMethodTypeBasic
	case OIDCAuthMethodTypePost:
		return domain.OIDCAuthMethodTypePost
	}
	return domain.OIDCAuthMethodTypeBasic
}
