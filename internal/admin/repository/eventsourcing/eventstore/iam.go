package eventstore

import (
	"context"
	"strings"

	"github.com/caos/logging"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_es "github.com/caos/zitadel/internal/user/repository/eventsourcing"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
)

type IAMRepository struct {
	SearchLimit uint64
	*iam_es.IAMEventstore
	OrgEvents      *org_es.OrgEventstore
	UserEvents     *usr_es.UserEventstore
	View           *admin_view.View
	SystemDefaults systemdefaults.SystemDefaults
	Roles          []string
}

func (repo *IAMRepository) IAMMemberByID(ctx context.Context, orgID, userID string) (*iam_model.IAMMemberView, error) {
	member, err := repo.View.IAMMemberByIDs(orgID, userID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.IAMMemberToModel(member), nil
}

func (repo *IAMRepository) AddIAMMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
	member.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddIAMMember(ctx, member)
}

func (repo *IAMRepository) ChangeIAMMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
	member.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangeIAMMember(ctx, member)
}

func (repo *IAMRepository) RemoveIAMMember(ctx context.Context, userID string) error {
	member := iam_model.NewIAMMember(repo.SystemDefaults.IamID, userID)
	return repo.IAMEventstore.RemoveIAMMember(ctx, member)
}

func (repo *IAMRepository) SearchIAMMembers(ctx context.Context, request *iam_model.IAMMemberSearchRequest) (*iam_model.IAMMemberSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, err := repo.View.GetLatestIAMMemberSequence()
	logging.Log("EVENT-Slkci").OnError(err).Warn("could not read latest iam sequence")
	members, count, err := repo.View.SearchIAMMembers(request)
	if err != nil {
		return nil, err
	}
	result := &iam_model.IAMMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      iam_es_model.IAMMembersToModel(members),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func (repo *IAMRepository) GetIAMMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "IAM") {
			roles = append(roles, roleMap)
		}
	}
	return roles
}

func (repo *IAMRepository) IDPConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IDPConfigView, error) {
	idp, err := repo.View.IDPConfigByID(idpConfigID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.IDPConfigViewToModel(idp), nil
}
func (repo *IAMRepository) AddOIDCIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	idp.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddIDPConfig(ctx, idp)
}

func (repo *IAMRepository) ChangeIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	idp.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangeIDPConfig(ctx, idp)
}

func (repo *IAMRepository) DeactivateIDPConfig(ctx context.Context, idpConfigID string) (*iam_model.IDPConfig, error) {
	return repo.IAMEventstore.DeactivateIDPConfig(ctx, repo.SystemDefaults.IamID, idpConfigID)
}

func (repo *IAMRepository) ReactivateIDPConfig(ctx context.Context, idpConfigID string) (*iam_model.IDPConfig, error) {
	return repo.IAMEventstore.ReactivateIDPConfig(ctx, repo.SystemDefaults.IamID, idpConfigID)
}

func (repo *IAMRepository) RemoveIDPConfig(ctx context.Context, idpConfigID string) error {
	aggregates := make([]*es_models.Aggregate, 0)
	idp := iam_model.NewIDPConfig(repo.SystemDefaults.IamID, idpConfigID)
	_, agg, err := repo.IAMEventstore.PrepareRemoveIDPConfig(ctx, idp)
	if err != nil {
		return err
	}
	aggregates = append(aggregates, agg)

	providers, err := repo.View.IDPProvidersByIdpConfigID(idpConfigID)
	if err != nil {
		return err
	}
	for _, p := range providers {
		if p.AggregateID == repo.SystemDefaults.IamID {
			continue
		}
		provider := &iam_model.IDPProvider{ObjectRoot: es_models.ObjectRoot{AggregateID: p.AggregateID}, IdpConfigID: p.IDPConfigID}
		providerAgg := new(es_models.Aggregate)
		_, providerAgg, err = repo.OrgEvents.PrepareRemoveIDPProviderFromLoginPolicy(ctx, provider, true)
		if err != nil {
			return err
		}
		aggregates = append(aggregates, providerAgg)
	}
	externalIDPs, err := repo.View.ExternalIDPsByIDPConfigID(idpConfigID)
	if err != nil {
		return err
	}
	for _, externalIDP := range externalIDPs {
		idpRemove := &usr_model.ExternalIDP{ObjectRoot: es_models.ObjectRoot{AggregateID: externalIDP.UserID}, IDPConfigID: externalIDP.IDPConfigID, UserID: externalIDP.ExternalUserID}
		idpAgg := make([]*es_models.Aggregate, 0)
		_, idpAgg, err = repo.UserEvents.PrepareRemoveExternalIDP(ctx, idpRemove, true)
		if err != nil {
			return err
		}
		aggregates = append(aggregates, idpAgg...)
	}
	return es_sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, nil, aggregates...)
}

