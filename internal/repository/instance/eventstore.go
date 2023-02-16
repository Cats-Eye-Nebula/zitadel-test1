package instance

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, DefaultOrgSetEventType, DefaultOrgSetMapper).
		RegisterFilterEventMapper(AggregateType, ProjectSetEventType, ProjectSetMapper).
		RegisterFilterEventMapper(AggregateType, ConsoleSetEventType, ConsoleSetMapper).
		RegisterFilterEventMapper(AggregateType, DefaultLanguageSetEventType, DefaultLanguageSetMapper).
		RegisterFilterEventMapper(AggregateType, SecretGeneratorAddedEventType, SecretGeneratorAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, SecretGeneratorChangedEventType, SecretGeneratorChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, SecretGeneratorRemovedEventType, SecretGeneratorRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, SMTPConfigAddedEventType, SMTPConfigAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, SMTPConfigChangedEventType, SMTPConfigChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, SMTPConfigPasswordChangedEventType, SMTPConfigPasswordChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, SMTPConfigRemovedEventType, SMTPConfigRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, SMSConfigTwilioAddedEventType, SMSConfigTwilioAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, SMSConfigTwilioChangedEventType, SMSConfigTwilioChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, SMSConfigTwilioTokenChangedEventType, SMSConfigTwilioTokenChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, SMSConfigActivatedEventType, SMSConfigActivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, SMSConfigDeactivatedEventType, SMSConfigDeactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, SMSConfigRemovedEventType, SMSConfigRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, DebugNotificationProviderFileAddedEventType, DebugNotificationProviderFileAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, DebugNotificationProviderFileChangedEventType, DebugNotificationProviderFileChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, DebugNotificationProviderFileRemovedEventType, DebugNotificationProviderFileRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, DebugNotificationProviderLogAddedEventType, DebugNotificationProviderLogAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, DebugNotificationProviderLogChangedEventType, DebugNotificationProviderLogChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, DebugNotificationProviderLogRemovedEventType, DebugNotificationProviderLogRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, OIDCSettingsAddedEventType, OIDCSettingsAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, OIDCSettingsChangedEventType, OIDCSettingsChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, SecurityPolicySetEventType, SecurityPolicySetEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyAddedEventType, LabelPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyChangedEventType, LabelPolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyActivatedEventType, LabelPolicyActivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyLogoAddedEventType, LabelPolicyLogoAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyLogoRemovedEventType, LabelPolicyLogoRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyIconAddedEventType, LabelPolicyIconAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyIconRemovedEventType, LabelPolicyIconRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyLogoDarkAddedEventType, LabelPolicyLogoDarkAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyLogoDarkRemovedEventType, LabelPolicyLogoDarkRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyIconDarkAddedEventType, LabelPolicyIconDarkAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyIconDarkRemovedEventType, LabelPolicyIconDarkRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyFontAddedEventType, LabelPolicyFontAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyFontRemovedEventType, LabelPolicyFontRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyAssetsRemovedEventType, LabelPolicyAssetsRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyAddedEventType, LoginPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyChangedEventType, LoginPolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, DomainPolicyAddedEventType, DomainPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, DomainPolicyChangedEventType, DomainPolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasswordAgePolicyAddedEventType, PasswordAgePolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasswordAgePolicyChangedEventType, PasswordAgePolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasswordComplexityPolicyAddedEventType, PasswordComplexityPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasswordComplexityPolicyChangedEventType, PasswordComplexityPolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, LockoutPolicyAddedEventType, LockoutPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LockoutPolicyChangedEventType, LockoutPolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, PrivacyPolicyAddedEventType, PrivacyPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, PrivacyPolicyChangedEventType, PrivacyPolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, MemberAddedEventType, MemberAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, MemberChangedEventType, MemberChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, MemberRemovedEventType, MemberRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, MemberCascadeRemovedEventType, MemberCascadeRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, OAuthIDPAddedEventType, OAuthIDPAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, OAuthIDPChangedEventType, OAuthIDPChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, OIDCIDPAddedEventType, OIDCIDPAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, OIDCIDPChangedEventType, OIDCIDPChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, JWTIDPAddedEventType, JWTIDPAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, JWTIDPChangedEventType, JWTIDPChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, AzureADIDPAddedEventType, AzureADIDPAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, AzureADIDPChangedEventType, AzureADIDPChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, GitHubIDPAddedEventType, GitHubIDPAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, GitHubIDPChangedEventType, GitHubIDPChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, GitHubEnterpriseIDPAddedEventType, GitHubEnterpriseIDPAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, GitHubEnterpriseIDPChangedEventType, GitHubEnterpriseIDPChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, GitLabIDPAddedEventType, GitLabIDPAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, GitLabIDPChangedEventType, GitLabIDPChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, GitLabSelfHostedIDPAddedEventType, GitLabSelfHostedIDPAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, GitLabSelfHostedIDPChangedEventType, GitLabSelfHostedIDPChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, GoogleIDPAddedEventType, GoogleIDPAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, GoogleIDPChangedEventType, GoogleIDPChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPRemovedEventType, IDPRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPConfigAddedEventType, IDPConfigAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPConfigChangedEventType, IDPConfigChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPConfigRemovedEventType, IDPConfigRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPConfigDeactivatedEventType, IDPConfigDeactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPConfigReactivatedEventType, IDPConfigReactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPOIDCConfigAddedEventType, IDPOIDCConfigAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPOIDCConfigChangedEventType, IDPOIDCConfigChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPJWTConfigAddedEventType, IDPJWTConfigAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPJWTConfigChangedEventType, IDPJWTConfigChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, LDAPIDPAddedEventType, LDAPIDPAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LDAPIDPChangedEventType, LDAPIDPChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPRemovedEventType, IDPRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyIDPProviderAddedEventType, IdentityProviderAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyIDPProviderRemovedEventType, IdentityProviderRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyIDPProviderCascadeRemovedEventType, IdentityProviderCascadeRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicySecondFactorAddedEventType, SecondFactorAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicySecondFactorRemovedEventType, SecondFactorRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyMultiFactorAddedEventType, MultiFactorAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyMultiFactorRemovedEventType, MultiFactorRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, MailTemplateAddedEventType, MailTemplateAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, MailTemplateChangedEventType, MailTemplateChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, MailTextAddedEventType, MailTextAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, MailTextChangedEventType, MailTextChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, CustomTextSetEventType, CustomTextSetEventMapper).
		RegisterFilterEventMapper(AggregateType, CustomTextRemovedEventType, CustomTextRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, CustomTextTemplateRemovedEventType, CustomTextTemplateRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, InstanceDomainAddedEventType, DomainAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, InstanceDomainPrimarySetEventType, DomainPrimarySetEventMapper).
		RegisterFilterEventMapper(AggregateType, InstanceDomainRemovedEventType, DomainRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, InstanceAddedEventType, InstanceAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, InstanceChangedEventType, InstanceChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, InstanceRemovedEventType, InstanceRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, NotificationPolicyAddedEventType, NotificationPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, NotificationPolicyChangedEventType, NotificationPolicyChangedEventMapper)
}
