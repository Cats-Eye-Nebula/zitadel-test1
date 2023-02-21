---
title: zitadel/admin.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)


## AdminService {#zitadeladminv1adminservice}


### Healthz

> **rpc** Healthz([HealthzRequest](#healthzrequest))
[HealthzResponse](#healthzresponse)

Indicates if ZITADEL is running.
It respondes as soon as ZITADEL started



    GET: /healthz


### GetSupportedLanguages

> **rpc** GetSupportedLanguages([GetSupportedLanguagesRequest](#getsupportedlanguagesrequest))
[GetSupportedLanguagesResponse](#getsupportedlanguagesresponse)

Returns the default languages



    GET: /languages


### SetDefaultLanguage

> **rpc** SetDefaultLanguage([SetDefaultLanguageRequest](#setdefaultlanguagerequest))
[SetDefaultLanguageResponse](#setdefaultlanguageresponse)

Set the default language



    PUT: /languages/default/{language}


### GetDefaultLanguage

> **rpc** GetDefaultLanguage([GetDefaultLanguageRequest](#getdefaultlanguagerequest))
[GetDefaultLanguageResponse](#getdefaultlanguageresponse)

Set the default language



    GET: /languages/default


### GetMyInstance

> **rpc** GetMyInstance([GetMyInstanceRequest](#getmyinstancerequest))
[GetMyInstanceResponse](#getmyinstanceresponse)

Returns the details of the instance



    GET: /instances/me


### ListInstanceDomains

> **rpc** ListInstanceDomains([ListInstanceDomainsRequest](#listinstancedomainsrequest))
[ListInstanceDomainsResponse](#listinstancedomainsresponse)

Returns the domains of the instance



    POST: /domains/_search


### ListSecretGenerators

> **rpc** ListSecretGenerators([ListSecretGeneratorsRequest](#listsecretgeneratorsrequest))
[ListSecretGeneratorsResponse](#listsecretgeneratorsresponse)

Set the default language



    POST: /secretgenerators/_search


### GetSecretGenerator

> **rpc** GetSecretGenerator([GetSecretGeneratorRequest](#getsecretgeneratorrequest))
[GetSecretGeneratorResponse](#getsecretgeneratorresponse)

Get Secret Generator by type (e.g PasswordResetCode)



    GET: /secretgenerators/{generator_type}


### UpdateSecretGenerator

> **rpc** UpdateSecretGenerator([UpdateSecretGeneratorRequest](#updatesecretgeneratorrequest))
[UpdateSecretGeneratorResponse](#updatesecretgeneratorresponse)

Update secret generator configuration



    PUT: /secretgenerators/{generator_type}


### GetSMTPConfig

> **rpc** GetSMTPConfig([GetSMTPConfigRequest](#getsmtpconfigrequest))
[GetSMTPConfigResponse](#getsmtpconfigresponse)

Get system smtp configuration



    GET: /smtp


### AddSMTPConfig

> **rpc** AddSMTPConfig([AddSMTPConfigRequest](#addsmtpconfigrequest))
[AddSMTPConfigResponse](#addsmtpconfigresponse)

Add system smtp configuration



    POST: /smtp


### UpdateSMTPConfig

> **rpc** UpdateSMTPConfig([UpdateSMTPConfigRequest](#updatesmtpconfigrequest))
[UpdateSMTPConfigResponse](#updatesmtpconfigresponse)

Update system smtp configuration



    PUT: /smtp


### UpdateSMTPConfigPassword

> **rpc** UpdateSMTPConfigPassword([UpdateSMTPConfigPasswordRequest](#updatesmtpconfigpasswordrequest))
[UpdateSMTPConfigPasswordResponse](#updatesmtpconfigpasswordresponse)

Update system smtp configuration password for host



    PUT: /smtp/password


### RemoveSMTPConfig

> **rpc** RemoveSMTPConfig([RemoveSMTPConfigRequest](#removesmtpconfigrequest))
[RemoveSMTPConfigResponse](#removesmtpconfigresponse)

Remove system smtp configuration



    DELETE: /smtp


### ListSMSProviders

> **rpc** ListSMSProviders([ListSMSProvidersRequest](#listsmsprovidersrequest))
[ListSMSProvidersResponse](#listsmsprovidersresponse)

list sms provider configurations



    POST: /sms/_search


### GetSMSProvider

> **rpc** GetSMSProvider([GetSMSProviderRequest](#getsmsproviderrequest))
[GetSMSProviderResponse](#getsmsproviderresponse)

Get sms provider



    GET: /sms/{id}


### AddSMSProviderTwilio

> **rpc** AddSMSProviderTwilio([AddSMSProviderTwilioRequest](#addsmsprovidertwiliorequest))
[AddSMSProviderTwilioResponse](#addsmsprovidertwilioresponse)

Add twilio sms provider



    POST: /sms/twilio


### UpdateSMSProviderTwilio

> **rpc** UpdateSMSProviderTwilio([UpdateSMSProviderTwilioRequest](#updatesmsprovidertwiliorequest))
[UpdateSMSProviderTwilioResponse](#updatesmsprovidertwilioresponse)

Update twilio sms provider



    PUT: /sms/twilio/{id}


### UpdateSMSProviderTwilioToken

> **rpc** UpdateSMSProviderTwilioToken([UpdateSMSProviderTwilioTokenRequest](#updatesmsprovidertwiliotokenrequest))
[UpdateSMSProviderTwilioTokenResponse](#updatesmsprovidertwiliotokenresponse)

Update twilio sms provider token



    PUT: /sms/twilio/{id}/token


### ActivateSMSProvider

> **rpc** ActivateSMSProvider([ActivateSMSProviderRequest](#activatesmsproviderrequest))
[ActivateSMSProviderResponse](#activatesmsproviderresponse)

Activate sms provider



    POST: /sms/{id}/_activate


### DeactivateSMSProvider

> **rpc** DeactivateSMSProvider([DeactivateSMSProviderRequest](#deactivatesmsproviderrequest))
[DeactivateSMSProviderResponse](#deactivatesmsproviderresponse)

Deactivate sms provider



    POST: /sms/{id}/_deactivate


### RemoveSMSProvider

> **rpc** RemoveSMSProvider([RemoveSMSProviderRequest](#removesmsproviderrequest))
[RemoveSMSProviderResponse](#removesmsproviderresponse)

Remove sms provider token



    DELETE: /sms/{id}


### GetOIDCSettings

> **rpc** GetOIDCSettings([GetOIDCSettingsRequest](#getoidcsettingsrequest))
[GetOIDCSettingsResponse](#getoidcsettingsresponse)

Get OIDC settings (e.g token lifetimes, etc.)



    GET: /settings/oidc


### AddOIDCSettings

> **rpc** AddOIDCSettings([AddOIDCSettingsRequest](#addoidcsettingsrequest))
[AddOIDCSettingsResponse](#addoidcsettingsresponse)

Add oidc settings (e.g token lifetimes, etc)



    POST: /settings/oidc


### UpdateOIDCSettings

> **rpc** UpdateOIDCSettings([UpdateOIDCSettingsRequest](#updateoidcsettingsrequest))
[UpdateOIDCSettingsResponse](#updateoidcsettingsresponse)

Update oidc settings (e.g token lifetimes, etc)



    PUT: /settings/oidc


### GetFileSystemNotificationProvider

> **rpc** GetFileSystemNotificationProvider([GetFileSystemNotificationProviderRequest](#getfilesystemnotificationproviderrequest))
[GetFileSystemNotificationProviderResponse](#getfilesystemnotificationproviderresponse)

Get file system notification provider



    GET: /notification/provider/file


### GetLogNotificationProvider

> **rpc** GetLogNotificationProvider([GetLogNotificationProviderRequest](#getlognotificationproviderrequest))
[GetLogNotificationProviderResponse](#getlognotificationproviderresponse)

Get log notification provider



    GET: /notification/provider/log


### GetSecurityPolicy

> **rpc** GetSecurityPolicy([GetSecurityPolicyRequest](#getsecuritypolicyrequest))
[GetSecurityPolicyResponse](#getsecuritypolicyresponse)

Get the security policy



    GET: /policies/security


### SetSecurityPolicy

> **rpc** SetSecurityPolicy([SetSecurityPolicyRequest](#setsecuritypolicyrequest))
[SetSecurityPolicyResponse](#setsecuritypolicyresponse)

set the security policy



    PUT: /policies/security


### GetOrgByID

> **rpc** GetOrgByID([GetOrgByIDRequest](#getorgbyidrequest))
[GetOrgByIDResponse](#getorgbyidresponse)

Returns an organisation by id



    GET: /orgs/{id}


### IsOrgUnique

> **rpc** IsOrgUnique([IsOrgUniqueRequest](#isorguniquerequest))
[IsOrgUniqueResponse](#isorguniqueresponse)

Checks whether an organisation exists by the given parameters



    GET: /orgs/_is_unique


### SetDefaultOrg

> **rpc** SetDefaultOrg([SetDefaultOrgRequest](#setdefaultorgrequest))
[SetDefaultOrgResponse](#setdefaultorgresponse)

Set the default org



    PUT: /orgs/default/{org_id}


### GetDefaultOrg

> **rpc** GetDefaultOrg([GetDefaultOrgRequest](#getdefaultorgrequest))
[GetDefaultOrgResponse](#getdefaultorgresponse)

Set the default org



    GET: /orgs/default


### ListOrgs

> **rpc** ListOrgs([ListOrgsRequest](#listorgsrequest))
[ListOrgsResponse](#listorgsresponse)

Returns all organisations matching the request
all queries need to match (AND)



    POST: /orgs/_search


### SetUpOrg

> **rpc** SetUpOrg([SetUpOrgRequest](#setuporgrequest))
[SetUpOrgResponse](#setuporgresponse)

Creates a new org and user
and adds the user to the orgs members as ORG_OWNER



    POST: /orgs/_setup


### RemoveOrg

> **rpc** RemoveOrg([RemoveOrgRequest](#removeorgrequest))
[RemoveOrgResponse](#removeorgresponse)

Sets the state of the organisation and all its resource (Users, Projects, Grants to and from the org) to removed
Users of this organisation will not be able login



    DELETE: /orgs/{org_id}


### GetIDPByID

> **rpc** GetIDPByID([GetIDPByIDRequest](#getidpbyidrequest))
[GetIDPByIDResponse](#getidpbyidresponse)

Returns a identity provider configuration of the IAM instance



    GET: /idps/{id}


### ListIDPs

> **rpc** ListIDPs([ListIDPsRequest](#listidpsrequest))
[ListIDPsResponse](#listidpsresponse)

Returns all identity provider configurations of the IAM instance



    POST: /idps/_search


### AddOIDCIDP

> **rpc** AddOIDCIDP([AddOIDCIDPRequest](#addoidcidprequest))
[AddOIDCIDPResponse](#addoidcidpresponse)

Adds a new oidc identity provider configuration the IAM instance



    POST: /idps/oidc


### AddJWTIDP

> **rpc** AddJWTIDP([AddJWTIDPRequest](#addjwtidprequest))
[AddJWTIDPResponse](#addjwtidpresponse)

Adds a new jwt identity provider configuration the IAM instance



    POST: /idps/jwt


### UpdateIDP

> **rpc** UpdateIDP([UpdateIDPRequest](#updateidprequest))
[UpdateIDPResponse](#updateidpresponse)

Updates the specified idp
all fields are updated. If no value is provided the field will be empty afterwards.



    PUT: /idps/{idp_id}


### DeactivateIDP

> **rpc** DeactivateIDP([DeactivateIDPRequest](#deactivateidprequest))
[DeactivateIDPResponse](#deactivateidpresponse)

Sets the state of the idp to IDP_STATE_INACTIVE
the state MUST be IDP_STATE_ACTIVE for this call



    POST: /idps/{idp_id}/_deactivate


### ReactivateIDP

> **rpc** ReactivateIDP([ReactivateIDPRequest](#reactivateidprequest))
[ReactivateIDPResponse](#reactivateidpresponse)

Sets the state of the idp to IDP_STATE_ACTIVE
the state MUST be IDP_STATE_INACTIVE for this call



    POST: /idps/{idp_id}/_reactivate


### RemoveIDP

> **rpc** RemoveIDP([RemoveIDPRequest](#removeidprequest))
[RemoveIDPResponse](#removeidpresponse)

RemoveIDP deletes the IDP permanetly



    DELETE: /idps/{idp_id}


### UpdateIDPOIDCConfig

> **rpc** UpdateIDPOIDCConfig([UpdateIDPOIDCConfigRequest](#updateidpoidcconfigrequest))
[UpdateIDPOIDCConfigResponse](#updateidpoidcconfigresponse)

Updates the oidc configuration of the specified idp
all fields are updated. If no value is provided the field will be empty afterwards.



    PUT: /idps/{idp_id}/oidc_config


### UpdateIDPJWTConfig

> **rpc** UpdateIDPJWTConfig([UpdateIDPJWTConfigRequest](#updateidpjwtconfigrequest))
[UpdateIDPJWTConfigResponse](#updateidpjwtconfigresponse)

Updates the jwt configuration of the specified idp
all fields are updated. If no value is provided the field will be empty afterwards.



    PUT: /idps/{idp_id}/jwt_config


### ListProviders

> **rpc** ListProviders([ListProvidersRequest](#listprovidersrequest))
[ListProvidersResponse](#listprovidersresponse)

Returns all identity providers, which match the query
Limit should always be set, there is a default limit set by the service



    POST: /idps/templates/_search


### GetProviderByID

> **rpc** GetProviderByID([GetProviderByIDRequest](#getproviderbyidrequest))
[GetProviderByIDResponse](#getproviderbyidresponse)

Returns an identity provider of the instance



    GET: /idps/templates/{id}


### AddGoogleProvider

> **rpc** AddGoogleProvider([AddGoogleProviderRequest](#addgoogleproviderrequest))
[AddGoogleProviderResponse](#addgoogleproviderresponse)

Add a new Google identity provider on the instance



    POST: /idps/google


### UpdateGoogleProvider

> **rpc** UpdateGoogleProvider([UpdateGoogleProviderRequest](#updategoogleproviderrequest))
[UpdateGoogleProviderResponse](#updategoogleproviderresponse)

Change an existing Google identity provider on the instance



    POST: /idps/google/{id}


### AddLDAPProvider

> **rpc** AddLDAPProvider([AddLDAPProviderRequest](#addldapproviderrequest))
[AddLDAPProviderResponse](#addldapproviderresponse)

Add a new LDAP identity provider on the instance



    POST: /idps/ldap


### UpdateLDAPProvider

> **rpc** UpdateLDAPProvider([UpdateLDAPProviderRequest](#updateldapproviderrequest))
[UpdateLDAPProviderResponse](#updateldapproviderresponse)

Change an existing LDAP identity provider on the instance



    PUT: /idps/ldap/{id}


### DeleteProvider

> **rpc** DeleteProvider([DeleteProviderRequest](#deleteproviderrequest))
[DeleteProviderResponse](#deleteproviderresponse)

Remove an identity provider
Will remove all linked providers of this configuration on the users



    DELETE: /idps/templates/{id}


### GetOrgIAMPolicy

> **rpc** GetOrgIAMPolicy([GetOrgIAMPolicyRequest](#getorgiampolicyrequest))
[GetOrgIAMPolicyResponse](#getorgiampolicyresponse)

deprecated: please use DomainPolicy instead
Returns the Org IAM policy defined by the administrators of ZITADEL



    GET: /policies/orgiam


### UpdateOrgIAMPolicy

> **rpc** UpdateOrgIAMPolicy([UpdateOrgIAMPolicyRequest](#updateorgiampolicyrequest))
[UpdateOrgIAMPolicyResponse](#updateorgiampolicyresponse)

deprecated: please use DomainPolicy instead
Updates the default OrgIAM policy.
it impacts all organisations without a customised policy



    PUT: /policies/orgiam


### GetCustomOrgIAMPolicy

> **rpc** GetCustomOrgIAMPolicy([GetCustomOrgIAMPolicyRequest](#getcustomorgiampolicyrequest))
[GetCustomOrgIAMPolicyResponse](#getcustomorgiampolicyresponse)

deprecated: please use DomainPolicy instead
Returns the customised policy or the default if not customised



    GET: /orgs/{org_id}/policies/orgiam


### AddCustomOrgIAMPolicy

> **rpc** AddCustomOrgIAMPolicy([AddCustomOrgIAMPolicyRequest](#addcustomorgiampolicyrequest))
[AddCustomOrgIAMPolicyResponse](#addcustomorgiampolicyresponse)

deprecated: please use DomainPolicy instead
Defines a custom OrgIAM policy as specified



    POST: /orgs/{org_id}/policies/orgiam


### UpdateCustomOrgIAMPolicy

> **rpc** UpdateCustomOrgIAMPolicy([UpdateCustomOrgIAMPolicyRequest](#updatecustomorgiampolicyrequest))
[UpdateCustomOrgIAMPolicyResponse](#updatecustomorgiampolicyresponse)

deprecated: please use DomainPolicy instead
Updates a custom OrgIAM policy as specified



    PUT: /orgs/{org_id}/policies/orgiam


### ResetCustomOrgIAMPolicyToDefault

> **rpc** ResetCustomOrgIAMPolicyToDefault([ResetCustomOrgIAMPolicyToDefaultRequest](#resetcustomorgiampolicytodefaultrequest))
[ResetCustomOrgIAMPolicyToDefaultResponse](#resetcustomorgiampolicytodefaultresponse)

deprecated: please use DomainPolicy instead
Resets the org iam policy of the organisation to default
ZITADEL will fallback to the default policy defined by the ZITADEL administrators



    DELETE: /orgs/{org_id}/policies/orgiam


### GetDomainPolicy

> **rpc** GetDomainPolicy([GetDomainPolicyRequest](#getdomainpolicyrequest))
[GetDomainPolicyResponse](#getdomainpolicyresponse)

Returns the Domain policy defined by the administrators of ZITADEL



    GET: /policies/domain


### UpdateDomainPolicy

> **rpc** UpdateDomainPolicy([UpdateDomainPolicyRequest](#updatedomainpolicyrequest))
[UpdateDomainPolicyResponse](#updatedomainpolicyresponse)

Updates the default Domain policy.
it impacts all organisations without a customised policy



    PUT: /policies/domain


### GetCustomDomainPolicy

> **rpc** GetCustomDomainPolicy([GetCustomDomainPolicyRequest](#getcustomdomainpolicyrequest))
[GetCustomDomainPolicyResponse](#getcustomdomainpolicyresponse)

Returns the customised policy or the default if not customised



    GET: /orgs/{org_id}/policies/domain


### AddCustomDomainPolicy

> **rpc** AddCustomDomainPolicy([AddCustomDomainPolicyRequest](#addcustomdomainpolicyrequest))
[AddCustomDomainPolicyResponse](#addcustomdomainpolicyresponse)

Defines a custom Domain policy as specified



    POST: /orgs/{org_id}/policies/domain


### UpdateCustomDomainPolicy

> **rpc** UpdateCustomDomainPolicy([UpdateCustomDomainPolicyRequest](#updatecustomdomainpolicyrequest))
[UpdateCustomDomainPolicyResponse](#updatecustomdomainpolicyresponse)

Updates a custom Domain policy as specified



    PUT: /orgs/{org_id}/policies/domain


### ResetCustomDomainPolicyToDefault

> **rpc** ResetCustomDomainPolicyToDefault([ResetCustomDomainPolicyToDefaultRequest](#resetcustomdomainpolicytodefaultrequest))
[ResetCustomDomainPolicyToDefaultResponse](#resetcustomdomainpolicytodefaultresponse)

Resets the org iam policy of the organisation to default
ZITADEL will fallback to the default policy defined by the ZITADEL administrators



    DELETE: /orgs/{org_id}/policies/domain


### GetLabelPolicy

> **rpc** GetLabelPolicy([GetLabelPolicyRequest](#getlabelpolicyrequest))
[GetLabelPolicyResponse](#getlabelpolicyresponse)

Returns the label policy defined by the administrators of ZITADEL



    GET: /policies/label


### GetPreviewLabelPolicy

> **rpc** GetPreviewLabelPolicy([GetPreviewLabelPolicyRequest](#getpreviewlabelpolicyrequest))
[GetPreviewLabelPolicyResponse](#getpreviewlabelpolicyresponse)

Returns the preview label policy defined by the administrators of ZITADEL



    GET: /policies/label/_preview


### UpdateLabelPolicy

> **rpc** UpdateLabelPolicy([UpdateLabelPolicyRequest](#updatelabelpolicyrequest))
[UpdateLabelPolicyResponse](#updatelabelpolicyresponse)

Updates the default label policy of ZITADEL
it impacts all organisations without a customised policy



    PUT: /policies/label


### ActivateLabelPolicy

> **rpc** ActivateLabelPolicy([ActivateLabelPolicyRequest](#activatelabelpolicyrequest))
[ActivateLabelPolicyResponse](#activatelabelpolicyresponse)

Activates all changes of the label policy



    POST: /policies/label/_activate


### RemoveLabelPolicyLogo

> **rpc** RemoveLabelPolicyLogo([RemoveLabelPolicyLogoRequest](#removelabelpolicylogorequest))
[RemoveLabelPolicyLogoResponse](#removelabelpolicylogoresponse)

Removes the logo of the label policy



    DELETE: /policies/label/logo


### RemoveLabelPolicyLogoDark

> **rpc** RemoveLabelPolicyLogoDark([RemoveLabelPolicyLogoDarkRequest](#removelabelpolicylogodarkrequest))
[RemoveLabelPolicyLogoDarkResponse](#removelabelpolicylogodarkresponse)

Removes the logo dark of the label policy



    DELETE: /policies/label/logo_dark


### RemoveLabelPolicyIcon

> **rpc** RemoveLabelPolicyIcon([RemoveLabelPolicyIconRequest](#removelabelpolicyiconrequest))
[RemoveLabelPolicyIconResponse](#removelabelpolicyiconresponse)

Removes the icon of the label policy



    DELETE: /policies/label/icon


### RemoveLabelPolicyIconDark

> **rpc** RemoveLabelPolicyIconDark([RemoveLabelPolicyIconDarkRequest](#removelabelpolicyicondarkrequest))
[RemoveLabelPolicyIconDarkResponse](#removelabelpolicyicondarkresponse)

Removes the logo dark of the label policy



    DELETE: /policies/label/icon_dark


### RemoveLabelPolicyFont

> **rpc** RemoveLabelPolicyFont([RemoveLabelPolicyFontRequest](#removelabelpolicyfontrequest))
[RemoveLabelPolicyFontResponse](#removelabelpolicyfontresponse)

Removes the font of the label policy



    DELETE: /policies/label/font


### GetLoginPolicy

> **rpc** GetLoginPolicy([GetLoginPolicyRequest](#getloginpolicyrequest))
[GetLoginPolicyResponse](#getloginpolicyresponse)

Returns the login policy defined by the administrators of ZITADEL



    GET: /policies/login


### UpdateLoginPolicy

> **rpc** UpdateLoginPolicy([UpdateLoginPolicyRequest](#updateloginpolicyrequest))
[UpdateLoginPolicyResponse](#updateloginpolicyresponse)

Updates the default login policy of ZITADEL
it impacts all organisations without a customised policy



    PUT: /policies/login


### ListLoginPolicyIDPs

> **rpc** ListLoginPolicyIDPs([ListLoginPolicyIDPsRequest](#listloginpolicyidpsrequest))
[ListLoginPolicyIDPsResponse](#listloginpolicyidpsresponse)

Returns the idps linked to the default login policy,
defined by the administrators of ZITADEL



    POST: /policies/login/idps/_search


### AddIDPToLoginPolicy

> **rpc** AddIDPToLoginPolicy([AddIDPToLoginPolicyRequest](#addidptologinpolicyrequest))
[AddIDPToLoginPolicyResponse](#addidptologinpolicyresponse)

Adds the povided idp to the default login policy.
It impacts all organisations without a customised policy



    POST: /policies/login/idps


### RemoveIDPFromLoginPolicy

> **rpc** RemoveIDPFromLoginPolicy([RemoveIDPFromLoginPolicyRequest](#removeidpfromloginpolicyrequest))
[RemoveIDPFromLoginPolicyResponse](#removeidpfromloginpolicyresponse)

Removes the povided idp from the default login policy.
It impacts all organisations without a customised policy



    DELETE: /policies/login/idps/{idp_id}


### ListLoginPolicySecondFactors

> **rpc** ListLoginPolicySecondFactors([ListLoginPolicySecondFactorsRequest](#listloginpolicysecondfactorsrequest))
[ListLoginPolicySecondFactorsResponse](#listloginpolicysecondfactorsresponse)

Returns the available second factors defined by the administrators of ZITADEL



    POST: /policies/login/second_factors/_search


### AddSecondFactorToLoginPolicy

> **rpc** AddSecondFactorToLoginPolicy([AddSecondFactorToLoginPolicyRequest](#addsecondfactortologinpolicyrequest))
[AddSecondFactorToLoginPolicyResponse](#addsecondfactortologinpolicyresponse)

Adds a second factor to the default login policy.
It impacts all organisations without a customised policy



    POST: /policies/login/second_factors


### RemoveSecondFactorFromLoginPolicy

> **rpc** RemoveSecondFactorFromLoginPolicy([RemoveSecondFactorFromLoginPolicyRequest](#removesecondfactorfromloginpolicyrequest))
[RemoveSecondFactorFromLoginPolicyResponse](#removesecondfactorfromloginpolicyresponse)

Removes a second factor from the default login policy.
It impacts all organisations without a customised policy



    DELETE: /policies/login/second_factors/{type}


### ListLoginPolicyMultiFactors

> **rpc** ListLoginPolicyMultiFactors([ListLoginPolicyMultiFactorsRequest](#listloginpolicymultifactorsrequest))
[ListLoginPolicyMultiFactorsResponse](#listloginpolicymultifactorsresponse)

Returns the available multi factors defined by the administrators of ZITADEL



    POST: /policies/login/multi_factors/_search


### AddMultiFactorToLoginPolicy

> **rpc** AddMultiFactorToLoginPolicy([AddMultiFactorToLoginPolicyRequest](#addmultifactortologinpolicyrequest))
[AddMultiFactorToLoginPolicyResponse](#addmultifactortologinpolicyresponse)

Adds a multi factor to the default login policy.
It impacts all organisations without a customised policy



    POST: /policies/login/multi_factors


### RemoveMultiFactorFromLoginPolicy

> **rpc** RemoveMultiFactorFromLoginPolicy([RemoveMultiFactorFromLoginPolicyRequest](#removemultifactorfromloginpolicyrequest))
[RemoveMultiFactorFromLoginPolicyResponse](#removemultifactorfromloginpolicyresponse)

Removes a multi factor from the default login policy.
It impacts all organisations without a customised policy



    DELETE: /policies/login/multi_factors/{type}


### GetPasswordComplexityPolicy

> **rpc** GetPasswordComplexityPolicy([GetPasswordComplexityPolicyRequest](#getpasswordcomplexitypolicyrequest))
[GetPasswordComplexityPolicyResponse](#getpasswordcomplexitypolicyresponse)

Returns the password complexity policy defined by the administrators of ZITADEL



    GET: /policies/password/complexity


### UpdatePasswordComplexityPolicy

> **rpc** UpdatePasswordComplexityPolicy([UpdatePasswordComplexityPolicyRequest](#updatepasswordcomplexitypolicyrequest))
[UpdatePasswordComplexityPolicyResponse](#updatepasswordcomplexitypolicyresponse)

Updates the default password complexity policy of ZITADEL
it impacts all organisations without a customised policy



    PUT: /policies/password/complexity


### GetPasswordAgePolicy

> **rpc** GetPasswordAgePolicy([GetPasswordAgePolicyRequest](#getpasswordagepolicyrequest))
[GetPasswordAgePolicyResponse](#getpasswordagepolicyresponse)

Returns the password age policy defined by the administrators of ZITADEL



    GET: /policies/password/age


### UpdatePasswordAgePolicy

> **rpc** UpdatePasswordAgePolicy([UpdatePasswordAgePolicyRequest](#updatepasswordagepolicyrequest))
[UpdatePasswordAgePolicyResponse](#updatepasswordagepolicyresponse)

Updates the default password age policy of ZITADEL
it impacts all organisations without a customised policy



    PUT: /policies/password/age


### GetLockoutPolicy

> **rpc** GetLockoutPolicy([GetLockoutPolicyRequest](#getlockoutpolicyrequest))
[GetLockoutPolicyResponse](#getlockoutpolicyresponse)

Returns the lockout policy defined by the administrators of ZITADEL



    GET: /policies/lockout


### UpdateLockoutPolicy

> **rpc** UpdateLockoutPolicy([UpdateLockoutPolicyRequest](#updatelockoutpolicyrequest))
[UpdateLockoutPolicyResponse](#updatelockoutpolicyresponse)

Updates the default lockout policy of ZITADEL
it impacts all organisations without a customised policy



    PUT: /policies/password/lockout


### GetPrivacyPolicy

> **rpc** GetPrivacyPolicy([GetPrivacyPolicyRequest](#getprivacypolicyrequest))
[GetPrivacyPolicyResponse](#getprivacypolicyresponse)

Returns the privacy policy defined by the administrators of ZITADEL



    GET: /policies/privacy


### UpdatePrivacyPolicy

> **rpc** UpdatePrivacyPolicy([UpdatePrivacyPolicyRequest](#updateprivacypolicyrequest))
[UpdatePrivacyPolicyResponse](#updateprivacypolicyresponse)

Updates the default privacy policy of ZITADEL
it impacts all organisations without a customised policy
Variable {{.Lang}} can be set to have different links based on the language



    PUT: /policies/privacy


### AddNotificationPolicy

> **rpc** AddNotificationPolicy([AddNotificationPolicyRequest](#addnotificationpolicyrequest))
[AddNotificationPolicyResponse](#addnotificationpolicyresponse)

Add a default notification policy for ZITADEL
it impacts all organisations without a customised policy



    POST: /policies/notification


### GetNotificationPolicy

> **rpc** GetNotificationPolicy([GetNotificationPolicyRequest](#getnotificationpolicyrequest))
[GetNotificationPolicyResponse](#getnotificationpolicyresponse)

Returns the notification policy defined by the administrators of ZITADEL



    GET: /policies/notification


### UpdateNotificationPolicy

> **rpc** UpdateNotificationPolicy([UpdateNotificationPolicyRequest](#updatenotificationpolicyrequest))
[UpdateNotificationPolicyResponse](#updatenotificationpolicyresponse)

Updates the default notification policy of ZITADEL
it impacts all organisations without a customised policy



    PUT: /policies/notification


### GetDefaultInitMessageText

> **rpc** GetDefaultInitMessageText([GetDefaultInitMessageTextRequest](#getdefaultinitmessagetextrequest))
[GetDefaultInitMessageTextResponse](#getdefaultinitmessagetextresponse)

Returns the default text for initial message (translation file)



    GET: /text/default/message/init/{language}


### GetCustomInitMessageText

> **rpc** GetCustomInitMessageText([GetCustomInitMessageTextRequest](#getcustominitmessagetextrequest))
[GetCustomInitMessageTextResponse](#getcustominitmessagetextresponse)

Returns the custom text for initial message (overwritten in eventstore)



    GET: /text/message/init/{language}


### SetDefaultInitMessageText

> **rpc** SetDefaultInitMessageText([SetDefaultInitMessageTextRequest](#setdefaultinitmessagetextrequest))
[SetDefaultInitMessageTextResponse](#setdefaultinitmessagetextresponse)

Sets the default custom text for initial message
it impacts all organisations without customized initial message text
The Following Variables can be used:
{{.Code}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}} {{.CreationDate}}



    PUT: /text/message/init/{language}


### ResetCustomInitMessageTextToDefault

> **rpc** ResetCustomInitMessageTextToDefault([ResetCustomInitMessageTextToDefaultRequest](#resetcustominitmessagetexttodefaultrequest))
[ResetCustomInitMessageTextToDefaultResponse](#resetcustominitmessagetexttodefaultresponse)

Removes the custom init message text of the system
The default text from the translation file will trigger after



    DELETE: /text/message/init/{language}


### GetDefaultPasswordResetMessageText

> **rpc** GetDefaultPasswordResetMessageText([GetDefaultPasswordResetMessageTextRequest](#getdefaultpasswordresetmessagetextrequest))
[GetDefaultPasswordResetMessageTextResponse](#getdefaultpasswordresetmessagetextresponse)

Returns the default text for password reset message (translation file)



    GET: /text/deafult/message/passwordreset/{language}


### GetCustomPasswordResetMessageText

> **rpc** GetCustomPasswordResetMessageText([GetCustomPasswordResetMessageTextRequest](#getcustompasswordresetmessagetextrequest))
[GetCustomPasswordResetMessageTextResponse](#getcustompasswordresetmessagetextresponse)

Returns the custom text for password reset message (overwritten in eventstore)



    GET: /text/message/passwordreset/{language}


### SetDefaultPasswordResetMessageText

> **rpc** SetDefaultPasswordResetMessageText([SetDefaultPasswordResetMessageTextRequest](#setdefaultpasswordresetmessagetextrequest))
[SetDefaultPasswordResetMessageTextResponse](#setdefaultpasswordresetmessagetextresponse)

Sets the default custom text for password reset message
it impacts all organisations without customized password reset message text
The Following Variables can be used:
{{.Code}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}} {{.CreationDate}}



    PUT: /text/message/passwordreset/{language}


### ResetCustomPasswordResetMessageTextToDefault

> **rpc** ResetCustomPasswordResetMessageTextToDefault([ResetCustomPasswordResetMessageTextToDefaultRequest](#resetcustompasswordresetmessagetexttodefaultrequest))
[ResetCustomPasswordResetMessageTextToDefaultResponse](#resetcustompasswordresetmessagetexttodefaultresponse)

Removes the custom password reset message text of the system
The default text from the translation file will trigger after



    DELETE: /text/message/verifyemail/{language}


### GetDefaultVerifyEmailMessageText

> **rpc** GetDefaultVerifyEmailMessageText([GetDefaultVerifyEmailMessageTextRequest](#getdefaultverifyemailmessagetextrequest))
[GetDefaultVerifyEmailMessageTextResponse](#getdefaultverifyemailmessagetextresponse)

Returns the default text for verify email message (translation files)



    GET: /text/default/message/verifyemail/{language}


### GetCustomVerifyEmailMessageText

> **rpc** GetCustomVerifyEmailMessageText([GetCustomVerifyEmailMessageTextRequest](#getcustomverifyemailmessagetextrequest))
[GetCustomVerifyEmailMessageTextResponse](#getcustomverifyemailmessagetextresponse)

Returns the custom text for verify email message (overwritten in eventstore)



    GET: /text/message/verifyemail/{language}


### SetDefaultVerifyEmailMessageText

> **rpc** SetDefaultVerifyEmailMessageText([SetDefaultVerifyEmailMessageTextRequest](#setdefaultverifyemailmessagetextrequest))
[SetDefaultVerifyEmailMessageTextResponse](#setdefaultverifyemailmessagetextresponse)

Sets the default custom text for verify email message
it impacts all organisations without customized verify email message text
The Following Variables can be used:
{{.Code}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}} {{.CreationDate}}



    PUT: /text/message/verifyemail/{language}


### ResetCustomVerifyEmailMessageTextToDefault

> **rpc** ResetCustomVerifyEmailMessageTextToDefault([ResetCustomVerifyEmailMessageTextToDefaultRequest](#resetcustomverifyemailmessagetexttodefaultrequest))
[ResetCustomVerifyEmailMessageTextToDefaultResponse](#resetcustomverifyemailmessagetexttodefaultresponse)

Removes the custom verify email message text of the system
The default text from the translation file will trigger after



    DELETE: /text/message/verifyemail/{language}


### GetDefaultVerifyPhoneMessageText

> **rpc** GetDefaultVerifyPhoneMessageText([GetDefaultVerifyPhoneMessageTextRequest](#getdefaultverifyphonemessagetextrequest))
[GetDefaultVerifyPhoneMessageTextResponse](#getdefaultverifyphonemessagetextresponse)

Returns the default text for verify phone message (translation file)



    GET: /text/default/message/verifyphone/{language}


### GetCustomVerifyPhoneMessageText

> **rpc** GetCustomVerifyPhoneMessageText([GetCustomVerifyPhoneMessageTextRequest](#getcustomverifyphonemessagetextrequest))
[GetCustomVerifyPhoneMessageTextResponse](#getcustomverifyphonemessagetextresponse)

Returns the custom text for verify phone message



    GET: /text/message/verifyphone/{language}


### SetDefaultVerifyPhoneMessageText

> **rpc** SetDefaultVerifyPhoneMessageText([SetDefaultVerifyPhoneMessageTextRequest](#setdefaultverifyphonemessagetextrequest))
[SetDefaultVerifyPhoneMessageTextResponse](#setdefaultverifyphonemessagetextresponse)

Sets the default custom text for verify phone message
it impacts all organisations without customized verify phone message text
The Following Variables can be used:
{{.Code}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}} {{.CreationDate}}



    PUT: /text/message/verifyphone/{language}


### ResetCustomVerifyPhoneMessageTextToDefault

> **rpc** ResetCustomVerifyPhoneMessageTextToDefault([ResetCustomVerifyPhoneMessageTextToDefaultRequest](#resetcustomverifyphonemessagetexttodefaultrequest))
[ResetCustomVerifyPhoneMessageTextToDefaultResponse](#resetcustomverifyphonemessagetexttodefaultresponse)

Removes the custom verify phone text of the system
The default text from the translation file will trigger after



    DELETE: /text/message/verifyphone/{language}


### GetDefaultDomainClaimedMessageText

> **rpc** GetDefaultDomainClaimedMessageText([GetDefaultDomainClaimedMessageTextRequest](#getdefaultdomainclaimedmessagetextrequest))
[GetDefaultDomainClaimedMessageTextResponse](#getdefaultdomainclaimedmessagetextresponse)

Returns the default text for domain claimed message (translation file)



    GET: /text/default/message/domainclaimed/{language}


### GetCustomDomainClaimedMessageText

> **rpc** GetCustomDomainClaimedMessageText([GetCustomDomainClaimedMessageTextRequest](#getcustomdomainclaimedmessagetextrequest))
[GetCustomDomainClaimedMessageTextResponse](#getcustomdomainclaimedmessagetextresponse)

Returns the custom text for domain claimed message (overwritten in eventstore)



    GET: /text/message/domainclaimed/{language}


### SetDefaultDomainClaimedMessageText

> **rpc** SetDefaultDomainClaimedMessageText([SetDefaultDomainClaimedMessageTextRequest](#setdefaultdomainclaimedmessagetextrequest))
[SetDefaultDomainClaimedMessageTextResponse](#setdefaultdomainclaimedmessagetextresponse)

Sets the default custom text for domain claimed message
it impacts all organisations without customized domain claimed message text
The Following Variables can be used:
{{.Domain}} {{.TempUsername}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}} {{.CreationDate}}



    PUT: /text/message/domainclaimed/{language}


### ResetCustomDomainClaimedMessageTextToDefault

> **rpc** ResetCustomDomainClaimedMessageTextToDefault([ResetCustomDomainClaimedMessageTextToDefaultRequest](#resetcustomdomainclaimedmessagetexttodefaultrequest))
[ResetCustomDomainClaimedMessageTextToDefaultResponse](#resetcustomdomainclaimedmessagetexttodefaultresponse)

Removes the custom domain claimed message text of the system
The default text from the translation file will trigger after



    DELETE: /text/message/domainclaimed/{language}


### GetDefaultPasswordlessRegistrationMessageText

> **rpc** GetDefaultPasswordlessRegistrationMessageText([GetDefaultPasswordlessRegistrationMessageTextRequest](#getdefaultpasswordlessregistrationmessagetextrequest))
[GetDefaultPasswordlessRegistrationMessageTextResponse](#getdefaultpasswordlessregistrationmessagetextresponse)

Returns the default text for passwordless registration message (translation file)



    GET: /text/default/message/passwordless_registration/{language}


### GetCustomPasswordlessRegistrationMessageText

> **rpc** GetCustomPasswordlessRegistrationMessageText([GetCustomPasswordlessRegistrationMessageTextRequest](#getcustompasswordlessregistrationmessagetextrequest))
[GetCustomPasswordlessRegistrationMessageTextResponse](#getcustompasswordlessregistrationmessagetextresponse)

Returns the custom text for passwordless registration message (overwritten in eventstore)



    GET: /text/message/passwordless_registration/{language}


### SetDefaultPasswordlessRegistrationMessageText

> **rpc** SetDefaultPasswordlessRegistrationMessageText([SetDefaultPasswordlessRegistrationMessageTextRequest](#setdefaultpasswordlessregistrationmessagetextrequest))
[SetDefaultPasswordlessRegistrationMessageTextResponse](#setdefaultpasswordlessregistrationmessagetextresponse)

Sets the default custom text for passwordless registration message
it impacts all organisations without customized passwordless registration message text
The Following Variables can be used:
{{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}} {{.CreationDate}}



    PUT: /text/message/passwordless_registration/{language}


### ResetCustomPasswordlessRegistrationMessageTextToDefault

> **rpc** ResetCustomPasswordlessRegistrationMessageTextToDefault([ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest](#resetcustompasswordlessregistrationmessagetexttodefaultrequest))
[ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse](#resetcustompasswordlessregistrationmessagetexttodefaultresponse)

Removes the custom passwordless link message text of the system
The default text from the translation file will trigger after



    DELETE: /text/message/passwordless_registration/{language}


### GetDefaultPasswordChangeMessageText

> **rpc** GetDefaultPasswordChangeMessageText([GetDefaultPasswordChangeMessageTextRequest](#getdefaultpasswordchangemessagetextrequest))
[GetDefaultPasswordChangeMessageTextResponse](#getdefaultpasswordchangemessagetextresponse)

Returns the default text for password change message (translation file)



    GET: /text/default/message/password_change/{language}


### GetCustomPasswordChangeMessageText

> **rpc** GetCustomPasswordChangeMessageText([GetCustomPasswordChangeMessageTextRequest](#getcustompasswordchangemessagetextrequest))
[GetCustomPasswordChangeMessageTextResponse](#getcustompasswordchangemessagetextresponse)

Returns the custom text for password change message (overwritten in eventstore)



    GET: /text/message/password_change/{language}


### SetDefaultPasswordChangeMessageText

> **rpc** SetDefaultPasswordChangeMessageText([SetDefaultPasswordChangeMessageTextRequest](#setdefaultpasswordchangemessagetextrequest))
[SetDefaultPasswordChangeMessageTextResponse](#setdefaultpasswordchangemessagetextresponse)

Sets the default custom text for password change message
it impacts all organisations without customized password change message text
The Following Variables can be used:
{{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}} {{.CreationDate}}



    PUT: /text/message/password_change/{language}


### ResetCustomPasswordChangeMessageTextToDefault

> **rpc** ResetCustomPasswordChangeMessageTextToDefault([ResetCustomPasswordChangeMessageTextToDefaultRequest](#resetcustompasswordchangemessagetexttodefaultrequest))
[ResetCustomPasswordChangeMessageTextToDefaultResponse](#resetcustompasswordchangemessagetexttodefaultresponse)

Removes the custom password change message text of the system
The default text from the translation file will trigger after



    DELETE: /text/message/password_change/{language}


### GetDefaultLoginTexts

> **rpc** GetDefaultLoginTexts([GetDefaultLoginTextsRequest](#getdefaultlogintextsrequest))
[GetDefaultLoginTextsResponse](#getdefaultlogintextsresponse)

Returns the default custom texts for login ui (translation file)



    GET: /text/default/login/{language}


### GetCustomLoginTexts

> **rpc** GetCustomLoginTexts([GetCustomLoginTextsRequest](#getcustomlogintextsrequest))
[GetCustomLoginTextsResponse](#getcustomlogintextsresponse)

Returns the custom texts for login ui



    GET: /text/login/{language}


### SetCustomLoginText

> **rpc** SetCustomLoginText([SetCustomLoginTextsRequest](#setcustomlogintextsrequest))
[SetCustomLoginTextsResponse](#setcustomlogintextsresponse)

Sets the custom text for login ui
it impacts all organisations without customized login ui texts



    PUT: /text/login/{language}


### ResetCustomLoginTextToDefault

> **rpc** ResetCustomLoginTextToDefault([ResetCustomLoginTextsToDefaultRequest](#resetcustomlogintextstodefaultrequest))
[ResetCustomLoginTextsToDefaultResponse](#resetcustomlogintextstodefaultresponse)

Removes the custom texts for login ui
it impacts all organisations without customized login ui texts
The default text form translation file will trigger after



    DELETE: /text/login/{language}


### ListIAMMemberRoles

> **rpc** ListIAMMemberRoles([ListIAMMemberRolesRequest](#listiammemberrolesrequest))
[ListIAMMemberRolesResponse](#listiammemberrolesresponse)

Returns the IAM roles visible for the requested user



    POST: /members/roles/_search


### ListIAMMembers

> **rpc** ListIAMMembers([ListIAMMembersRequest](#listiammembersrequest))
[ListIAMMembersResponse](#listiammembersresponse)

Returns all members matching the request
all queries need to match (ANDed)



    POST: /members/_search


### AddIAMMember

> **rpc** AddIAMMember([AddIAMMemberRequest](#addiammemberrequest))
[AddIAMMemberResponse](#addiammemberresponse)

Adds a user to the membership list of ZITADEL with the given roles
undefined roles will be dropped



    POST: /members


### UpdateIAMMember

> **rpc** UpdateIAMMember([UpdateIAMMemberRequest](#updateiammemberrequest))
[UpdateIAMMemberResponse](#updateiammemberresponse)

Sets the given roles on a member.
The member has only roles provided by this call



    PUT: /members/{user_id}


### RemoveIAMMember

> **rpc** RemoveIAMMember([RemoveIAMMemberRequest](#removeiammemberrequest))
[RemoveIAMMemberResponse](#removeiammemberresponse)

Removes the user from the membership list of ZITADEL



    DELETE: /members/{user_id}


### ListViews

> **rpc** ListViews([ListViewsRequest](#listviewsrequest))
[ListViewsResponse](#listviewsresponse)

Returns all stored read models of ZITADEL
views are used for search optimisation and optimise request latencies
they represent the delta of the event happend on the objects



    POST: /views/_search


### ListFailedEvents

> **rpc** ListFailedEvents([ListFailedEventsRequest](#listfailedeventsrequest))
[ListFailedEventsResponse](#listfailedeventsresponse)

Returns event descriptions which cannot be processed.
It's possible that some events need some retries.
For example if the SMTP-API wasn't able to send an email at the first time



    POST: /failedevents/_search


### RemoveFailedEvent

> **rpc** RemoveFailedEvent([RemoveFailedEventRequest](#removefailedeventrequest))
[RemoveFailedEventResponse](#removefailedeventresponse)

Deletes the event from failed events view.
the event is not removed from the change stream
This call is usefull if the system was able to process the event later.
e.g. if the second try of sending an email was successful. the first try produced a
failed event. You can find out if it worked on the `failure_count`



    DELETE: /failedevents/{database}/{view_name}/{failed_sequence}


### ImportData

> **rpc** ImportData([ImportDataRequest](#importdatarequest))
[ImportDataResponse](#importdataresponse)

Imports data into instance and creates different objects



    POST: /import


### ExportData

> **rpc** ExportData([ExportDataRequest](#exportdatarequest))
[ExportDataResponse](#exportdataresponse)

Exports data from instance



    POST: /export


### ListEventTypes

> **rpc** ListEventTypes([ListEventTypesRequest](#listeventtypesrequest))
[ListEventTypesResponse](#listeventtypesresponse)





    POST: /events/types/_search


### ListEvents

> **rpc** ListEvents([ListEventsRequest](#listeventsrequest))
[ListEventsResponse](#listeventsresponse)





    POST: /events/_search


### ListAggregateTypes

> **rpc** ListAggregateTypes([ListAggregateTypesRequest](#listaggregatetypesrequest))
[ListAggregateTypesResponse](#listaggregatetypesresponse)





    POST: /aggregates/types/_search







## Messages


### ActivateLabelPolicyRequest
This is an empty request




### ActivateLabelPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ActivateSMSProviderRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ActivateSMSProviderResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddCustomDomainPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_login_must_be_domain |  bool | the username has to end with the domain of it's organisation (uniqueness is organisation based) |  |
| validate_org_domains |  bool | - |  |
| smtp_sender_address_matches_instance_domain |  bool | - |  |




### AddCustomDomainPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddCustomOrgIAMPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_login_must_be_domain |  bool | the username has to end with the domain of it's organisation (uniqueness is organisation based) |  |




### AddCustomOrgIAMPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddGoogleProviderRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.max_len: 200<br />  |
| client_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| client_secret |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| scopes | repeated string | - | repeated.max_items: 20<br /> repeated.items.string.min_len: 1<br /> repeated.items.string.max_len: 100<br />  |
| provider_options |  zitadel.idp.v1.Options | - |  |




### AddGoogleProviderResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| id |  string | - |  |




### AddIAMMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| roles | repeated string | if no roles provided the user won't have any rights |  |




### AddIAMMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddIDPToLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | Id of the predefined idp configuration | string.min_len: 1<br /> string.max_len: 200<br />  |




### AddIDPToLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddJWTIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| styling_type |  zitadel.idp.v1.IDPStylingType | - | enum.defined_only: true<br />  |
| jwt_endpoint |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| issuer |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| keys_endpoint |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| header_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| auto_register |  bool | - |  |




### AddJWTIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| idp_id |  string | - |  |




### AddLDAPProviderRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| host |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| port |  string | - | string.max_len: 5<br />  |
| tls |  bool | - |  |
| base_dn |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_object_class |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_unique_attribute |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| admin |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| password |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| attributes |  zitadel.idp.v1.LDAPAttributes | - |  |
| provider_options |  zitadel.idp.v1.Options | - |  |




### AddLDAPProviderResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| id |  string | - |  |




### AddMultiFactorToLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  zitadel.policy.v1.MultiFactorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### AddMultiFactorToLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddNotificationPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| password_change |  bool | - |  |




### AddNotificationPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddOIDCIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| styling_type |  zitadel.idp.v1.IDPStylingType | - | enum.defined_only: true<br />  |
| client_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| client_secret |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| issuer |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| scopes | repeated string | - |  |
| display_name_mapping |  zitadel.idp.v1.OIDCMappingField | - | enum.defined_only: true<br />  |
| username_mapping |  zitadel.idp.v1.OIDCMappingField | - | enum.defined_only: true<br />  |
| auto_register |  bool | - |  |




### AddOIDCIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| idp_id |  string | - |  |




### AddOIDCSettingsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| access_token_lifetime |  google.protobuf.Duration | - |  |
| id_token_lifetime |  google.protobuf.Duration | - |  |
| refresh_token_idle_expiration |  google.protobuf.Duration | - |  |
| refresh_token_expiration |  google.protobuf.Duration | - |  |




### AddOIDCSettingsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddSMSProviderTwilioRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| sid |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| token |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| sender_number |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### AddSMSProviderTwilioResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| id |  string | - |  |




### AddSMTPConfigRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| sender_address |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| sender_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| tls |  bool | - |  |
| host |  string | - | string.min_len: 1<br /> string.max_len: 500<br />  |
| user |  string | - |  |
| password |  string | - |  |




### AddSMTPConfigResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddSecondFactorToLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  zitadel.policy.v1.SecondFactorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### AddSecondFactorToLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### DataOrg



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - |  |
| org |  zitadel.management.v1.AddOrgRequest | - |  |
| domain_policy |  AddCustomDomainPolicyRequest | - |  |
| label_policy |  zitadel.management.v1.AddCustomLabelPolicyRequest | - |  |
| lockout_policy |  zitadel.management.v1.AddCustomLockoutPolicyRequest | - |  |
| login_policy |  zitadel.management.v1.AddCustomLoginPolicyRequest | - |  |
| password_complexity_policy |  zitadel.management.v1.AddCustomPasswordComplexityPolicyRequest | - |  |
| privacy_policy |  zitadel.management.v1.AddCustomPrivacyPolicyRequest | - |  |
| projects | repeated zitadel.v1.v1.DataProject | - |  |
| project_roles | repeated zitadel.management.v1.AddProjectRoleRequest | - |  |
| api_apps | repeated zitadel.v1.v1.DataAPIApplication | - |  |
| oidc_apps | repeated zitadel.v1.v1.DataOIDCApplication | - |  |
| human_users | repeated zitadel.v1.v1.DataHumanUser | - |  |
| machine_users | repeated zitadel.v1.v1.DataMachineUser | - |  |
| trigger_actions | repeated zitadel.management.v1.SetTriggerActionsRequest | - |  |
| actions | repeated zitadel.v1.v1.DataAction | - |  |
| project_grants | repeated zitadel.v1.v1.DataProjectGrant | - |  |
| user_grants | repeated zitadel.management.v1.AddUserGrantRequest | - |  |
| org_members | repeated zitadel.management.v1.AddOrgMemberRequest | - |  |
| project_members | repeated zitadel.management.v1.AddProjectMemberRequest | - |  |
| project_grant_members | repeated zitadel.management.v1.AddProjectGrantMemberRequest | - |  |
| user_metadata | repeated zitadel.management.v1.SetUserMetadataRequest | - |  |
| login_texts | repeated zitadel.management.v1.SetCustomLoginTextsRequest | - |  |
| init_messages | repeated zitadel.management.v1.SetCustomInitMessageTextRequest | - |  |
| password_reset_messages | repeated zitadel.management.v1.SetCustomPasswordResetMessageTextRequest | - |  |
| verify_email_messages | repeated zitadel.management.v1.SetCustomVerifyEmailMessageTextRequest | - |  |
| verify_phone_messages | repeated zitadel.management.v1.SetCustomVerifyPhoneMessageTextRequest | - |  |
| domain_claimed_messages | repeated zitadel.management.v1.SetCustomDomainClaimedMessageTextRequest | - |  |
| passwordless_registration_messages | repeated zitadel.management.v1.SetCustomPasswordlessRegistrationMessageTextRequest | - |  |
| oidc_idps | repeated zitadel.v1.v1.DataOIDCIDP | - |  |
| jwt_idps | repeated zitadel.v1.v1.DataJWTIDP | - |  |
| user_links | repeated zitadel.idp.v1.IDPUserLink | - |  |
| domains | repeated zitadel.org.v1.Domain | - |  |
| app_keys | repeated zitadel.v1.v1.DataAppKey | - |  |
| machine_keys | repeated zitadel.v1.v1.DataMachineKey | - |  |




### DeactivateIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### DeactivateIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### DeactivateSMSProviderRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### DeactivateSMSProviderResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### DeleteProviderRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### DeleteProviderResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ExportDataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_ids | repeated string | - |  |
| excluded_org_ids | repeated string | - |  |
| with_passwords |  bool | - |  |
| with_otp |  bool | - |  |
| response_output |  bool | - |  |
| local_output |  ExportDataRequest.LocalOutput | - |  |
| s3_output |  ExportDataRequest.S3Output | - |  |
| gcs_output |  ExportDataRequest.GCSOutput | - |  |
| timeout |  string | - |  |




### ExportDataRequest.GCSOutput



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| bucket |  string | - |  |
| serviceaccount_json |  string | - |  |
| path |  string | - |  |




### ExportDataRequest.LocalOutput



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| path |  string | - |  |




### ExportDataRequest.S3Output



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| path |  string | - |  |
| endpoint |  string | - |  |
| access_key_id |  string | - |  |
| secret_access_key |  string | - |  |
| ssl |  bool | - |  |
| bucket |  string | - |  |




### ExportDataResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| orgs | repeated DataOrg | - |  |




### FailedEvent



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| database |  string | - |  |
| view_name |  string | - |  |
| failed_sequence |  uint64 | - |  |
| failure_count |  uint64 | - |  |
| error_message |  string | - |  |
| last_failed |  google.protobuf.Timestamp | - |  |




### GetCustomDomainClaimedMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomDomainClaimedMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetCustomDomainPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomDomainPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.DomainPolicy | - |  |
| is_default |  bool | deprecated: is_default is also defined in zitadel.policy.v1.DomainPolicy |  |




### GetCustomInitMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomInitMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetCustomLoginTextsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomLoginTextsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.LoginCustomText | - |  |




### GetCustomOrgIAMPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomOrgIAMPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.OrgIAMPolicy | - |  |
| is_default |  bool | deprecated: is_default is also defined in zitadel.policy.v1.OrgIAMPolicy |  |




### GetCustomPasswordChangeMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomPasswordChangeMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetCustomPasswordResetMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomPasswordResetMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetCustomPasswordlessRegistrationMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomPasswordlessRegistrationMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetCustomVerifyEmailMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomVerifyEmailMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetCustomVerifyPhoneMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomVerifyPhoneMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultDomainClaimedMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultDomainClaimedMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultInitMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultInitMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultLanguageRequest
This is an empty request




### GetDefaultLanguageResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - |  |




### GetDefaultLoginTextsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultLoginTextsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.LoginCustomText | - |  |




### GetDefaultOrgRequest
This is an empty request




### GetDefaultOrgResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org |  zitadel.org.v1.Org | - |  |




### GetDefaultPasswordChangeMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultPasswordChangeMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultPasswordResetMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultPasswordResetMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultPasswordlessRegistrationMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultPasswordlessRegistrationMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultVerifyEmailMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultVerifyEmailMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultVerifyPhoneMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultVerifyPhoneMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDomainPolicyRequest





### GetDomainPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.DomainPolicy | - |  |




### GetFileSystemNotificationProviderRequest
This is an empty request




### GetFileSystemNotificationProviderResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| provider |  zitadel.settings.v1.DebugNotificationProvider | - |  |




### GetIDPByIDRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetIDPByIDResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp |  zitadel.idp.v1.IDP | - |  |




### GetLabelPolicyRequest
This is an empty request




### GetLabelPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.LabelPolicy | - |  |




### GetLockoutPolicyRequest
This is an empty request




### GetLockoutPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.LockoutPolicy | - |  |




### GetLogNotificationProviderRequest
This is an empty request




### GetLogNotificationProviderResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| provider |  zitadel.settings.v1.DebugNotificationProvider | - |  |




### GetLoginPolicyRequest
This is an empty request




### GetLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.LoginPolicy | - |  |




### GetMyInstanceRequest
This is an empty request




### GetMyInstanceResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| instance |  zitadel.instance.v1.InstanceDetail | - |  |




### GetNotificationPolicyRequest
This is an empty request




### GetNotificationPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.NotificationPolicy | - |  |




### GetOIDCSettingsRequest
This is an empty request




### GetOIDCSettingsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| settings |  zitadel.settings.v1.OIDCSettings | - |  |




### GetOrgByIDRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetOrgByIDResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org |  zitadel.org.v1.Org | - |  |




### GetOrgIAMPolicyRequest





### GetOrgIAMPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.OrgIAMPolicy | - |  |




### GetPasswordAgePolicyRequest
This is an empty request




### GetPasswordAgePolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PasswordAgePolicy | - |  |




### GetPasswordComplexityPolicyRequest





### GetPasswordComplexityPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PasswordComplexityPolicy | - |  |




### GetPreviewLabelPolicyRequest
This is an empty request




### GetPreviewLabelPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.LabelPolicy | - |  |




### GetPrivacyPolicyRequest
This is an empty request




### GetPrivacyPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PrivacyPolicy | - |  |




### GetProviderByIDRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetProviderByIDResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp |  zitadel.idp.v1.Provider | - |  |




### GetSMSProviderRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 100<br />  |




### GetSMSProviderResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| config |  zitadel.settings.v1.SMSProvider | - |  |




### GetSMTPConfigRequest
This is an empty request




### GetSMTPConfigResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| smtp_config |  zitadel.settings.v1.SMTPConfig | - |  |




### GetSecretGeneratorRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| generator_type |  zitadel.settings.v1.SecretGeneratorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### GetSecretGeneratorResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| secret_generator |  zitadel.settings.v1.SecretGenerator | - |  |




### GetSecurityPolicyRequest
This is an empty request




### GetSecurityPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.settings.v1.SecurityPolicy | - |  |




### GetSupportedLanguagesRequest
This is an empty request




### GetSupportedLanguagesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| languages | repeated string | - |  |




### HealthzRequest
This is an empty request




### HealthzResponse
This is an empty response




### IDPQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_id_query |  zitadel.idp.v1.IDPIDQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_name_query |  zitadel.idp.v1.IDPNameQuery | - |  |




### ImportDataError



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  string | - |  |
| id |  string | - |  |
| message |  string | - |  |




### ImportDataOrg



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| orgs | repeated DataOrg | - |  |




### ImportDataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) data.data_orgs |  ImportDataOrg | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) data.data_orgsv1 |  zitadel.v1.v1.ImportDataOrg | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) data.data_orgs_local |  ImportDataRequest.LocalInput | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) data.data_orgsv1_local |  ImportDataRequest.LocalInput | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) data.data_orgs_s3 |  ImportDataRequest.S3Input | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) data.data_orgsv1_s3 |  ImportDataRequest.S3Input | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) data.data_orgs_gcs |  ImportDataRequest.GCSInput | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) data.data_orgsv1_gcs |  ImportDataRequest.GCSInput | - |  |
| timeout |  string | - |  |




### ImportDataRequest.GCSInput



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| bucket |  string | - |  |
| serviceaccount_json |  string | - |  |
| path |  string | - |  |




### ImportDataRequest.LocalInput



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| path |  string | - |  |




### ImportDataRequest.S3Input



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| path |  string | - |  |
| endpoint |  string | - |  |
| access_key_id |  string | - |  |
| secret_access_key |  string | - |  |
| ssl |  bool | - |  |
| bucket |  string | - |  |




### ImportDataResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| errors | repeated ImportDataError | - |  |
| success |  ImportDataSuccess | - |  |




### ImportDataSuccess



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| orgs | repeated ImportDataSuccessOrg | - |  |




### ImportDataSuccessOrg



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - |  |
| project_ids | repeated string | - |  |
| project_roles | repeated string | - |  |
| oidc_app_ids | repeated string | - |  |
| api_app_ids | repeated string | - |  |
| human_user_ids | repeated string | - |  |
| machine_user_ids | repeated string | - |  |
| action_ids | repeated string | - |  |
| trigger_actions | repeated zitadel.management.v1.SetTriggerActionsRequest | - |  |
| project_grants | repeated ImportDataSuccessProjectGrant | - |  |
| user_grants | repeated ImportDataSuccessUserGrant | - |  |
| org_members | repeated string | - |  |
| project_members | repeated ImportDataSuccessProjectMember | - |  |
| project_grant_members | repeated ImportDataSuccessProjectGrantMember | - |  |
| oidc_ipds | repeated string | - |  |
| jwt_idps | repeated string | - |  |
| idp_links | repeated string | - |  |
| user_links | repeated ImportDataSuccessUserLinks | - |  |
| user_metadata | repeated ImportDataSuccessUserMetadata | - |  |
| domains | repeated string | - |  |
| app_keys | repeated string | - |  |
| machine_keys | repeated string | - |  |




### ImportDataSuccessProjectGrant



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| grant_id |  string | - |  |
| project_id |  string | - |  |
| org_id |  string | - |  |




### ImportDataSuccessProjectGrantMember



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - |  |
| grant_id |  string | - |  |
| user_id |  string | - |  |




### ImportDataSuccessProjectMember



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - |  |
| user_id |  string | - |  |




### ImportDataSuccessUserGrant



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| project_id |  string | - |  |
| user_id |  string | - |  |




### ImportDataSuccessUserLinks



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - |  |
| external_user_id |  string | - |  |
| display_name |  string | - |  |
| idp_id |  string | - |  |




### ImportDataSuccessUserMetadata



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - |  |
| key |  string | - |  |




### IsOrgUniqueRequest
if name or domain is already in use, org is not unique
at least one argument has to be provided


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.max_len: 200<br />  |
| domain |  string | - | string.max_len: 200<br />  |




### IsOrgUniqueResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| is_unique |  bool | - |  |




### ListAggregateTypesRequest





### ListAggregateTypesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| aggregate_types | repeated zitadel.event.v1.AggregateType | - |  |




### ListEventTypesRequest





### ListEventTypesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| event_types | repeated zitadel.event.v1.EventType | - |  |




### ListEventsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| sequence |  uint64 | sequence represents the order of events. It's always upcounting if asc is false sequence is used as less than filter if asc is true sequence is used as greater than filter if sequence is 0 the field is ignored |  |
| limit |  uint32 | - |  |
| asc |  bool | - |  |
| editor_user_id |  string | - | string.min_len: 0<br /> string.max_len: 200<br />  |
| event_types | repeated string | the types are or filtered and must match the type exatly | repeated.max_items: 30<br />  |
| aggregate_id |  string | - | string.min_len: 0<br /> string.max_len: 200<br />  |
| aggregate_types | repeated string | - | repeated.max_items: 10<br />  |
| resource_owner |  string | - | string.min_len: 0<br /> string.max_len: 200<br />  |
| creation_date |  google.protobuf.Timestamp | if asc is false creation_date is used as less than filter if asc is true creation_date is used as greater than filter if creation_date is not set the field is ignored |  |




### ListEventsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| events | repeated zitadel.event.v1.Event | - |  |




### ListFailedEventsRequest
This is an empty request




### ListFailedEventsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated FailedEvent | TODO: list details |  |




### ListIAMMemberRolesRequest
This is an empty request




### ListIAMMemberRolesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| roles | repeated string | - |  |




### ListIAMMembersRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.member.v1.SearchQuery | criterias the client is looking for |  |




### ListIAMMembersResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.member.v1.Member | - |  |




### ListIDPsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| sorting_column |  zitadel.idp.v1.IDPFieldName | the field the result is sorted |  |
| queries | repeated IDPQuery | criterias the client is looking for |  |




### ListIDPsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| sorting_column |  zitadel.idp.v1.IDPFieldName | - |  |
| result | repeated zitadel.idp.v1.IDP | - |  |




### ListInstanceDomainsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | - |  |
| sorting_column |  zitadel.instance.v1.DomainFieldName | the field the result is sorted |  |
| queries | repeated zitadel.instance.v1.DomainSearchQuery | criterias the client is looking for |  |




### ListInstanceDomainsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| sorting_column |  zitadel.instance.v1.DomainFieldName | - |  |
| result | repeated zitadel.instance.v1.Domain | - |  |




### ListLoginPolicyIDPsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |




### ListLoginPolicyIDPsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.idp.v1.IDPLoginPolicyLink | - |  |




### ListLoginPolicyMultiFactorsRequest
This is an empty request




### ListLoginPolicyMultiFactorsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.policy.v1.MultiFactorType | - |  |




### ListLoginPolicySecondFactorsRequest
This is an empty request




### ListLoginPolicySecondFactorsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.policy.v1.SecondFactorType | - |  |




### ListOrgsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| sorting_column |  zitadel.org.v1.OrgFieldName | the field the result is sorted |  |
| queries | repeated zitadel.org.v1.OrgQuery | criterias the client is looking for |  |




### ListOrgsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| sorting_column |  zitadel.org.v1.OrgFieldName | - |  |
| result | repeated zitadel.org.v1.Org | - |  |




### ListProvidersRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated ProviderQuery | criteria the client is looking for |  |




### ListProvidersResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.idp.v1.Provider | - |  |




### ListSMSProvidersRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |




### ListSMSProvidersResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.settings.v1.SMSProvider | - |  |




### ListSecretGeneratorsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.settings.v1.SecretGeneratorQuery | criterias the client is looking for |  |




### ListSecretGeneratorsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.settings.v1.SecretGenerator | - |  |




### ListViewsRequest
This is an empty request




### ListViewsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated View | TODO: list details |  |




### ProviderQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_id_query |  zitadel.idp.v1.IDPIDQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_name_query |  zitadel.idp.v1.IDPNameQuery | - |  |




### ReactivateIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ReactivateIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveFailedEventRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| database |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| view_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| failed_sequence |  uint64 | - |  |




### RemoveFailedEventResponse
This is an empty response




### RemoveIAMMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveIAMMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveIDPFromLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveIDPFromLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveLabelPolicyFontRequest
This is an empty request




### RemoveLabelPolicyFontResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveLabelPolicyIconDarkRequest
This is an empty request




### RemoveLabelPolicyIconDarkResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveLabelPolicyIconRequest
This is an empty request




### RemoveLabelPolicyIconResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveLabelPolicyLogoDarkRequest
This is an empty request




### RemoveLabelPolicyLogoDarkResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveLabelPolicyLogoRequest
This is an empty request




### RemoveLabelPolicyLogoResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveMultiFactorFromLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  zitadel.policy.v1.MultiFactorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### RemoveMultiFactorFromLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveOrgRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveOrgResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveSMSProviderRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveSMSProviderResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveSMTPConfigRequest
this is en empty request




### RemoveSMTPConfigResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveSecondFactorFromLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  zitadel.policy.v1.SecondFactorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### RemoveSecondFactorFromLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomDomainClaimedMessageTextToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomDomainClaimedMessageTextToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomDomainPolicyToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomDomainPolicyToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomInitMessageTextToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomInitMessageTextToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomLoginTextsToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomLoginTextsToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomOrgIAMPolicyToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomOrgIAMPolicyToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomPasswordChangeMessageTextToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomPasswordChangeMessageTextToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomPasswordResetMessageTextToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomPasswordResetMessageTextToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomVerifyEmailMessageTextToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomVerifyEmailMessageTextToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomVerifyPhoneMessageTextToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomVerifyPhoneMessageTextToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetCustomLoginTextsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| select_account_text |  zitadel.text.v1.SelectAccountScreenText | - |  |
| login_text |  zitadel.text.v1.LoginScreenText | - |  |
| password_text |  zitadel.text.v1.PasswordScreenText | - |  |
| username_change_text |  zitadel.text.v1.UsernameChangeScreenText | - |  |
| username_change_done_text |  zitadel.text.v1.UsernameChangeDoneScreenText | - |  |
| init_password_text |  zitadel.text.v1.InitPasswordScreenText | - |  |
| init_password_done_text |  zitadel.text.v1.InitPasswordDoneScreenText | - |  |
| email_verification_text |  zitadel.text.v1.EmailVerificationScreenText | - |  |
| email_verification_done_text |  zitadel.text.v1.EmailVerificationDoneScreenText | - |  |
| initialize_user_text |  zitadel.text.v1.InitializeUserScreenText | - |  |
| initialize_done_text |  zitadel.text.v1.InitializeUserDoneScreenText | - |  |
| init_mfa_prompt_text |  zitadel.text.v1.InitMFAPromptScreenText | - |  |
| init_mfa_otp_text |  zitadel.text.v1.InitMFAOTPScreenText | - |  |
| init_mfa_u2f_text |  zitadel.text.v1.InitMFAU2FScreenText | - |  |
| init_mfa_done_text |  zitadel.text.v1.InitMFADoneScreenText | - |  |
| mfa_providers_text |  zitadel.text.v1.MFAProvidersText | - |  |
| verify_mfa_otp_text |  zitadel.text.v1.VerifyMFAOTPScreenText | - |  |
| verify_mfa_u2f_text |  zitadel.text.v1.VerifyMFAU2FScreenText | - |  |
| passwordless_text |  zitadel.text.v1.PasswordlessScreenText | - |  |
| password_change_text |  zitadel.text.v1.PasswordChangeScreenText | - |  |
| password_change_done_text |  zitadel.text.v1.PasswordChangeDoneScreenText | - |  |
| password_reset_done_text |  zitadel.text.v1.PasswordResetDoneScreenText | - |  |
| registration_option_text |  zitadel.text.v1.RegistrationOptionScreenText | - |  |
| registration_user_text |  zitadel.text.v1.RegistrationUserScreenText | - |  |
| registration_org_text |  zitadel.text.v1.RegistrationOrgScreenText | - |  |
| linking_user_done_text |  zitadel.text.v1.LinkingUserDoneScreenText | - |  |
| external_user_not_found_text |  zitadel.text.v1.ExternalUserNotFoundScreenText | - |  |
| success_login_text |  zitadel.text.v1.SuccessLoginScreenText | - |  |
| logout_text |  zitadel.text.v1.LogoutDoneScreenText | - |  |
| footer_text |  zitadel.text.v1.FooterText | - |  |
| passwordless_prompt_text |  zitadel.text.v1.PasswordlessPromptScreenText | - |  |
| passwordless_registration_text |  zitadel.text.v1.PasswordlessRegistrationScreenText | - |  |
| passwordless_registration_done_text |  zitadel.text.v1.PasswordlessRegistrationDoneScreenText | - |  |
| external_registration_user_overview_text |  zitadel.text.v1.ExternalRegistrationUserOverviewScreenText | - |  |




### SetCustomLoginTextsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultDomainClaimedMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetDefaultDomainClaimedMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultInitMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 1000<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetDefaultInitMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultLanguageRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 10<br />  |




### SetDefaultLanguageResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultOrgRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### SetDefaultOrgResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultPasswordChangeMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetDefaultPasswordChangeMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultPasswordResetMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetDefaultPasswordResetMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultPasswordlessRegistrationMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetDefaultPasswordlessRegistrationMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultVerifyEmailMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetDefaultVerifyEmailMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultVerifyPhoneMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetDefaultVerifyPhoneMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetSecurityPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| enable_iframe_embedding |  bool | states if iframe embedding is enabled or disabled |  |
| allowed_origins | repeated string | origins allowed to load ZITADEL in an iframe if enable_iframe_embedding is true |  |




### SetSecurityPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetUpOrgRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org |  SetUpOrgRequest.Org | - | message.required: true<br />  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) user.human |  SetUpOrgRequest.Human | oneof field for the user managing the organisation |  |
| roles | repeated string | specify Org Member Roles for the provided user (default is ORG_OWNER if roles are empty) |  |




### SetUpOrgRequest.Human



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| profile |  SetUpOrgRequest.Human.Profile | - | message.required: true<br />  |
| email |  SetUpOrgRequest.Human.Email | - | message.required: true<br />  |
| phone |  SetUpOrgRequest.Human.Phone | - |  |
| password |  string | - |  |




### SetUpOrgRequest.Human.Email



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| email |  string | - | string.email: true<br />  |
| is_email_verified |  bool | - |  |




### SetUpOrgRequest.Human.Phone



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| phone |  string | has to be a global number | string.min_len: 1<br /> string.max_len: 50<br /> string.prefix: +<br />  |
| is_phone_verified |  bool | - |  |




### SetUpOrgRequest.Human.Profile



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| first_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| last_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| nick_name |  string | - | string.max_len: 200<br />  |
| display_name |  string | - | string.max_len: 200<br />  |
| preferred_language |  string | - | string.max_len: 10<br />  |
| gender |  zitadel.user.v1.Gender | - |  |




### SetUpOrgRequest.Org



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| domain |  string | - | string.max_len: 200<br />  |




### SetUpOrgResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| org_id |  string | - |  |
| user_id |  string | - |  |




### UpdateCustomDomainPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_login_must_be_domain |  bool | - |  |
| validate_org_domains |  bool | - |  |
| smtp_sender_address_matches_instance_domain |  bool | - |  |




### UpdateCustomDomainPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateCustomOrgIAMPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_login_must_be_domain |  bool | - |  |




### UpdateCustomOrgIAMPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateDomainPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_login_must_be_domain |  bool | - |  |
| validate_org_domains |  bool | - |  |
| smtp_sender_address_matches_instance_domain |  bool | - |  |




### UpdateDomainPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateGoogleProviderRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| name |  string | - | string.max_len: 200<br />  |
| client_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| client_secret |  string | - | string.max_len: 200<br />  |
| scopes | repeated string | - | repeated.max_items: 20<br /> repeated.items.string.min_len: 1<br /> repeated.items.string.max_len: 100<br />  |
| provider_options |  zitadel.idp.v1.Options | - |  |




### UpdateGoogleProviderResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateIAMMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| roles | repeated string | if no roles provided the user won't have any rights |  |




### UpdateIAMMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateIDPJWTConfigRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| jwt_endpoint |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| issuer |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| keys_endpoint |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| header_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### UpdateIDPJWTConfigResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateIDPOIDCConfigRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| issuer |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| client_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| client_secret |  string | - | string.max_len: 200<br />  |
| scopes | repeated string | - |  |
| display_name_mapping |  zitadel.idp.v1.OIDCMappingField | - | enum.defined_only: true<br />  |
| username_mapping |  zitadel.idp.v1.OIDCMappingField | - | enum.defined_only: true<br />  |




### UpdateIDPOIDCConfigResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| styling_type |  zitadel.idp.v1.IDPStylingType | - | enum.defined_only: true<br />  |
| auto_register |  bool | - |  |




### UpdateIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateLDAPProviderRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| host |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| port |  string | - | string.max_len: 5<br />  |
| tls |  bool | - |  |
| base_dn |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_object_class |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_unique_attribute |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| admin |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| password |  string | - | string.max_len: 200<br />  |
| attributes |  zitadel.idp.v1.LDAPAttributes | - |  |
| provider_options |  zitadel.idp.v1.Options | - |  |




### UpdateLDAPProviderResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateLabelPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| primary_color |  string | - | string.max_len: 50<br />  |
| hide_login_name_suffix |  bool | - |  |
| warn_color |  string | - | string.max_len: 50<br />  |
| background_color |  string | - | string.max_len: 50<br />  |
| font_color |  string | - | string.max_len: 50<br />  |
| primary_color_dark |  string | - | string.max_len: 50<br />  |
| background_color_dark |  string | - | string.max_len: 50<br />  |
| warn_color_dark |  string | - | string.max_len: 50<br />  |
| font_color_dark |  string | - | string.max_len: 50<br />  |
| disable_watermark |  bool | - |  |




### UpdateLabelPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateLockoutPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| max_password_attempts |  uint32 | failed attempts until a user gets locked |  |




### UpdateLockoutPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| allow_username_password |  bool | - |  |
| allow_register |  bool | - |  |
| allow_external_idp |  bool | - |  |
| force_mfa |  bool | - |  |
| passwordless_type |  zitadel.policy.v1.PasswordlessType | - | enum.defined_only: true<br />  |
| hide_password_reset |  bool | - |  |
| ignore_unknown_usernames |  bool | - |  |
| default_redirect_uri |  string | - |  |
| password_check_lifetime |  google.protobuf.Duration | - |  |
| external_login_check_lifetime |  google.protobuf.Duration | - |  |
| mfa_init_skip_lifetime |  google.protobuf.Duration | - |  |
| second_factor_check_lifetime |  google.protobuf.Duration | - |  |
| multi_factor_check_lifetime |  google.protobuf.Duration | - |  |
| allow_domain_discovery |  bool | If set to true, the suffix (@domain.com) of an unknown username input on the login screen will be matched against the org domains and will redirect to the registration of that organisation on success. |  |
| disable_login_with_email |  bool | - |  |
| disable_login_with_phone |  bool | - |  |




### UpdateLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateNotificationPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| password_change |  bool | - |  |




### UpdateNotificationPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateOIDCSettingsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| access_token_lifetime |  google.protobuf.Duration | - |  |
| id_token_lifetime |  google.protobuf.Duration | - |  |
| refresh_token_idle_expiration |  google.protobuf.Duration | - |  |
| refresh_token_expiration |  google.protobuf.Duration | - |  |




### UpdateOIDCSettingsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateOrgIAMPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_login_must_be_domain |  bool | - |  |




### UpdateOrgIAMPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdatePasswordAgePolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| max_age_days |  uint32 | - |  |
| expire_warn_days |  uint32 | - |  |




### UpdatePasswordAgePolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdatePasswordComplexityPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| min_length |  uint32 | - |  |
| has_uppercase |  bool | - |  |
| has_lowercase |  bool | - |  |
| has_number |  bool | - |  |
| has_symbol |  bool | - |  |




### UpdatePasswordComplexityPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdatePrivacyPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| tos_link |  string | - |  |
| privacy_link |  string | - |  |
| help_link |  string | - |  |




### UpdatePrivacyPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateSMSProviderTwilioRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| sid |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| sender_number |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### UpdateSMSProviderTwilioResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateSMSProviderTwilioTokenRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| token |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### UpdateSMSProviderTwilioTokenResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateSMTPConfigPasswordRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| password |  string | - |  |




### UpdateSMTPConfigPasswordResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateSMTPConfigRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| sender_address |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| sender_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| tls |  bool | - |  |
| host |  string | - | string.min_len: 1<br /> string.max_len: 500<br />  |
| user |  string | - |  |




### UpdateSMTPConfigResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateSecretGeneratorRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| generator_type |  zitadel.settings.v1.SecretGeneratorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |
| length |  uint32 | - |  |
| expiry |  google.protobuf.Duration | - |  |
| include_lower_letters |  bool | - |  |
| include_upper_letters |  bool | - |  |
| include_digits |  bool | - |  |
| include_symbols |  bool | - |  |




### UpdateSecretGeneratorResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### View



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| database |  string | - |  |
| view_name |  string | - |  |
| processed_sequence |  uint64 | - |  |
| event_timestamp |  google.protobuf.Timestamp | The timestamp the event occurred |  |
| last_successful_spooler_run |  google.protobuf.Timestamp | - |  |






