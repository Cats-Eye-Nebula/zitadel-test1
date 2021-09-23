package query

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
)

func readModelToIAM(readModel *ReadModel) *model.IAM {
	return &model.IAM{
		ObjectRoot:                      readModelToObjectRoot(readModel.ReadModel),
		GlobalOrgID:                     readModel.GlobalOrgID,
		IAMProjectID:                    readModel.ProjectID,
		SetUpDone:                       readModel.SetUpDone,
		SetUpStarted:                    readModel.SetUpStarted,
		Members:                         readModelToMembers(&readModel.Members),
		DefaultLabelPolicy:              readModelToLabelPolicy(&readModel.DefaultLabelPolicy),
		DefaultLoginPolicy:              readModelToLoginPolicy(&readModel.DefaultLoginPolicy),
		DefaultOrgIAMPolicy:             readModelToOrgIAMPolicy(&readModel.DefaultOrgIAMPolicy),
		DefaultPasswordAgePolicy:        readModelToPasswordAgePolicy(&readModel.DefaultPasswordAgePolicy),
		DefaultPasswordComplexityPolicy: readModelToPasswordComplexityPolicy(&readModel.DefaultPasswordComplexityPolicy),
		DefaultLockoutPolicy:            readModelToPasswordLockoutPolicy(&readModel.DefaultPasswordLockoutPolicy),
		IDPs:                            readModelToIDPConfigs(&readModel.IDPs),
	}
}

func readModelToIDPConfigView(rm *IAMIDPConfigReadModel) *domain.IDPConfigView {
	converted := &domain.IDPConfigView{
		AggregateID:     rm.AggregateID,
		ChangeDate:      rm.ChangeDate,
		CreationDate:    rm.CreationDate,
		IDPConfigID:     rm.ConfigID,
		IDPProviderType: rm.ProviderType,
		IsOIDC:          rm.OIDCConfig != nil,
		Name:            rm.Name,
		Sequence:        rm.ProcessedSequence,
		State:           rm.State,
		StylingType:     rm.StylingType,
	}
	if rm.OIDCConfig != nil {
		converted.OIDCClientID = rm.OIDCConfig.ClientID
		converted.OIDCClientSecret = rm.OIDCConfig.ClientSecret
		converted.OIDCIDPDisplayNameMapping = rm.OIDCConfig.IDPDisplayNameMapping
		converted.OIDCIssuer = rm.OIDCConfig.Issuer
		converted.OIDCScopes = rm.OIDCConfig.Scopes
		converted.OIDCUsernameMapping = rm.OIDCConfig.UserNameMapping
		converted.OAuthAuthorizationEndpoint = rm.OIDCConfig.AuthorizationEndpoint
		converted.OAuthTokenEndpoint = rm.OIDCConfig.TokenEndpoint
	}
	if rm.JWTConfig != nil {
		converted.JWTEndpoint = rm.JWTConfig.JWTEndpoint
		converted.JWTIssuer = rm.JWTConfig.Issuer
		converted.JWTKeysEndpoint = rm.JWTConfig.KeysEndpoint
	}
	return converted
}

func readModelToMember(readModel *MemberReadModel) *model.IAMMember {
	return &model.IAMMember{
		ObjectRoot: readModelToObjectRoot(readModel.ReadModel),
		Roles:      readModel.Roles,
		UserID:     readModel.UserID,
	}
}

func readModelToMembers(readModel *IAMMembersReadModel) []*model.IAMMember {
	members := make([]*model.IAMMember, len(readModel.Members))

	for i, member := range readModel.Members {
		members[i] = &model.IAMMember{
			ObjectRoot: readModelToObjectRoot(member.ReadModel),
			Roles:      member.Roles,
			UserID:     member.UserID,
		}
	}

	return members
}

func readModelToLabelPolicy(readModel *IAMLabelPolicyReadModel) *model.LabelPolicy {
	return &model.LabelPolicy{
		ObjectRoot:          readModelToObjectRoot(readModel.LabelPolicyReadModel.ReadModel),
		PrimaryColor:        readModel.PrimaryColor,
		BackgroundColor:     readModel.BackgroundColor,
		WarnColor:           readModel.WarnColor,
		FontColor:           readModel.FontColor,
		PrimaryColorDark:    readModel.PrimaryColorDark,
		BackgroundColorDark: readModel.BackgroundColorDark,
		WarnColorDark:       readModel.WarnColorDark,
		FontColorDark:       readModel.FontColorDark,
		Default:             true,
	}
}

