// @generated by protoc-gen-es v1.0.0
// @generated from file zitadel/policy.proto (package zitadel.policy.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, Duration, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import type { ObjectDetails } from "./object_pb.js";
import type { IDPLoginPolicyLink } from "./idp_pb.js";

/**
 * @generated from enum zitadel.policy.v1.SecondFactorType
 */
export declare enum SecondFactorType {
  /**
   * @generated from enum value: SECOND_FACTOR_TYPE_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: SECOND_FACTOR_TYPE_OTP = 1;
   */
  OTP = 1,

  /**
   * @generated from enum value: SECOND_FACTOR_TYPE_U2F = 2;
   */
  U2F = 2,
}

/**
 * @generated from enum zitadel.policy.v1.MultiFactorType
 */
export declare enum MultiFactorType {
  /**
   * @generated from enum value: MULTI_FACTOR_TYPE_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION = 1;
   */
  U2F_WITH_VERIFICATION = 1,
}

/**
 * @generated from enum zitadel.policy.v1.PasswordlessType
 */
export declare enum PasswordlessType {
  /**
   * @generated from enum value: PASSWORDLESS_TYPE_NOT_ALLOWED = 0;
   */
  NOT_ALLOWED = 0,

  /**
   * PLANNED: PASSWORDLESS_TYPE_WITH_CERT
   *
   * @generated from enum value: PASSWORDLESS_TYPE_ALLOWED = 1;
   */
  ALLOWED = 1,
}

/**
 * deprecated: please use DomainPolicy instead
 *
 * @generated from message zitadel.policy.v1.OrgIAMPolicy
 */
export declare class OrgIAMPolicy extends Message<OrgIAMPolicy> {
  /**
   * @generated from field: zitadel.v1.ObjectDetails details = 1;
   */
  details?: ObjectDetails;

  /**
   * @generated from field: bool user_login_must_be_domain = 2;
   */
  userLoginMustBeDomain: boolean;

  /**
   * @generated from field: bool is_default = 3;
   */
  isDefault: boolean;

  constructor(data?: PartialMessage<OrgIAMPolicy>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.policy.v1.OrgIAMPolicy";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): OrgIAMPolicy;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): OrgIAMPolicy;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): OrgIAMPolicy;

  static equals(a: OrgIAMPolicy | PlainMessage<OrgIAMPolicy> | undefined, b: OrgIAMPolicy | PlainMessage<OrgIAMPolicy> | undefined): boolean;
}

/**
 * @generated from message zitadel.policy.v1.DomainPolicy
 */
export declare class DomainPolicy extends Message<DomainPolicy> {
  /**
   * @generated from field: zitadel.v1.ObjectDetails details = 1;
   */
  details?: ObjectDetails;

  /**
   * @generated from field: bool user_login_must_be_domain = 2;
   */
  userLoginMustBeDomain: boolean;

  /**
   * @generated from field: bool is_default = 3;
   */
  isDefault: boolean;

  /**
   * @generated from field: bool validate_org_domains = 4;
   */
  validateOrgDomains: boolean;

  /**
   * @generated from field: bool smtp_sender_address_matches_instance_domain = 5;
   */
  smtpSenderAddressMatchesInstanceDomain: boolean;

  constructor(data?: PartialMessage<DomainPolicy>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.policy.v1.DomainPolicy";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DomainPolicy;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DomainPolicy;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DomainPolicy;

  static equals(a: DomainPolicy | PlainMessage<DomainPolicy> | undefined, b: DomainPolicy | PlainMessage<DomainPolicy> | undefined): boolean;
}

/**
 * @generated from message zitadel.policy.v1.LabelPolicy
 */
export declare class LabelPolicy extends Message<LabelPolicy> {
  /**
   * @generated from field: zitadel.v1.ObjectDetails details = 1;
   */
  details?: ObjectDetails;

  /**
   * hex value for primary color
   *
   * @generated from field: string primary_color = 2;
   */
  primaryColor: string;

  /**
   * defines if the organisation's admin changed the policy
   *
   * @generated from field: bool is_default = 4;
   */
  isDefault: boolean;

  /**
   * hides the org suffix on the login form if the scope \"urn:zitadel:iam:org:domain:primary:{domainname}\" is set
   *
   * @generated from field: bool hide_login_name_suffix = 5;
   */
  hideLoginNameSuffix: boolean;

  /**
   * hex value for secondary color
   *
   * @generated from field: string warn_color = 6;
   */
  warnColor: string;

  /**
   * hex value for background color
   *
   * @generated from field: string background_color = 7;
   */
  backgroundColor: string;

  /**
   * hex value for font color
   *
   * @generated from field: string font_color = 8;
   */
  fontColor: string;

  /**
   * hex value for primary color dark theme
   *
   * @generated from field: string primary_color_dark = 9;
   */
  primaryColorDark: string;

