package eventstore

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	es_proj_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	usr_grant_model "github.com/caos/zitadel/internal/usergrant/model"
	usr_grant_event "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing"
	"strings"

	"github.com/caos/zitadel/internal/api/auth"
	global_model "github.com/caos/zitadel/internal/model"

	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"

	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
)

type ProjectRepo struct {
	es_int.Eventstore
	SearchLimit     uint64
	ProjectEvents   *proj_event.ProjectEventstore
	UserGrantEvents *usr_grant_event.UserGrantEventStore
	View            *view.View
	Roles           []string
}

func (repo *ProjectRepo) ProjectByID(ctx context.Context, id string) (*proj_model.ProjectView, error) {
	project, err := repo.View.ProjectByID(id)
	if err != nil {
		return nil, err
	}
	return model.ProjectToModel(project), nil
}

func (repo *ProjectRepo) CreateProject(ctx context.Context, name string) (*proj_model.Project, error) {
	project := &proj_model.Project{Name: name}
	return repo.ProjectEvents.CreateProject(ctx, project)
}

func (repo *ProjectRepo) UpdateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	return repo.ProjectEvents.UpdateProject(ctx, project)
}

func (repo *ProjectRepo) DeactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	return repo.ProjectEvents.DeactivateProject(ctx, id)
}

func (repo *ProjectRepo) ReactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	return repo.ProjectEvents.ReactivateProject(ctx, id)
}

func (repo *ProjectRepo) SearchProjects(ctx context.Context, request *proj_model.ProjectViewSearchRequest) (*proj_model.ProjectViewSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)

	permissions := auth.GetPermissionsFromCtx(ctx)
	if !auth.HasGlobalPermission(permissions) {
		ids := auth.GetPermissionCtxIDs(permissions)
		request.Queries = append(request.Queries, &proj_model.ProjectViewSearchQuery{Key: proj_model.PROJECTSEARCHKEY_PROJECTID, Method: global_model.SEARCHMETHOD_IS_ONE_OF, Value: ids})
	}

	projects, count, err := repo.View.SearchProjects(request)
	if err != nil {
		return nil, err
	}
	return &proj_model.ProjectViewSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.ProjectsToModel(projects),
	}, nil
}

func (repo *ProjectRepo) ProjectGrantViewByID(ctx context.Context, grantID string) (project *proj_model.ProjectGrantView, err error) {
	p, err := repo.View.ProjectGrantByID(grantID)
	if err != nil {
		return nil, err
	}
	return model.ProjectGrantToModel(p), nil
}

func (repo *ProjectRepo) ProjectMemberByID(ctx context.Context, projectID, userID string) (member *proj_model.ProjectMember, err error) {
	member = proj_model.NewProjectMember(projectID, userID)
	return repo.ProjectEvents.ProjectMemberByIDs(ctx, member)
}

func (repo *ProjectRepo) AddProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	return repo.ProjectEvents.AddProjectMember(ctx, member)
}

func (repo *ProjectRepo) ChangeProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	return repo.ProjectEvents.ChangeProjectMember(ctx, member)
}

func (repo *ProjectRepo) RemoveProjectMember(ctx context.Context, projectID, userID string) error {
	member := proj_model.NewProjectMember(projectID, userID)
	return repo.ProjectEvents.RemoveProjectMember(ctx, member)
}

func (repo *ProjectRepo) SearchProjectMembers(ctx context.Context, request *proj_model.ProjectMemberSearchRequest) (*proj_model.ProjectMemberSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	members, count, err := repo.View.SearchProjectMembers(request)
	if err != nil {
		return nil, err
	}
	return &proj_model.ProjectMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.ProjectMembersToModel(members),
	}, nil
}

func (repo *ProjectRepo) AddProjectRole(ctx context.Context, role *proj_model.ProjectRole) (*proj_model.ProjectRole, error) {
	return repo.ProjectEvents.AddProjectRoles(ctx, role)
}

func (repo *ProjectRepo) BulkAddProjectRole(ctx context.Context, roles []*proj_model.ProjectRole) error {
	_, err := repo.ProjectEvents.AddProjectRoles(ctx, roles...)
	return err
}

func (repo *ProjectRepo) ChangeProjectRole(ctx context.Context, member *proj_model.ProjectRole) (*proj_model.ProjectRole, error) {
	return repo.ProjectEvents.ChangeProjectRole(ctx, member)
}

