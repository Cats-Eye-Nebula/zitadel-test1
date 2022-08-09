---
title: zitadel/policy.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)




## Messages


### DomainPolicy



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| user_login_must_be_domain |  bool | - |  |
| is_default |  bool | - |  |
| validate_org_domains |  bool | - |  |
| smtp_sender_address_matches_instance_domain |  bool | - |  |




### LabelPolicy



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| primary_color |  string | hex value for primary color |  |
| is_default |  bool | defines if the organisation's admin changed the policy |  |
| hide_login_name_suffix |  bool | hides the org suffix on the login form if the scope \"urn:zitadel:iam:org:domain:primary:{domainname}\" is set. Details about this [scope in](../openidoauth/scopes) |  |
| warn_color |  string | hex value for secondary color |  |
| background_color |  string | hex value for background color |  |
| font_color |  string | hex value for font color |  |
| primary_color_dark |  string | hex value for primary color dark theme |  |
| background_color_dark |  string | hex value for background color dark theme |  |
| warn_color_dark |  string | hex value for warn color dark theme |  |
| font_color_dark |  string | hex value for font color dark theme |  |
| disable_watermark |  bool | - |  |
| logo_url |  string | - |  |
| icon_url |  string | - |  |
| logo_url_dark |  string | - |  |
| icon_url_dark |  string | - |  |
| font_url |  string | - |  |




### LockoutPolicy



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| max_password_attempts |  uint64 | - |  |
| is_default |  bool | - |  |




### LoginPolicy



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| allow_username_password |  bool | - |  |
| allow_register |  bool | - |  |
| allow_external_idp |  bool | - |  |
| force_mfa |  bool | - |  |
| passwordless_type |  PasswordlessType | - |  |
| is_default |  bool | - |  |
| hide_password_reset |  bool | - |  |
| ignore_unknown_usernames |  bool | - |  |
| default_redirect_uri |  string | - |  |
| password_check_lifetime |  google.protobuf.Duration | - |  |
| external_login_check_lifetime |  google.protobuf.Duration | - |  |
| mfa_init_skip_lifetime |  google.protobuf.Duration | - |  |
| second_factor_check_lifetime |  google.protobuf.Duration | - |  |
| multi_factor_check_lifetime |  google.protobuf.Duration | - |  |
| second_factors | repeated SecondFactorType | - |  |
| multi_factors | repeated MultiFactorType | - |  |
| idps | repeated zitadel.idp.v1.IDPLoginPolicyLink | - |  |




### OrgIAMPolicy
deprecated: please use DomainPolicy instead


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| user_login_must_be_domain |  bool | - |  |
| is_default |  bool | - |  |




### PasswordAgePolicy



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| max_age_days |  uint64 | - |  |
| expire_warn_days |  uint64 | - |  |
| is_default |  bool | - |  |




### PasswordComplexityPolicy



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| min_length |  uint64 | - |  |
| has_uppercase |  bool | - |  |
| has_lowercase |  bool | - |  |
| has_number |  bool | - |  |
| has_symbol |  bool | - |  |
| is_default |  bool | - |  |




### PrivacyPolicy



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| tos_link |  string | - |  |
| privacy_link |  string | - |  |
| is_default |  bool | - |  |
| help_link |  string | - |  |






## Enums


### MultiFactorType {#multifactortype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| MULTI_FACTOR_TYPE_UNSPECIFIED | 0 | - |
| MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION | 1 | - |




### PasswordlessType {#passwordlesstype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| PASSWORDLESS_TYPE_NOT_ALLOWED | 0 | - |
| PASSWORDLESS_TYPE_ALLOWED | 1 | PLANNED: PASSWORDLESS_TYPE_WITH_CERT |




### SecondFactorType {#secondfactortype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| SECOND_FACTOR_TYPE_UNSPECIFIED | 0 | - |
| SECOND_FACTOR_TYPE_OTP | 1 | - |
| SECOND_FACTOR_TYPE_U2F | 2 | - |




