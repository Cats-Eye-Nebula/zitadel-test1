package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/sethvargo/go-password/password"
)

type ProjectEventstore struct {
	es_int.Eventstore
	projectCache *ProjectCache
	pwGenerator  password.PasswordGenerator
	PasswordAlg  crypto.HashAlgorithm
}

type ProjectConfig struct {
	es_int.Eventstore
	Cache            *config.CacheConfig
	PasswordSaltCost int
}

func StartProject(conf ProjectConfig) (*ProjectEventstore, error) {
	projectCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}
	pwGenerator, err := password.NewGenerator(&password.GeneratorInput{})
	if err != nil {
		return nil, err
	}
	passwordAlg := crypto.NewBCrypt(conf.PasswordSaltCost)
	return &ProjectEventstore{
		Eventstore:   conf.Eventstore,
		projectCache: projectCache,
		pwGenerator:  pwGenerator,
		PasswordAlg:  passwordAlg,
	}, nil
}

func (es *ProjectEventstore) ProjectByID(ctx context.Context, id string) (*proj_model.Project, error) {
	project, sequence := es.projectCache.getProject(id)

	query, err := ProjectByIDQuery(project.ID, sequence)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, project.AppendEvents, query)
	if err != nil && !(caos_errs.IsNotFound(err) && project.Sequence != 0) {
		return nil, err
	}
	es.projectCache.cacheProject(project)
	return ProjectToModel(project), nil
}

func (es *ProjectEventstore) CreateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	if !project.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Name is required")
	}
	project.State = proj_model.PROJECTSTATE_ACTIVE
	repoProject := ProjectFromModel(project)

	createAggregate := ProjectCreateAggregate(es.AggregateCreator(), repoProject)
	err := es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, createAggregate)
	if err != nil {
		return nil, err
	}

	es.projectCache.cacheProject(repoProject)
	return ProjectToModel(repoProject), nil
}

func (es *ProjectEventstore) UpdateProject(ctx context.Context, project *proj_model.Project) (*proj_model.Project, error) {
	if !project.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Name is required")
	}
	existingProject, err := es.ProjectByID(ctx, project.ID)
	if err != nil {
		return nil, err
	}
	repoExisting := ProjectFromModel(existingProject)
	repoNew := ProjectFromModel(project)

	updateAggregate := ProjectUpdateAggregate(es.AggregateCreator(), repoExisting, repoNew)
	err = es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, updateAggregate)
	if err != nil {
		return nil, err
	}

	es.projectCache.cacheProject(repoExisting)
	return ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) DeactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	existing, err := es.ProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !existing.IsActive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die45", "project must be active")
	}

	repoExisting := ProjectFromModel(existing)
	aggregate := ProjectDeactivateAggregate(es.AggregateCreator(), repoExisting)
	es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)

	es.projectCache.cacheProject(repoExisting)
	return ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) ReactivateProject(ctx context.Context, id string) (*proj_model.Project, error) {
	existing, err := es.ProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.IsActive() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die45", "project must be inactive")
	}

	repoExisting := ProjectFromModel(existing)
	aggregate := ProjectReactivateAggregate(es.AggregateCreator(), repoExisting)
	es_sdk.Push(ctx, es.PushAggregates, repoExisting.AppendEvents, aggregate)

	es.projectCache.cacheProject(repoExisting)
	return ProjectToModel(repoExisting), nil
}

