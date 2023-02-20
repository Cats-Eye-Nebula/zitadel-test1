// @generated by protoc-gen-es v1.0.0
// @generated from file zitadel/project.proto (package zitadel.project.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import type { ObjectDetails, TextQueryMethod } from "./object_pb.js";

/**
 * @generated from enum zitadel.project.v1.ProjectState
 */
export declare enum ProjectState {
  /**
   * @generated from enum value: PROJECT_STATE_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: PROJECT_STATE_ACTIVE = 1;
   */
  ACTIVE = 1,

  /**
   * @generated from enum value: PROJECT_STATE_INACTIVE = 2;
   */
  INACTIVE = 2,
}

/**
 * @generated from enum zitadel.project.v1.PrivateLabelingSetting
 */
export declare enum PrivateLabelingSetting {
  /**
   * @generated from enum value: PRIVATE_LABELING_SETTING_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY = 1;
   */
  ENFORCE_PROJECT_RESOURCE_OWNER_POLICY = 1,

  /**
   * @generated from enum value: PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY = 2;
   */
  ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY = 2,
}

/**
 * @generated from enum zitadel.project.v1.ProjectGrantState
 */
export declare enum ProjectGrantState {
  /**
   * @generated from enum value: PROJECT_GRANT_STATE_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: PROJECT_GRANT_STATE_ACTIVE = 1;
   */
  ACTIVE = 1,

  /**
   * @generated from enum value: PROJECT_GRANT_STATE_INACTIVE = 2;
   */
  INACTIVE = 2,
}

/**
 * @generated from message zitadel.project.v1.Project
 */
export declare class Project extends Message<Project> {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: zitadel.v1.ObjectDetails details = 2;
   */
  details?: ObjectDetails;

  /**
   * @generated from field: string name = 3;
   */
  name: string;

  /**
   * @generated from field: zitadel.project.v1.ProjectState state = 4;
   */
  state: ProjectState;

  /**
   * describes if roles of user should be added in token
   *
   * @generated from field: bool project_role_assertion = 5;
   */
  projectRoleAssertion: boolean;

  /**
   * ZITADEL checks if the user has at least one on this project
   *
   * @generated from field: bool project_role_check = 6;
   */
  projectRoleCheck: boolean;

  /**
   * ZITADEL checks if the org of the user has permission to this project
   *
   * @generated from field: bool has_project_check = 7;
   */
  hasProjectCheck: boolean;

  /**
   * Defines from where the private labeling should be triggered
   *
   * @generated from field: zitadel.project.v1.PrivateLabelingSetting private_labeling_setting = 8;
   */
  privateLabelingSetting: PrivateLabelingSetting;

