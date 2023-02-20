// @generated by protoc-gen-es v1.0.0
// @generated from file zitadel/app.proto (package zitadel.app.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { Duration, proto3 } from "@bufbuild/protobuf";
import { ObjectDetails, TextQueryMethod } from "./object_pb.js";
import { LocalizedMessage } from "./message_pb.js";

/**
 * @generated from enum zitadel.app.v1.AppState
 */
export const AppState = proto3.makeEnum(
  "zitadel.app.v1.AppState",
  [
    {no: 0, name: "APP_STATE_UNSPECIFIED", localName: "UNSPECIFIED"},
    {no: 1, name: "APP_STATE_ACTIVE", localName: "ACTIVE"},
    {no: 2, name: "APP_STATE_INACTIVE", localName: "INACTIVE"},
  ],
);

/**
 * @generated from enum zitadel.app.v1.OIDCResponseType
 */
export const OIDCResponseType = proto3.makeEnum(
  "zitadel.app.v1.OIDCResponseType",
  [
    {no: 0, name: "OIDC_RESPONSE_TYPE_CODE"},
    {no: 1, name: "OIDC_RESPONSE_TYPE_ID_TOKEN"},
    {no: 2, name: "OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN"},
  ],
);

/**
 * @generated from enum zitadel.app.v1.OIDCGrantType
 */
export const OIDCGrantType = proto3.makeEnum(
  "zitadel.app.v1.OIDCGrantType",
  [
    {no: 0, name: "OIDC_GRANT_TYPE_AUTHORIZATION_CODE"},
    {no: 1, name: "OIDC_GRANT_TYPE_IMPLICIT"},
    {no: 2, name: "OIDC_GRANT_TYPE_REFRESH_TOKEN"},
  ],
);

/**
 * @generated from enum zitadel.app.v1.OIDCAppType
 */
export const OIDCAppType = proto3.makeEnum(
  "zitadel.app.v1.OIDCAppType",
  [
    {no: 0, name: "OIDC_APP_TYPE_WEB"},
    {no: 1, name: "OIDC_APP_TYPE_USER_AGENT"},
    {no: 2, name: "OIDC_APP_TYPE_NATIVE"},
  ],
);

/**
 * @generated from enum zitadel.app.v1.OIDCAuthMethodType
 */
export const OIDCAuthMethodType = proto3.makeEnum(
  "zitadel.app.v1.OIDCAuthMethodType",
  [
    {no: 0, name: "OIDC_AUTH_METHOD_TYPE_BASIC"},
    {no: 1, name: "OIDC_AUTH_METHOD_TYPE_POST"},
    {no: 2, name: "OIDC_AUTH_METHOD_TYPE_NONE"},
    {no: 3, name: "OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT"},
  ],
);

/**
 * @generated from enum zitadel.app.v1.OIDCVersion
 */
export const OIDCVersion = proto3.makeEnum(
  "zitadel.app.v1.OIDCVersion",
  [
    {no: 0, name: "OIDC_VERSION_1_0"},
  ],
);

/**
 * @generated from enum zitadel.app.v1.OIDCTokenType
 */
export const OIDCTokenType = proto3.makeEnum(
  "zitadel.app.v1.OIDCTokenType",
  [
    {no: 0, name: "OIDC_TOKEN_TYPE_BEARER"},
    {no: 1, name: "OIDC_TOKEN_TYPE_JWT"},
  ],
);

/**
 * @generated from enum zitadel.app.v1.APIAuthMethodType
 */
export const APIAuthMethodType = proto3.makeEnum(
  "zitadel.app.v1.APIAuthMethodType",
  [
    {no: 0, name: "API_AUTH_METHOD_TYPE_BASIC"},
    {no: 1, name: "API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT"},
  ],
);

/**
 * @generated from message zitadel.app.v1.App
 */
export const App = proto3.makeMessageType(
  "zitadel.app.v1.App",
  () => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "details", kind: "message", T: ObjectDetails },
    { no: 3, name: "state", kind: "enum", T: proto3.getEnumType(AppState) },
    { no: 4, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "oidc_config", kind: "message", T: OIDCConfig, oneof: "config" },
    { no: 6, name: "api_config", kind: "message", T: APIConfig, oneof: "config" },
    { no: 7, name: "saml_config", kind: "message", T: SAMLConfig, oneof: "config" },
  ],
);

/**
 * @generated from message zitadel.app.v1.AppQuery
 */
export const AppQuery = proto3.makeMessageType(
  "zitadel.app.v1.AppQuery",
  () => [
    { no: 1, name: "name_query", kind: "message", T: AppNameQuery, oneof: "query" },
  ],
);

/**
 * @generated from message zitadel.app.v1.AppNameQuery
 */
export const AppNameQuery = proto3.makeMessageType(
  "zitadel.app.v1.AppNameQuery",
  () => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "method", kind: "enum", T: proto3.getEnumType(TextQueryMethod) },
  ],
);

/**
 * @generated from message zitadel.app.v1.OIDCConfig
 */
export const OIDCConfig = proto3.makeMessageType(
  "zitadel.app.v1.OIDCConfig",
  () => [
    { no: 1, name: "redirect_uris", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 2, name: "response_types", kind: "enum", T: proto3.getEnumType(OIDCResponseType), repeated: true },
    { no: 3, name: "grant_types", kind: "enum", T: proto3.getEnumType(OIDCGrantType), repeated: true },
    { no: 4, name: "app_type", kind: "enum", T: proto3.getEnumType(OIDCAppType) },
    { no: 5, name: "client_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 7, name: "auth_method_type", kind: "enum", T: proto3.getEnumType(OIDCAuthMethodType) },
    { no: 8, name: "post_logout_redirect_uris", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 9, name: "version", kind: "enum", T: proto3.getEnumType(OIDCVersion) },
    { no: 10, name: "none_compliant", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 11, name: "compliance_problems", kind: "message", T: LocalizedMessage, repeated: true },
    { no: 12, name: "dev_mode", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 13, name: "access_token_type", kind: "enum", T: proto3.getEnumType(OIDCTokenType) },
    { no: 14, name: "access_token_role_assertion", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 15, name: "id_token_role_assertion", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 16, name: "id_token_userinfo_assertion", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 17, name: "clock_skew", kind: "message", T: Duration },
    { no: 18, name: "additional_origins", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 19, name: "allowed_origins", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
  ],
);

/**
 * @generated from message zitadel.app.v1.SAMLConfig
 */
export const SAMLConfig = proto3.makeMessageType(
  "zitadel.app.v1.SAMLConfig",
  () => [
    { no: 1, name: "metadata_xml", kind: "scalar", T: 12 /* ScalarType.BYTES */, oneof: "metadata" },
    { no: 2, name: "metadata_url", kind: "scalar", T: 9 /* ScalarType.STRING */, oneof: "metadata" },
  ],
);

/**
 * @generated from message zitadel.app.v1.APIConfig
 */
export const APIConfig = proto3.makeMessageType(
  "zitadel.app.v1.APIConfig",
  () => [
    { no: 1, name: "client_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "auth_method_type", kind: "enum", T: proto3.getEnumType(APIAuthMethodType) },
  ],
);