func (repo *ProjectRepo) RemoveProjectRole(ctx context.Context, projectID, key string) error {
	role := proj_model.NewProjectRole(projectID, key)
	aggregates := make([]*es_models.Aggregate, 0)
	project, agg, err := repo.ProjectEvents.PrepareRemoveProjectRole(ctx, role)
	if err != nil {
		return err
	}
	aggregates = append(aggregates, agg)

	usergrants, err := repo.View.UserGrantsByProjectIDAndRoleKey(projectID, key)
	if err != nil {
		return err
	}
	for _, grant := range usergrants {
		changed := &usr_grant_model.UserGrant{
			ObjectRoot: models.ObjectRoot{AggregateID: grant.ID, Sequence: grant.Sequence, ResourceOwner: grant.ResourceOwner},
			RoleKeys:   grant.RoleKeys,
			ProjectID:  grant.ProjectID,
			UserID:     grant.UserID,
		}
		changed.RemoveRoleKeyIfExisting(key)
		_, agg, err := repo.UserGrantEvents.PrepareChangeUserGrant(ctx, changed, true)
		if err != nil {
			return err
		}
		aggregates = append(aggregates, agg)
	}
	if err != nil {
		return err
	}
	err = es_sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, project.AppendEvents, aggregates...)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProjectRepo) SearchProjectRoles(ctx context.Context, request *proj_model.ProjectRoleSearchRequest) (*proj_model.ProjectRoleSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	roles, count, err := repo.View.SearchProjectRoles(request)
	if err != nil {
		return nil, err
	}
	return &proj_model.ProjectRoleSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.ProjectRolesToModel(roles),
	}, nil
}

func (repo *ProjectRepo) ProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64) (*proj_model.ProjectChanges, error) {
	changes, err := repo.ProjectEvents.ProjectChanges(ctx, id, lastSequence, limit)
	if err != nil {
		return nil, err
	}
	return changes, nil
}

func (repo *ProjectRepo) ApplicationByID(ctx context.Context, projectID, appID string) (app *proj_model.Application, err error) {
	return repo.ProjectEvents.ApplicationByIDs(ctx, projectID, appID)
}

func (repo *ProjectRepo) AddApplication(ctx context.Context, app *proj_model.Application) (*proj_model.Application, error) {
	return repo.ProjectEvents.AddApplication(ctx, app)
}

func (repo *ProjectRepo) ChangeApplication(ctx context.Context, app *proj_model.Application) (*proj_model.Application, error) {
	return repo.ProjectEvents.ChangeApplication(ctx, app)
}

func (repo *ProjectRepo) DeactivateApplication(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	return repo.ProjectEvents.DeactivateApplication(ctx, projectID, appID)
}

func (repo *ProjectRepo) ReactivateApplication(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	return repo.ProjectEvents.ReactivateApplication(ctx, projectID, appID)
}

func (repo *ProjectRepo) RemoveApplication(ctx context.Context, projectID, appID string) error {
	app := proj_model.NewApplication(projectID, appID)
	return repo.ProjectEvents.RemoveApplication(ctx, app)
}

func (repo *ProjectRepo) SearchApplications(ctx context.Context, request *proj_model.ApplicationSearchRequest) (*proj_model.ApplicationSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	apps, count, err := repo.View.SearchApplications(request)
	if err != nil {
		return nil, err
	}
	return &proj_model.ApplicationSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.ApplicationViewsToModel(apps),
	}, nil
}

func (repo *ProjectRepo) ApplicationChanges(ctx context.Context, id string, appId string, lastSequence uint64, limit uint64) (*proj_model.ApplicationChanges, error) {
	changes, err := repo.ProjectEvents.ApplicationChanges(ctx, id, appId, lastSequence, limit)
	if err != nil {
		return nil, err
	}
	return changes, nil
}

func (repo *ProjectRepo) ChangeOIDCConfig(ctx context.Context, config *proj_model.OIDCConfig) (*proj_model.OIDCConfig, error) {
	return repo.ProjectEvents.ChangeOIDCConfig(ctx, config)
}

func (repo *ProjectRepo) ChangeOIDConfigSecret(ctx context.Context, projectID, appID string) (*proj_model.OIDCConfig, error) {
	return repo.ProjectEvents.ChangeOIDCConfigSecret(ctx, projectID, appID)
}

func (repo *ProjectRepo) ProjectGrantByID(ctx context.Context, projectID, appID string) (app *proj_model.ProjectGrant, err error) {
	return repo.ProjectEvents.ProjectGrantByIDs(ctx, projectID, appID)
}

func (repo *ProjectRepo) SearchProjectGrants(ctx context.Context, request *proj_model.ProjectGrantViewSearchRequest) (*proj_model.ProjectGrantViewSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	projects, count, err := repo.View.SearchProjectGrants(request)
	if err != nil {
		return nil, err
	}
	return &proj_model.ProjectGrantViewSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.ProjectGrantsToModel(projects),
	}, nil
}

func (repo *ProjectRepo) AddProjectGrant(ctx context.Context, grant *proj_model.ProjectGrant) (*proj_model.ProjectGrant, error) {
	return repo.ProjectEvents.AddProjectGrant(ctx, grant)
}