  constructor(data?: PartialMessage<Project>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.Project";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Project;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Project;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Project;

  static equals(a: Project | PlainMessage<Project> | undefined, b: Project | PlainMessage<Project> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.GrantedProject
 */
export declare class GrantedProject extends Message<GrantedProject> {
  /**
   * @generated from field: string grant_id = 1;
   */
  grantId: string;

  /**
   * @generated from field: string granted_org_id = 2;
   */
  grantedOrgId: string;

  /**
   * @generated from field: string granted_org_name = 3;
   */
  grantedOrgName: string;

  /**
   * @generated from field: repeated string granted_role_keys = 4;
   */
  grantedRoleKeys: string[];

  /**
   * @generated from field: zitadel.project.v1.ProjectGrantState state = 5;
   */
  state: ProjectGrantState;

  /**
   * @generated from field: string project_id = 6;
   */
  projectId: string;

  /**
   * @generated from field: string project_name = 7;
   */
  projectName: string;

  /**
   * @generated from field: string project_owner_id = 8;
   */
  projectOwnerId: string;

  /**
   * @generated from field: string project_owner_name = 9;
   */
  projectOwnerName: string;

  /**
   * @generated from field: zitadel.v1.ObjectDetails details = 10;
   */
  details?: ObjectDetails;

  constructor(data?: PartialMessage<GrantedProject>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.GrantedProject";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GrantedProject;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GrantedProject;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GrantedProject;

  static equals(a: GrantedProject | PlainMessage<GrantedProject> | undefined, b: GrantedProject | PlainMessage<GrantedProject> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.ProjectQuery
 */
export declare class ProjectQuery extends Message<ProjectQuery> {
  /**
   * @generated from oneof zitadel.project.v1.ProjectQuery.query
   */
  query: {
    /**
     * @generated from field: zitadel.project.v1.ProjectNameQuery name_query = 1;
     */
    value: ProjectNameQuery;
    case: "nameQuery";
  } | {
    /**
     * @generated from field: zitadel.project.v1.ProjectResourceOwnerQuery project_resource_owner_query = 2;
     */
    value: ProjectResourceOwnerQuery;
    case: "projectResourceOwnerQuery";
  } | { case: undefined; value?: undefined };

  constructor(data?: PartialMessage<ProjectQuery>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.ProjectQuery";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ProjectQuery;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ProjectQuery;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ProjectQuery;

  static equals(a: ProjectQuery | PlainMessage<ProjectQuery> | undefined, b: ProjectQuery | PlainMessage<ProjectQuery> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.ProjectNameQuery
 */
export declare class ProjectNameQuery extends Message<ProjectNameQuery> {
  /**
   * @generated from field: string name = 1;
   */
  name: string;

  /**
   * @generated from field: zitadel.v1.TextQueryMethod method = 2;
   */
  method: TextQueryMethod;

  constructor(data?: PartialMessage<ProjectNameQuery>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.ProjectNameQuery";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ProjectNameQuery;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ProjectNameQuery;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ProjectNameQuery;

  static equals(a: ProjectNameQuery | PlainMessage<ProjectNameQuery> | undefined, b: ProjectNameQuery | PlainMessage<ProjectNameQuery> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.ProjectResourceOwnerQuery
 */
export declare class ProjectResourceOwnerQuery extends Message<ProjectResourceOwnerQuery> {
  /**
   * @generated from field: string resource_owner = 1;
   */
  resourceOwner: string;

  constructor(data?: PartialMessage<ProjectResourceOwnerQuery>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.ProjectResourceOwnerQuery";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ProjectResourceOwnerQuery;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ProjectResourceOwnerQuery;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ProjectResourceOwnerQuery;

  static equals(a: ProjectResourceOwnerQuery | PlainMessage<ProjectResourceOwnerQuery> | undefined, b: ProjectResourceOwnerQuery | PlainMessage<ProjectResourceOwnerQuery> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.Role
 */
export declare class Role extends Message<Role> {
  /**
   * @generated from field: string key = 1;
   */
  key: string;

  /**
   * @generated from field: zitadel.v1.ObjectDetails details = 2;
   */
  details?: ObjectDetails;

  /**
   * @generated from field: string display_name = 3;
   */
  displayName: string;

  /**
   * @generated from field: string group = 4;
   */
  group: string;

  constructor(data?: PartialMessage<Role>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.Role";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Role;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Role;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Role;

  static equals(a: Role | PlainMessage<Role> | undefined, b: Role | PlainMessage<Role> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.RoleQuery
 */
export declare class RoleQuery extends Message<RoleQuery> {
  /**
   * @generated from oneof zitadel.project.v1.RoleQuery.query
   */
  query: {
    /**
     * @generated from field: zitadel.project.v1.RoleKeyQuery key_query = 1;
     */
    value: RoleKeyQuery;
    case: "keyQuery";
  } | {
    /**
     * @generated from field: zitadel.project.v1.RoleDisplayNameQuery display_name_query = 2;
     */
    value: RoleDisplayNameQuery;
    case: "displayNameQuery";
  } | { case: undefined; value?: undefined };

  constructor(data?: PartialMessage<RoleQuery>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.RoleQuery";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RoleQuery;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RoleQuery;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RoleQuery;

  static equals(a: RoleQuery | PlainMessage<RoleQuery> | undefined, b: RoleQuery | PlainMessage<RoleQuery> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.RoleKeyQuery
 */
export declare class RoleKeyQuery extends Message<RoleKeyQuery> {
  /**
   * @generated from field: string key = 1;
   */
  key: string;

  /**
   * @generated from field: zitadel.v1.TextQueryMethod method = 2;
   */
  method: TextQueryMethod;

  constructor(data?: PartialMessage<RoleKeyQuery>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.RoleKeyQuery";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RoleKeyQuery;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RoleKeyQuery;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RoleKeyQuery;

  static equals(a: RoleKeyQuery | PlainMessage<RoleKeyQuery> | undefined, b: RoleKeyQuery | PlainMessage<RoleKeyQuery> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.RoleDisplayNameQuery
 */
export declare class RoleDisplayNameQuery extends Message<RoleDisplayNameQuery> {
  /**
   * @generated from field: string display_name = 1;
   */
  displayName: string;

  /**
   * @generated from field: zitadel.v1.TextQueryMethod method = 2;
   */
  method: TextQueryMethod;

  constructor(data?: PartialMessage<RoleDisplayNameQuery>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.RoleDisplayNameQuery";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RoleDisplayNameQuery;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RoleDisplayNameQuery;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RoleDisplayNameQuery;

  static equals(a: RoleDisplayNameQuery | PlainMessage<RoleDisplayNameQuery> | undefined, b: RoleDisplayNameQuery | PlainMessage<RoleDisplayNameQuery> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.ProjectGrantQuery
 */
export declare class ProjectGrantQuery extends Message<ProjectGrantQuery> {
  /**
   * @generated from oneof zitadel.project.v1.ProjectGrantQuery.query
   */
  query: {
    /**
     * @generated from field: zitadel.project.v1.GrantProjectNameQuery project_name_query = 1;
     */
    value: GrantProjectNameQuery;
    case: "projectNameQuery";
  } | {
    /**
     * @generated from field: zitadel.project.v1.GrantRoleKeyQuery role_key_query = 2;
     */
    value: GrantRoleKeyQuery;
    case: "roleKeyQuery";
  } | { case: undefined; value?: undefined };

  constructor(data?: PartialMessage<ProjectGrantQuery>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.ProjectGrantQuery";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ProjectGrantQuery;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ProjectGrantQuery;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ProjectGrantQuery;

  static equals(a: ProjectGrantQuery | PlainMessage<ProjectGrantQuery> | undefined, b: ProjectGrantQuery | PlainMessage<ProjectGrantQuery> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.AllProjectGrantQuery
 */
export declare class AllProjectGrantQuery extends Message<AllProjectGrantQuery> {
  /**
   * @generated from oneof zitadel.project.v1.AllProjectGrantQuery.query
   */
  query: {
    /**
     * @generated from field: zitadel.project.v1.GrantProjectNameQuery project_name_query = 1;
     */
    value: GrantProjectNameQuery;
    case: "projectNameQuery";
  } | {
    /**
     * @generated from field: zitadel.project.v1.GrantRoleKeyQuery role_key_query = 2;
     */
    value: GrantRoleKeyQuery;
    case: "roleKeyQuery";
  } | {
    /**
     * @generated from field: zitadel.project.v1.ProjectIDQuery project_id_query = 3;
     */
    value: ProjectIDQuery;
    case: "projectIdQuery";
  } | {
    /**
     * @generated from field: zitadel.project.v1.GrantedOrgIDQuery granted_org_id_query = 4;
     */
    value: GrantedOrgIDQuery;
    case: "grantedOrgIdQuery";
  } | { case: undefined; value?: undefined };

  constructor(data?: PartialMessage<AllProjectGrantQuery>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.AllProjectGrantQuery";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AllProjectGrantQuery;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AllProjectGrantQuery;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AllProjectGrantQuery;

  static equals(a: AllProjectGrantQuery | PlainMessage<AllProjectGrantQuery> | undefined, b: AllProjectGrantQuery | PlainMessage<AllProjectGrantQuery> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.GrantProjectNameQuery
 */
export declare class GrantProjectNameQuery extends Message<GrantProjectNameQuery> {
  /**
   * @generated from field: string name = 1;
   */
  name: string;

  /**
   * @generated from field: zitadel.v1.TextQueryMethod method = 2;
   */
  method: TextQueryMethod;

  constructor(data?: PartialMessage<GrantProjectNameQuery>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.GrantProjectNameQuery";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GrantProjectNameQuery;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GrantProjectNameQuery;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GrantProjectNameQuery;

  static equals(a: GrantProjectNameQuery | PlainMessage<GrantProjectNameQuery> | undefined, b: GrantProjectNameQuery | PlainMessage<GrantProjectNameQuery> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.GrantRoleKeyQuery
 */
export declare class GrantRoleKeyQuery extends Message<GrantRoleKeyQuery> {
  /**
   * @generated from field: string role_key = 1;
   */
  roleKey: string;

  /**
   * @generated from field: zitadel.v1.TextQueryMethod method = 2;
   */
  method: TextQueryMethod;

  constructor(data?: PartialMessage<GrantRoleKeyQuery>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.GrantRoleKeyQuery";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GrantRoleKeyQuery;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GrantRoleKeyQuery;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GrantRoleKeyQuery;

  static equals(a: GrantRoleKeyQuery | PlainMessage<GrantRoleKeyQuery> | undefined, b: GrantRoleKeyQuery | PlainMessage<GrantRoleKeyQuery> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.ProjectIDQuery
 */
export declare class ProjectIDQuery extends Message<ProjectIDQuery> {
  /**
   * @generated from field: string project_id = 1;
   */
  projectId: string;

  constructor(data?: PartialMessage<ProjectIDQuery>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.ProjectIDQuery";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ProjectIDQuery;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ProjectIDQuery;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ProjectIDQuery;

  static equals(a: ProjectIDQuery | PlainMessage<ProjectIDQuery> | undefined, b: ProjectIDQuery | PlainMessage<ProjectIDQuery> | undefined): boolean;
}

/**
 * @generated from message zitadel.project.v1.GrantedOrgIDQuery
 */
export declare class GrantedOrgIDQuery extends Message<GrantedOrgIDQuery> {
  /**
   * @generated from field: string granted_org_id = 1;
   */
  grantedOrgId: string;

  constructor(data?: PartialMessage<GrantedOrgIDQuery>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zitadel.project.v1.GrantedOrgIDQuery";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GrantedOrgIDQuery;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GrantedOrgIDQuery;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GrantedOrgIDQuery;

  static equals(a: GrantedOrgIDQuery | PlainMessage<GrantedOrgIDQuery> | undefined, b: GrantedOrgIDQuery | PlainMessage<GrantedOrgIDQuery> | undefined): boolean;
}

