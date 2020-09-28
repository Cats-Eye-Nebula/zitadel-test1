// Code generated by protoc-gen-authmethod. DO NOT EDIT.

package management

import (
	"github.com/caos/zitadel/internal/api/authz"
)

/**
 * ManagementService
 */

const ManagementService_MethodPrefix = "caos.zitadel.management.api.v1.ManagementService"

var ManagementService_AuthMethods = authz.MethodMapping{

	"/caos.zitadel.management.api.v1.ManagementService/GetIam": authz.Option{
		Permission: "authenticated",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/IsUserUnique": authz.Option{
		Permission: "user.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetUserByID": authz.Option{
		Permission: "user.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetUserByLoginNameGlobal": authz.Option{
		Permission: "user.global.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchUsers": authz.Option{
		Permission: "user.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/CreateUser": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/DeactivateUser": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ReactivateUser": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/LockUser": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UnlockUser": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/DeleteUser": authz.Option{
		Permission: "user.delete",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UserChanges": authz.Option{
		Permission: "user.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/AddMachineKey": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/DeleteMachineKey": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchMachineKeys": authz.Option{
		Permission: "user.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetMachineKey": authz.Option{
		Permission: "user.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetUserProfile": authz.Option{
		Permission: "user.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdateUserProfile": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetUserEmail": authz.Option{
		Permission: "user.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ChangeUserUserName": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ChangeUserEmail": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ResendEmailVerificationMail": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetUserPhone": authz.Option{
		Permission: "user.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ChangeUserPhone": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveUserPhone": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ResendPhoneVerificationCode": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetUserAddress": authz.Option{
		Permission: "user.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdateUserAddress": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdateUserMachine": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchUserExternalIDPs": authz.Option{
		Permission: "user.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveExternalIDP": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetUserMfas": authz.Option{
		Permission: "user.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SendSetPasswordNotification": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SetInitialPassword": authz.Option{
		Permission: "user.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchUserMemberships": authz.Option{
		Permission: "user.membership.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetPasswordAgePolicy": authz.Option{
		Permission: "policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/CreatePasswordAgePolicy": authz.Option{
		Permission: "policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdatePasswordAgePolicy": authz.Option{
		Permission: "policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/DeletePasswordAgePolicy": authz.Option{
		Permission: "policy.delete",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetPasswordLockoutPolicy": authz.Option{
		Permission: "policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/CreatePasswordLockoutPolicy": authz.Option{
		Permission: "policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdatePasswordLockoutPolicy": authz.Option{
		Permission: "policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/DeletePasswordLockoutPolicy": authz.Option{
		Permission: "policy.delete",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/CreateOrg": authz.Option{
		Permission: "org.create",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/OrgChanges": authz.Option{
		Permission: "org.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetMyOrg": authz.Option{
		Permission: "org.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetOrgByDomainGlobal": authz.Option{
		Permission: "org.global.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/DeactivateMyOrg": authz.Option{
		Permission: "org.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ReactivateMyOrg": authz.Option{
		Permission: "org.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchMyOrgDomains": authz.Option{
		Permission: "org.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/AddMyOrgDomain": authz.Option{
		Permission: "org.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GenerateMyOrgDomainValidation": authz.Option{
		Permission: "org.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ValidateMyOrgDomain": authz.Option{
		Permission: "org.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SetMyPrimaryOrgDomain": authz.Option{
		Permission: "org.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveMyOrgDomain": authz.Option{
		Permission: "org.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetMyOrgIamPolicy": authz.Option{
		Permission: "authenticated",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetOrgMemberRoles": authz.Option{
		Permission: "org.member.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/AddMyOrgMember": authz.Option{
		Permission: "org.member.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ChangeMyOrgMember": authz.Option{
		Permission: "org.member.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveMyOrgMember": authz.Option{
		Permission: "org.member.delete",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchMyOrgMembers": authz.Option{
		Permission: "org.member.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ProjectChanges": authz.Option{
		Permission: "project.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchProjects": authz.Option{
		Permission: "project.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ProjectByID": authz.Option{
		Permission: "project.read",
		CheckParam: "Id",
	},

	"/caos.zitadel.management.api.v1.ManagementService/CreateProject": authz.Option{
		Permission: "project.create",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdateProject": authz.Option{
		Permission: "project.write",
		CheckParam: "Id",
	},

	"/caos.zitadel.management.api.v1.ManagementService/DeactivateProject": authz.Option{
		Permission: "project.write",
		CheckParam: "Id",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ReactivateProject": authz.Option{
		Permission: "project.write",
		CheckParam: "Id",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveProject": authz.Option{
		Permission: "project.delete",
		CheckParam: "Id",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchGrantedProjects": authz.Option{
		Permission: "project.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetGrantedProjectByID": authz.Option{
		Permission: "project.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetProjectMemberRoles": authz.Option{
		Permission: "project.member.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchProjectMembers": authz.Option{
		Permission: "project.member.read",
		CheckParam: "ProjectId",
	},

	"/caos.zitadel.management.api.v1.ManagementService/AddProjectMember": authz.Option{
		Permission: "project.member.write",
		CheckParam: "Id",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ChangeProjectMember": authz.Option{
		Permission: "project.member.write",
		CheckParam: "Id",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveProjectMember": authz.Option{
		Permission: "project.member.delete",
		CheckParam: "Id",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchProjectRoles": authz.Option{
		Permission: "project.role.read",
		CheckParam: "ProjectId",
	},

	"/caos.zitadel.management.api.v1.ManagementService/AddProjectRole": authz.Option{
		Permission: "project.role.write",
		CheckParam: "Id",
	},

	"/caos.zitadel.management.api.v1.ManagementService/BulkAddProjectRole": authz.Option{
		Permission: "project.role.write",
		CheckParam: "Id",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ChangeProjectRole": authz.Option{
		Permission: "project.role.write",
		CheckParam: "Id",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveProjectRole": authz.Option{
		Permission: "project.role.delete",
		CheckParam: "Id",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchApplications": authz.Option{
		Permission: "project.app.read",
		CheckParam: "ProjectId",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ApplicationByID": authz.Option{
		Permission: "project.app.read",
		CheckParam: "ProjectId",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ApplicationChanges": authz.Option{
		Permission: "project.app.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/CreateOIDCApplication": authz.Option{
		Permission: "project.app.write",
		CheckParam: "ProjectId",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdateApplication": authz.Option{
		Permission: "project.app.write",
		CheckParam: "ProjectId",
	},

	"/caos.zitadel.management.api.v1.ManagementService/DeactivateApplication": authz.Option{
		Permission: "project.app.write",
		CheckParam: "ProjectId",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ReactivateApplication": authz.Option{
		Permission: "project.app.write",
		CheckParam: "ProjectId",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveApplication": authz.Option{
		Permission: "project.app.delete",
		CheckParam: "ProjectId",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdateApplicationOIDCConfig": authz.Option{
		Permission: "project.app.write",
		CheckParam: "ProjectId",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RegenerateOIDCClientSecret": authz.Option{
		Permission: "project.app.write",
		CheckParam: "ProjectId",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchProjectGrants": authz.Option{
		Permission: "project.grant.read",
		CheckParam: "ProjectId",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ProjectGrantByID": authz.Option{
		Permission: "project.grant.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/CreateProjectGrant": authz.Option{
		Permission: "project.grant.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdateProjectGrant": authz.Option{
		Permission: "project.grant.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/DeactivateProjectGrant": authz.Option{
		Permission: "project.grant.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ReactivateProjectGrant": authz.Option{
		Permission: "project.grant.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveProjectGrant": authz.Option{
		Permission: "project.grant.delete",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetProjectGrantMemberRoles": authz.Option{
		Permission: "project.grant.member.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchProjectGrantMembers": authz.Option{
		Permission: "project.grant.member.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/AddProjectGrantMember": authz.Option{
		Permission: "project.grant.member.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ChangeProjectGrantMember": authz.Option{
		Permission: "project.grant.member.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveProjectGrantMember": authz.Option{
		Permission: "project.grant.member.delete",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchUserGrants": authz.Option{
		Permission: "user.grant.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UserGrantByID": authz.Option{
		Permission: "user.grant.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/CreateUserGrant": authz.Option{
		Permission: "user.grant.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdateUserGrant": authz.Option{
		Permission: "user.grant.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/DeactivateUserGrant": authz.Option{
		Permission: "user.grant.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ReactivateUserGrant": authz.Option{
		Permission: "user.grant.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveUserGrant": authz.Option{
		Permission: "user.grant.delete",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/BulkRemoveUserGrant": authz.Option{
		Permission: "user.grant.delete",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/IdpByID": authz.Option{
		Permission: "org.idp.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/CreateOidcIdp": authz.Option{
		Permission: "org.idp.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdateIdpConfig": authz.Option{
		Permission: "org.idp.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/DeactivateIdpConfig": authz.Option{
		Permission: "org.idp.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/ReactivateIdpConfig": authz.Option{
		Permission: "org.idp.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveIdpConfig": authz.Option{
		Permission: "org.idp.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdateOidcIdpConfig": authz.Option{
		Permission: "org.idp.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/SearchIdps": authz.Option{
		Permission: "org.idp.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetLoginPolicy": authz.Option{
		Permission: "policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetDefaultLoginPolicy": authz.Option{
		Permission: "policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/CreateLoginPolicy": authz.Option{
		Permission: "policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdateLoginPolicy": authz.Option{
		Permission: "policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveLoginPolicy": authz.Option{
		Permission: "policy.delete",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetLoginPolicyIdpProviders": authz.Option{
		Permission: "policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/AddIdpProviderToLoginPolicy": authz.Option{
		Permission: "policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemoveIdpProviderFromLoginPolicy": authz.Option{
		Permission: "policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetPasswordComplexityPolicy": authz.Option{
		Permission: "policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/GetDefaultPasswordComplexityPolicy": authz.Option{
		Permission: "policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/CreatePasswordComplexityPolicy": authz.Option{
		Permission: "policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/UpdatePasswordComplexityPolicy": authz.Option{
		Permission: "policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.management.api.v1.ManagementService/RemovePasswordComplexityPolicy": authz.Option{
		Permission: "policy.delete",
		CheckParam: "",
	},
}