func (repo *ProjectRepo) ChangeProjectGrant(ctx context.Context, grant *proj_model.ProjectGrant) (*proj_model.ProjectGrant, error) {
	project, aggFunc, removedRoles, err := repo.ProjectEvents.PrepareChangeProjectGrant(ctx, grant)
	if err != nil {
		return nil, err
	}
	agg, err := aggFunc(ctx)
	if err != nil {
		return nil, err
	}
	aggregates := make([]*es_models.Aggregate, 0)
	aggregates = append(aggregates, agg)

	usergrants, err := repo.View.UserGrantsByProjectID(grant.AggregateID)
	if err != nil {
		return nil, err
	}
	for _, grant := range usergrants {
		changed := &usr_grant_model.UserGrant{
			ObjectRoot: models.ObjectRoot{AggregateID: grant.ID, Sequence: grant.Sequence, ResourceOwner: grant.ResourceOwner},
			RoleKeys:   grant.RoleKeys,
			ProjectID:  grant.ProjectID,
			UserID:     grant.UserID,
		}
		existing := changed.RemoveRoleKeysIfExisting(removedRoles)
		if existing {
			_, agg, err := repo.UserGrantEvents.PrepareChangeUserGrant(ctx, changed, true)
			if err != nil {
				return nil, err
			}
			aggregates = append(aggregates, agg)
		}
	}
	if err != nil {
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, project.AppendEvents, aggregates...)
	if err != nil {
		return nil, err
	}
	if _, g := es_proj_model.GetProjectGrant(project.Grants, grant.GrantID); g != nil {
		return es_proj_model.GrantToModel(g), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-dksi8", "Could not find app in list")
}

func (repo *ProjectRepo) DeactivateProjectGrant(ctx context.Context, projectID, grantID string) (*proj_model.ProjectGrant, error) {
	return repo.ProjectEvents.DeactivateProjectGrant(ctx, projectID, grantID)
}

func (repo *ProjectRepo) ReactivateProjectGrant(ctx context.Context, projectID, grantID string) (*proj_model.ProjectGrant, error) {
	return repo.ProjectEvents.ReactivateProjectGrant(ctx, projectID, grantID)
}

func (repo *ProjectRepo) RemoveProjectGrant(ctx context.Context, projectID, grantID string) error {
	grant, err := repo.ProjectEvents.ProjectGrantByIDs(ctx, projectID, grantID)
	if err != nil {
		return err
	}
	aggregates := make([]*es_models.Aggregate, 0)
	project, aggFunc, err := repo.ProjectEvents.PrepareRemoveProjectGrant(ctx, grant)
	if err != nil {
		return err
	}
	agg, err := aggFunc(ctx)
	if err != nil {
		return err
	}
	aggregates = append(aggregates, agg)

	usergrants, err := repo.View.UserGrantsByOrgIDAndProjectID(grant.GrantedOrgID, projectID)
	if err != nil {
		return err
	}
	for _, grant := range usergrants {
		_, grantAggregates, err := repo.UserGrantEvents.PrepareRemoveUserGrant(ctx, grant.ID, true)
		if err != nil {
			return err
		}
		for _, agg := range grantAggregates {
			aggregates = append(aggregates, agg)
		}
	}
	if err != nil {
		return err
	}
	err = es_sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, project.AppendEvents, aggregates...)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProjectRepo) ProjectGrantMemberByID(ctx context.Context, projectID, grantID, userID string) (member *proj_model.ProjectGrantMember, err error) {
	member = proj_model.NewProjectGrantMember(projectID, grantID, userID)
	return repo.ProjectEvents.ProjectGrantMemberByIDs(ctx, member)
}

func (repo *ProjectRepo) AddProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	return repo.ProjectEvents.AddProjectGrantMember(ctx, member)
}

func (repo *ProjectRepo) ChangeProjectGrantMember(ctx context.Context, member *proj_model.ProjectGrantMember) (*proj_model.ProjectGrantMember, error) {
	return repo.ProjectEvents.ChangeProjectGrantMember(ctx, member)
}

func (repo *ProjectRepo) RemoveProjectGrantMember(ctx context.Context, projectID, grantID, userID string) error {
	member := proj_model.NewProjectGrantMember(projectID, grantID, userID)
	return repo.ProjectEvents.RemoveProjectGrantMember(ctx, member)
}

func (repo *ProjectRepo) SearchProjectGrantMembers(ctx context.Context, request *proj_model.ProjectGrantMemberSearchRequest) (*proj_model.ProjectGrantMemberSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	members, count, err := repo.View.SearchProjectGrantMembers(request)
	if err != nil {
		return nil, err
	}
	return &proj_model.ProjectGrantMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.ProjectGrantMembersToModel(members),
	}, nil
}

func (repo *ProjectRepo) GetProjectMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "PROJECT") && !strings.HasPrefix(roleMap, "PROJECT_GRANT") {
			roles = append(roles, roleMap)
		}
	}
	return roles
}

func (repo *ProjectRepo) GetProjectGrantMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "PROJECT_GRANT") {
			roles = append(roles, roleMap)
		}
	}
	return roles
}