  /**
   * hex value for background color dark theme
   *
   * @generated from field: string background_color_dark = 10;
   */
  backgroundColorDark: string;

  /**
   * hex value for warn color dark theme
   *
   * @generated from field: string warn_color_dark = 11;
   */
  warnColorDark: string;

  /**
   * hex value for font color dark theme
   *
   * @generated from field: string font_color_dark = 12;
   */
  fontColorDark: string;

  /**
   * @generated from field: bool disable_watermark = 13;
   */
  disableWatermark: boolean;

  /**
   * @generated from field: string logo_url = 14;
   */
  logoUrl: string;

  /**
   * @generated from field: string icon_url = 15;
   */
  iconUrl: string;

  /**
   * @generated from field: string logo_url_dark = 16;
   */
  logoUrlDark: string;

  /**
   * @generated from field: string icon_url_dark = 17;
   */
  iconUrlDark: string;

  /**
   * @generated from field: string font_url = 18;
   */
  fontUrl: string;

  constructor(data?: PartialMessage<LabelPolicy>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.policy.v1.LabelPolicy";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): LabelPolicy;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): LabelPolicy;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): LabelPolicy;

  static equals(a: LabelPolicy | PlainMessage<LabelPolicy> | undefined, b: LabelPolicy | PlainMessage<LabelPolicy> | undefined): boolean;
}

/**
 * @generated from message zitadel.policy.v1.LoginPolicy
 */
export declare class LoginPolicy extends Message<LoginPolicy> {
  /**
   * @generated from field: zitadel.v1.ObjectDetails details = 1;
   */
  details?: ObjectDetails;

  /**
   * @generated from field: bool allow_username_password = 2;
   */
  allowUsernamePassword: boolean;

  /**
   * @generated from field: bool allow_register = 3;
   */
  allowRegister: boolean;

  /**
   * @generated from field: bool allow_external_idp = 4;
   */
  allowExternalIdp: boolean;

  /**
   * @generated from field: bool force_mfa = 5;
   */
  forceMfa: boolean;

  /**
   * @generated from field: zitadel.policy.v1.PasswordlessType passwordless_type = 6;
   */
  passwordlessType: PasswordlessType;

  /**
   * @generated from field: bool is_default = 7;
   */
  isDefault: boolean;

  /**
   * @generated from field: bool hide_password_reset = 8;
   */
  hidePasswordReset: boolean;

  /**
   * @generated from field: bool ignore_unknown_usernames = 9;
   */
  ignoreUnknownUsernames: boolean;

  /**
   * @generated from field: string default_redirect_uri = 10;
   */
  defaultRedirectUri: string;

  /**
   * @generated from field: google.protobuf.Duration password_check_lifetime = 11;
   */
  passwordCheckLifetime?: Duration;

  /**
   * @generated from field: google.protobuf.Duration external_login_check_lifetime = 12;
   */
  externalLoginCheckLifetime?: Duration;

  /**
   * @generated from field: google.protobuf.Duration mfa_init_skip_lifetime = 13;
   */
  mfaInitSkipLifetime?: Duration;

  /**
   * @generated from field: google.protobuf.Duration second_factor_check_lifetime = 14;
   */
  secondFactorCheckLifetime?: Duration;

  /**
   * @generated from field: google.protobuf.Duration multi_factor_check_lifetime = 15;
   */
  multiFactorCheckLifetime?: Duration;

  /**
   * @generated from field: repeated zitadel.policy.v1.SecondFactorType second_factors = 16;
   */
  secondFactors: SecondFactorType[];

  /**
   * @generated from field: repeated zitadel.policy.v1.MultiFactorType multi_factors = 17;
   */
  multiFactors: MultiFactorType[];

  /**
   * @generated from field: repeated zitadel.idp.v1.IDPLoginPolicyLink idps = 18;
   */
  idps: IDPLoginPolicyLink[];

  /**
   * If set to true, the suffix (@domain.com) of an unknown username input on the login screen will be matched against the org domains and will redirect to the registration of that organisation on success.
   *
   * @generated from field: bool allow_domain_discovery = 19;
   */
  allowDomainDiscovery: boolean;

  /**
   * @generated from field: bool disable_login_with_email = 20;
   */
  disableLoginWithEmail: boolean;

  /**
   * @generated from field: bool disable_login_with_phone = 21;
   */
  disableLoginWithPhone: boolean;

  constructor(data?: PartialMessage<LoginPolicy>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.policy.v1.LoginPolicy";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): LoginPolicy;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): LoginPolicy;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): LoginPolicy;

  static equals(a: LoginPolicy | PlainMessage<LoginPolicy> | undefined, b: LoginPolicy | PlainMessage<LoginPolicy> | undefined): boolean;
}

/**
 * @generated from message zitadel.policy.v1.PasswordComplexityPolicy
 */
export declare class PasswordComplexityPolicy extends Message<PasswordComplexityPolicy> {
  /**
   * @generated from field: zitadel.v1.ObjectDetails details = 1;
   */
  details?: ObjectDetails;