func (es *ProjectEventstore) ProjectMemberByIDs(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	if member.UserID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-ld93d", "userID missing")
	}
	project, err := es.ProjectByID(ctx, member.ID)
	if err != nil {
		return nil, err
	}
	for _, m := range project.Members {
		if m.UserID == member.UserID {
			return m, nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Could not find member in list")
}

func (es *ProjectEventstore) AddProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "UserID and Roles are required")
	}
	existing, err := es.ProjectByID(ctx, member.ID)
	if err != nil {
		return nil, err
	}
	if existing.ContainsMember(member) {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-idke6", "User is already member of this Project")
	}
	repoProject := ProjectFromModel(existing)
	repoMember := ProjectMemberFromModel(member)

	addAggregate := ProjectMemberAddedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, addAggregate)
	es.projectCache.cacheProject(repoProject)
	for _, m := range repoProject.Members {
		if m.UserID == member.UserID {
			return ProjectMemberToModel(m), nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Could not find member in list")
}

func (es *ProjectEventstore) ChangeProjectMember(ctx context.Context, member *proj_model.ProjectMember) (*proj_model.ProjectMember, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9dk45", "UserID and Roles are required")
	}
	existing, err := es.ProjectByID(ctx, member.ID)
	if err != nil {
		return nil, err
	}
	if !existing.ContainsMember(member) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-oe39f", "User is not member of this project")
	}
	repoProject := ProjectFromModel(existing)
	repoMember := ProjectMemberFromModel(member)

	projectAggregate := ProjectMemberChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)
	for _, m := range repoProject.Members {
		if m.UserID == member.UserID {
			return ProjectMemberToModel(m), nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Could not find member in list")
}

func (es *ProjectEventstore) RemoveProjectMember(ctx context.Context, member *proj_model.ProjectMember) error {
	if member.UserID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-d43fs", "UserID and Roles are required")
	}
	existing, err := es.ProjectByID(ctx, member.ID)
	if err != nil {
		return err
	}
	if !existing.ContainsMember(member) {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-swf34", "User is not member of this project")
	}
	repoProject := ProjectFromModel(existing)
	repoMember := ProjectMemberFromModel(member)

	projectAggregate := ProjectMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	es.projectCache.cacheProject(repoProject)
	return err
}

func (es *ProjectEventstore) AddProjectRole(ctx context.Context, role *proj_model.ProjectRole) (*proj_model.ProjectRole, error) {
	if !role.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-idue3", "Key is required")
	}
	existing, err := es.ProjectByID(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	if existing.ContainsRole(role) {
		return nil, caos_errs.ThrowAlreadyExists(nil, "EVENT-sk35t", "Project contains role with same key")
	}
	repoProject := ProjectFromModel(existing)
	repoRole := ProjectRoleFromModel(role)
	projectAggregate := ProjectRoleAddedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoRole)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}

	es.projectCache.cacheProject(repoProject)
	for _, r := range repoProject.Roles {
		if r.Key == role.Key {
			return ProjectRoleToModel(r), nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sie83", "Could not find role in list")
}

func (es *ProjectEventstore) ChangeProjectRole(ctx context.Context, role *proj_model.ProjectRole) (*proj_model.ProjectRole, error) {
	if !role.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9die3", "Key is required")
	}
	existing, err := es.ProjectByID(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	if !existing.ContainsRole(role) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-die34", "Role doesn't exist on this project")
	}
	repoProject := ProjectFromModel(existing)
	repoRole := ProjectRoleFromModel(role)
	projectAggregate := ProjectRoleChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoRole)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}

	es.projectCache.cacheProject(repoProject)
	for _, r := range repoProject.Roles {
		if r.Key == role.Key {
			return ProjectRoleToModel(r), nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-sl1or", "Could not find role in list")
}

func (es *ProjectEventstore) RemoveProjectRole(ctx context.Context, role *proj_model.ProjectRole) error {
	if role.Key == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-id823", "Key is required")
	}
	existing, err := es.ProjectByID(ctx, role.ID)
	if err != nil {
		return err
	}
	if !existing.ContainsRole(role) {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-oe823", "Role doesn't exist on project")
	}
	repoProject := ProjectFromModel(existing)
	repoRole := ProjectRoleFromModel(role)
	projectAggregate := ProjectRoleRemovedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoRole)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return nil
}

func (es *ProjectEventstore) ApplicationByIDs(ctx context.Context, projectID, appID string) (*proj_model.Application, error) {
	project, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	for _, a := range project.Applications {
		if a.AppID == appID {
			return a, nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-8ei2s", "Could not find app in list")
}

func (es *ProjectEventstore) AddApplication(ctx context.Context, app *proj_model.Application) (*proj_model.Application, error) {
	if !app.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9eidw", "Some required fields are missing")
	}
	existing, err := es.ProjectByID(ctx, app.ID)
	if err != nil {
		return nil, err
	}
	repoProject := ProjectFromModel(existing)
	repoApp := AppFromModel(app)

	addAggregate := ApplicationAddedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoApp, es.pwGenerator)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, addAggregate)
	for _, a := range repoProject.Applications {
		if a.AppID == app.AppID {
			return AppToModel(a), nil
		}
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-3udjs", "Could not find member in list")
}