func (repo *IAMRepository) ChangeOidcIDPConfig(ctx context.Context, oidcConfig *iam_model.OIDCIDPConfig) (*iam_model.OIDCIDPConfig, error) {
	oidcConfig.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangeIDPOIDCConfig(ctx, oidcConfig)
}

func (repo *IAMRepository) SearchIDPConfigs(ctx context.Context, request *iam_model.IDPConfigSearchRequest) (*iam_model.IDPConfigSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, err := repo.View.GetLatestIDPConfigSequence()
	logging.Log("EVENT-Dk8si").OnError(err).Warn("could not read latest idp config sequence")
	idps, count, err := repo.View.SearchIDPConfigs(request)
	if err != nil {
		return nil, err
	}
	result := &iam_model.IDPConfigSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      iam_es_model.IdpConfigViewsToModel(idps),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func (repo *IAMRepository) AddDefaultLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddLoginPolicy(ctx, policy)
}

func (repo *IAMRepository) ChangeDefaultLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangeLoginPolicy(ctx, policy)
}

func (repo *IAMRepository) SearchDefaultIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	request.AppendAggregateIDQuery(repo.SystemDefaults.IamID)
	sequence, err := repo.View.GetLatestIDPProviderSequence()
	logging.Log("EVENT-Tuiks").OnError(err).Warn("could not read latest iam sequence")
	providers, count, err := repo.View.SearchIDPProviders(request)
	if err != nil {
		return nil, err
	}
	result := &iam_model.IDPProviderSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      iam_es_model.IDPProviderViewsToModel(providers),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func (repo *IAMRepository) AddIDPProviderToLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) (*iam_model.IDPProvider, error) {
	provider.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddIDPProviderToLoginPolicy(ctx, provider)
}

func (repo *IAMRepository) RemoveIDPProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) error {
	aggregates := make([]*es_models.Aggregate, 0)
	provider.AggregateID = repo.SystemDefaults.IamID
	_, removeAgg, err := repo.IAMEventstore.PrepareRemoveIDPProviderFromLoginPolicy(ctx, provider)
	if err != nil {
		return err
	}
	aggregates = append(aggregates, removeAgg)

	externalIDPs, err := repo.View.ExternalIDPsByIDPConfigID(provider.IdpConfigID)
	if err != nil {
		return err
	}
	for _, externalIDP := range externalIDPs {
		idpRemove := &usr_model.ExternalIDP{ObjectRoot: es_models.ObjectRoot{AggregateID: externalIDP.UserID}, IDPConfigID: externalIDP.IDPConfigID, UserID: externalIDP.ExternalUserID}
		idpAgg := make([]*es_models.Aggregate, 0)
		_, idpAgg, err = repo.UserEvents.PrepareRemoveExternalIDP(ctx, idpRemove, true)
		if err != nil {
			return err
		}
		aggregates = append(aggregates, idpAgg...)
	}
	return es_sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, nil, aggregates...)
}

func (repo *IAMRepository) GetDefaultPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error) {
	policy, err := repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.PasswordComplexityViewToModel(policy), nil
}