  /**
   * @generated from field: uint64 min_length = 2;
   */
  minLength: bigint;

  /**
   * @generated from field: bool has_uppercase = 3;
   */
  hasUppercase: boolean;

  /**
   * @generated from field: bool has_lowercase = 4;
   */
  hasLowercase: boolean;

  /**
   * @generated from field: bool has_number = 5;
   */
  hasNumber: boolean;

  /**
   * @generated from field: bool has_symbol = 6;
   */
  hasSymbol: boolean;

  /**
   * @generated from field: bool is_default = 7;
   */
  isDefault: boolean;

  constructor(data?: PartialMessage<PasswordComplexityPolicy>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.policy.v1.PasswordComplexityPolicy";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PasswordComplexityPolicy;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PasswordComplexityPolicy;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PasswordComplexityPolicy;

  static equals(a: PasswordComplexityPolicy | PlainMessage<PasswordComplexityPolicy> | undefined, b: PasswordComplexityPolicy | PlainMessage<PasswordComplexityPolicy> | undefined): boolean;
}

/**
 * @generated from message zitadel.policy.v1.PasswordAgePolicy
 */
export declare class PasswordAgePolicy extends Message<PasswordAgePolicy> {
  /**
   * @generated from field: zitadel.v1.ObjectDetails details = 1;
   */
  details?: ObjectDetails;

  /**
   * @generated from field: uint64 max_age_days = 2;
   */
  maxAgeDays: bigint;

  /**
   * @generated from field: uint64 expire_warn_days = 3;
   */
  expireWarnDays: bigint;

  /**
   * @generated from field: bool is_default = 4;
   */
  isDefault: boolean;

  constructor(data?: PartialMessage<PasswordAgePolicy>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.policy.v1.PasswordAgePolicy";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PasswordAgePolicy;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PasswordAgePolicy;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PasswordAgePolicy;

  static equals(a: PasswordAgePolicy | PlainMessage<PasswordAgePolicy> | undefined, b: PasswordAgePolicy | PlainMessage<PasswordAgePolicy> | undefined): boolean;
}

/**
 * @generated from message zitadel.policy.v1.LockoutPolicy
 */
export declare class LockoutPolicy extends Message<LockoutPolicy> {
  /**
   * @generated from field: zitadel.v1.ObjectDetails details = 1;
   */
  details?: ObjectDetails;

  /**
   * @generated from field: uint64 max_password_attempts = 2;
   */
  maxPasswordAttempts: bigint;

  /**
   * @generated from field: bool is_default = 4;
   */
  isDefault: boolean;

  constructor(data?: PartialMessage<LockoutPolicy>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.policy.v1.LockoutPolicy";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): LockoutPolicy;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): LockoutPolicy;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): LockoutPolicy;

  static equals(a: LockoutPolicy | PlainMessage<LockoutPolicy> | undefined, b: LockoutPolicy | PlainMessage<LockoutPolicy> | undefined): boolean;
}

/**
 * @generated from message zitadel.policy.v1.PrivacyPolicy
 */
export declare class PrivacyPolicy extends Message<PrivacyPolicy> {
  /**
   * @generated from field: zitadel.v1.ObjectDetails details = 1;
   */
  details?: ObjectDetails;

  /**
   * @generated from field: string tos_link = 2;
   */
  tosLink: string;

  /**
   * @generated from field: string privacy_link = 3;
   */
  privacyLink: string;

  /**
   * @generated from field: bool is_default = 4;
   */
  isDefault: boolean;

  /**
   * @generated from field: string help_link = 5;
   */
  helpLink: string;

  constructor(data?: PartialMessage<PrivacyPolicy>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.policy.v1.PrivacyPolicy";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PrivacyPolicy;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PrivacyPolicy;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PrivacyPolicy;

  static equals(a: PrivacyPolicy | PlainMessage<PrivacyPolicy> | undefined, b: PrivacyPolicy | PlainMessage<PrivacyPolicy> | undefined): boolean;
}

/**
 * @generated from message zitadel.policy.v1.NotificationPolicy
 */
export declare class NotificationPolicy extends Message<NotificationPolicy> {
  /**
   * @generated from field: zitadel.v1.ObjectDetails details = 1;
   */
  details?: ObjectDetails;

  /**
   * @generated from field: bool is_default = 2;
   */
  isDefault: boolean;

  /**
   * @generated from field: bool password_change = 3;
   */
  passwordChange: boolean;

  constructor(data?: PartialMessage<NotificationPolicy>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.policy.v1.NotificationPolicy";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): NotificationPolicy;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): NotificationPolicy;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): NotificationPolicy;

  static equals(a: NotificationPolicy | PlainMessage<NotificationPolicy> | undefined, b: NotificationPolicy | PlainMessage<NotificationPolicy> | undefined): boolean;
}

