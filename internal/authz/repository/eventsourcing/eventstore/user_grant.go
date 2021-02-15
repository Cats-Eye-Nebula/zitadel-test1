package eventstore

import (
	"context"
	global_model "github.com/caos/zitadel/internal/model"
	user_model "github.com/caos/zitadel/internal/user/model"
	user_view_model "github.com/caos/zitadel/internal/user/repository/view/model"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/usergrant/repository/view/model"
)

type UserGrantRepo struct {
	View         *view.View
	IamID        string
	IamProjectID string
	Auth         authz.Config
	IamEvents    *iam_event.IAMEventstore
}

func (repo *UserGrantRepo) Health() error {
	return repo.View.Health()
}

func (repo *UserGrantRepo) SearchMyMemberships(ctx context.Context) ([]*authz.Membership, error) {
	ctxData := authz.GetCtxData(ctx)
	orgMemberships, orgCount, err := repo.View.SearchUserMemberships(&user_model.UserMembershipSearchRequest{
		Queries: []*user_model.UserMembershipSearchQuery{
			{
				Key:    user_model.UserMembershipSearchKeyUserID,
				Method: global_model.SearchMethodEquals,
				Value:  ctxData.UserID,
			},
			{
				Key:    user_model.UserMembershipSearchKeyResourceOwner,
				Method: global_model.SearchMethodEquals,
				Value:  ctxData.OrgID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	iamMemberships, iamCount, err := repo.View.SearchUserMemberships(&user_model.UserMembershipSearchRequest{
		Queries: []*user_model.UserMembershipSearchQuery{
			{
				Key:    user_model.UserMembershipSearchKeyUserID,
				Method: global_model.SearchMethodEquals,
				Value:  ctxData.UserID,
			},
			{
				Key:    user_model.UserMembershipSearchKeyAggregateID,
				Method: global_model.SearchMethodEquals,
				Value:  repo.IamID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if orgCount == 0 && iamCount == 0 {
		return []*authz.Membership{}, nil
	}
	orgMemberships = append(orgMemberships, iamMemberships...)
	return userMembershipsToMemberships(orgMemberships), nil
}

func (repo *UserGrantRepo) SearchMyZitadelPermissions(ctx context.Context) ([]string, error) {
	ctxData := authz.GetCtxData(ctx)
	orgMemberships, orgCount, err := repo.View.SearchUserMemberships(&user_model.UserMembershipSearchRequest{
		Queries: []*user_model.UserMembershipSearchQuery{
			{
				Key:    user_model.UserMembershipSearchKeyUserID,
				Method: global_model.SearchMethodEquals,
				Value:  ctxData.UserID,
			},
			{
				Key:    user_model.UserMembershipSearchKeyResourceOwner,
				Method: global_model.SearchMethodEquals,
				Value:  ctxData.OrgID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	iamMemberships, iamCount, err := repo.View.SearchUserMemberships(&user_model.UserMembershipSearchRequest{
		Queries: []*user_model.UserMembershipSearchQuery{
			{
				Key:    user_model.UserMembershipSearchKeyUserID,
				Method: global_model.SearchMethodEquals,
				Value:  ctxData.UserID,
			},
			{
				Key:    user_model.UserMembershipSearchKeyAggregateID,
				Method: global_model.SearchMethodEquals,
				Value:  repo.IamID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if orgCount == 0 && iamCount == 0 {
		return []string{}, nil
	}
	orgMemberships = append(orgMemberships, iamMemberships...)
	permissions := &grant_model.Permissions{Permissions: []string{}}
	for _, membership := range orgMemberships {
		for _, role := range membership.Roles {
			permissions = repo.mapRoleToPermission(permissions, membership, role)
		}
	}
	return permissions.Permissions, nil
}

func (repo *UserGrantRepo) FillIamProjectID(ctx context.Context) error {
	if repo.IamProjectID != "" {
		return nil
	}
	iam, err := repo.IamEvents.IAMByID(ctx, repo.IamID)
	if err != nil {
		return err
	}
	if iam.SetUpDone < iam_model.StepCount-1 {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-skiwS", "Setup not done")
	}
	repo.IamProjectID = iam.IAMProjectID
	return nil
}

func mergeOrgAndAdminGrant(ctxData authz.CtxData, orgGrant, iamAdminGrant *model.UserGrantView) (grant *authz.Grant) {
	if orgGrant != nil {
		roles := orgGrant.RoleKeys
		if iamAdminGrant != nil {
			roles = addIamAdminRoles(roles, iamAdminGrant.RoleKeys)
		}
		grant = &authz.Grant{OrgID: orgGrant.ResourceOwner, Roles: roles}
	} else if iamAdminGrant != nil {
		grant = &authz.Grant{
			OrgID: ctxData.OrgID,
			Roles: iamAdminGrant.RoleKeys,
		}
	}
	return grant
}

func addIamAdminRoles(orgRoles, iamAdminRoles []string) []string {
	result := make([]string, 0)
	result = append(result, iamAdminRoles...)
	for _, role := range orgRoles {
		if !authz.ExistsPerm(result, role) {
			result = append(result, role)
		}
	}
	return result
}

func (repo *UserGrantRepo) mapRoleToPermission(permissions *grant_model.Permissions, membership *user_view_model.UserMembershipView, role string) *grant_model.Permissions {
	for _, mapping := range repo.Auth.RolePermissionMappings {
		if mapping.Role == role {
			ctxID := ""
			if membership.MemberType == int32(user_model.MemberTypeProject) || membership.MemberType == int32(user_model.MemberTypeProjectGrant) {
				ctxID = membership.ObjectID
			}
			permissions.AppendPermissions(ctxID, mapping.Permissions...)
		}
	}
	return permissions
}

func userMembershipToMembership(membership *user_view_model.UserMembershipView) *authz.Membership {
	return &authz.Membership{
		MemberType:  authz.MemberType(membership.MemberType),
		AggregateID: membership.AggregateID,
		ObjectID:    membership.ObjectID,
		Roles:       membership.Roles,
	}
}

func userMembershipsToMemberships(memberships []*user_view_model.UserMembershipView) []*authz.Membership {
	result := make([]*authz.Membership, len(memberships))
	for i, m := range memberships {
		result[i] = userMembershipToMembership(m)
	}
	return result
}