func (repo *IAMRepository) AddDefaultPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddPasswordComplexityPolicy(ctx, policy)
}

func (repo *IAMRepository) ChangeDefaultPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangePasswordComplexityPolicy(ctx, policy)
}

func (repo *IAMRepository) GetDefaultPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error) {
	policy, err := repo.View.PasswordAgePolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.PasswordAgeViewToModel(policy), nil
}

func (repo *IAMRepository) AddDefaultPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddPasswordAgePolicy(ctx, policy)
}

func (repo *IAMRepository) ChangeDefaultPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangePasswordAgePolicy(ctx, policy)
}

func (repo *IAMRepository) GetDefaultPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error) {
	policy, err := repo.View.PasswordLockoutPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.PasswordLockoutViewToModel(policy), nil
}

func (repo *IAMRepository) AddDefaultPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddPasswordLockoutPolicy(ctx, policy)
}

func (repo *IAMRepository) ChangeDefaultPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangePasswordLockoutPolicy(ctx, policy)
}

func (repo *IAMRepository) GetOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error) {
	policy, err := repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.OrgIAMViewToModel(policy), nil
}

func (repo *IAMRepository) AddDefaultOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddOrgIAMPolicy(ctx, policy)
}

func (repo *IAMRepository) ChangeDefaultOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangeOrgIAMPolicy(ctx, policy)
}

func (repo *IAMRepository) GetDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error) {
	policy, err := repo.View.LabelPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.LabelPolicyViewToModel(policy), err
}

func (repo *IAMRepository) AddDefaultLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddLabelPolicy(ctx, policy)
}

func (repo *IAMRepository) ChangeDefaultLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangeLabelPolicy(ctx, policy)
}

func (repo *IAMRepository) GetDefaultMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error) {
	template, err := repo.View.MailTemplateByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.MailTemplateViewToModel(template), err
}

func (repo *IAMRepository) AddDefaultMailTemplate(ctx context.Context, template *iam_model.MailTemplate) (*iam_model.MailTemplate, error) {
	template.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddMailTemplate(ctx, template)
}

func (repo *IAMRepository) ChangeDefaultMailTemplate(ctx context.Context, template *iam_model.MailTemplate) (*iam_model.MailTemplate, error) {
	template.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangeMailTemplate(ctx, template)
}

func (repo *IAMRepository) SearchIAMMembersx(ctx context.Context, request *iam_model.IAMMemberSearchRequest) (*iam_model.IAMMemberSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, err := repo.View.GetLatestIAMMemberSequence()
	logging.Log("EVENT-Slkci").OnError(err).Warn("could not read latest iam sequence")
	members, count, err := repo.View.SearchIAMMembers(request)
	if err != nil {
		return nil, err
	}
	result := &iam_model.IAMMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      iam_es_model.IAMMembersToModel(members),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func (repo *IAMRepository) GetDefaultMailTexts(ctx context.Context) (*iam_model.MailTextsView, error) {
	text, err := repo.View.MailTexts(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.MailTextsViewToModel(text), err
}

func (repo *IAMRepository) GetDefaultMailText(ctx context.Context, textType string, language string) (*iam_model.MailTextView, error) {
	text, err := repo.View.MailTextByIDs(repo.SystemDefaults.IamID, textType, language)
	if err != nil {
		return nil, err
	}
	return iam_es_model.MailTextViewToModel(text), err
}

func (repo *IAMRepository) AddDefaultMailText(ctx context.Context, text *iam_model.MailText) (*iam_model.MailText, error) {
	text.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddMailText(ctx, text)
}

func (repo *IAMRepository) ChangeDefaultMailText(ctx context.Context, text *iam_model.MailText) (*iam_model.MailText, error) {
	text.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangeMailText(ctx, text)
}

func (repo *IAMRepository) GetDefaultLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error) {
	policy, err := repo.View.LoginPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.LoginPolicyViewToModel(policy), err
}