func readModelToLoginPolicy(readModel *IAMLoginPolicyReadModel) *model.LoginPolicy {
	return &model.LoginPolicy{
		ObjectRoot:            readModelToObjectRoot(readModel.LoginPolicyReadModel.ReadModel),
		AllowExternalIdp:      readModel.AllowExternalIDP,
		AllowRegister:         readModel.AllowRegister,
		AllowUsernamePassword: readModel.AllowUserNamePassword,
		Default:               true,
	}
}
func readModelToOrgIAMPolicy(readModel *IAMOrgIAMPolicyReadModel) *model.OrgIAMPolicy {
	return &model.OrgIAMPolicy{
		ObjectRoot:            readModelToObjectRoot(readModel.OrgIAMPolicyReadModel.ReadModel),
		UserLoginMustBeDomain: readModel.UserLoginMustBeDomain,
		Default:               true,
	}
}
func readModelToPasswordAgePolicy(readModel *IAMPasswordAgePolicyReadModel) *model.PasswordAgePolicy {
	return &model.PasswordAgePolicy{
		ObjectRoot:     readModelToObjectRoot(readModel.PasswordAgePolicyReadModel.ReadModel),
		ExpireWarnDays: readModel.ExpireWarnDays,
		MaxAgeDays:     readModel.MaxAgeDays,
	}
}
func readModelToPasswordComplexityPolicy(readModel *IAMPasswordComplexityPolicyReadModel) *model.PasswordComplexityPolicy {
	return &model.PasswordComplexityPolicy{
		ObjectRoot:   readModelToObjectRoot(readModel.PasswordComplexityPolicyReadModel.ReadModel),
		HasLowercase: readModel.HasLowercase,
		HasNumber:    readModel.HasNumber,
		HasSymbol:    readModel.HasSymbol,
		HasUppercase: readModel.HasUpperCase,
		MinLength:    readModel.MinLength,
	}
}
func readModelToPasswordLockoutPolicy(readModel *IAMLockoutPolicyReadModel) *model.LockoutPolicy {
	return &model.LockoutPolicy{
		ObjectRoot:          readModelToObjectRoot(readModel.LockoutPolicyReadModel.ReadModel),
		MaxPasswordAttempts: readModel.MaxAttempts,
		ShowLockOutFailures: readModel.ShowLockOutFailures,
	}
}

func readModelToIDPConfigs(rm *IAMIDPConfigsReadModel) []*model.IDPConfig {
	configs := make([]*model.IDPConfig, len(rm.Configs))
	for i, config := range rm.Configs {
		configs[i] = readModelToIDPConfig(&IAMIDPConfigReadModel{IDPConfigReadModel: *config})
	}
	return configs
}

func readModelToIDPConfig(rm *IAMIDPConfigReadModel) *model.IDPConfig {
	config := &model.IDPConfig{
		ObjectRoot:  readModelToObjectRoot(rm.ReadModel),
		IDPConfigID: rm.ConfigID,
		Name:        rm.Name,
		State:       model.IDPConfigState(rm.State),
		StylingType: model.IDPStylingType(rm.StylingType),
	}
	if rm.OIDCConfig != nil {
		config.OIDCConfig = readModelToIDPOIDCConfig(rm.OIDCConfig)
	}
	if rm.JWTConfig != nil {
		config.JWTIDPConfig = readModelToIDPJWTConfig(rm.JWTConfig)
	}
	return config
}

func readModelToIDPOIDCConfig(rm *OIDCConfigReadModel) *model.OIDCIDPConfig {
	return &model.OIDCIDPConfig{
		ObjectRoot:            readModelToObjectRoot(rm.ReadModel),
		ClientID:              rm.ClientID,
		ClientSecret:          rm.ClientSecret,
		ClientSecretString:    string(rm.ClientSecret.Crypted),
		IDPConfigID:           rm.IDPConfigID,
		IDPDisplayNameMapping: model.OIDCMappingField(rm.IDPDisplayNameMapping),
		Issuer:                rm.Issuer,
		Scopes:                rm.Scopes,
		UsernameMapping:       model.OIDCMappingField(rm.UserNameMapping),
	}
}

func readModelToIDPJWTConfig(rm *JWTConfigReadModel) *model.JWTIDPConfig {
	return &model.JWTIDPConfig{
		ObjectRoot:   readModelToObjectRoot(rm.ReadModel),
		IDPConfigID:  rm.IDPConfigID,
		JWTEndpoint:  rm.JWTEndpoint,
		Issuer:       rm.Issuer,
		KeysEndpoint: rm.KeysEndpoint,
	}
}

func readModelToObjectRoot(readModel eventstore.ReadModel) models.ObjectRoot {
	return models.ObjectRoot{
		AggregateID:   readModel.AggregateID,
		ChangeDate:    readModel.ChangeDate,
		CreationDate:  readModel.CreationDate,
		ResourceOwner: readModel.ResourceOwner,
		Sequence:      readModel.ProcessedSequence,
	}
}
