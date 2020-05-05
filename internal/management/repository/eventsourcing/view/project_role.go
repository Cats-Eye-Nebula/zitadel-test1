package view

import (
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view"
)

const (
	projectRoleTable = "management.project_roles"
)

func (v *View) ProjectRoleByIDs(projectID, orgID, key string) (*model.ProjectRoleView, error) {
	return view.ProjectRoleByIDs(v.Db, projectRoleTable, projectID, orgID, key)
}

func (v *View) ResourceOwnerProjectRolesByKey(projectID, resourceowner, key string) ([]*model.ProjectRoleView, error) {
	return view.ResourceOwnerProjectRolesByKey(v.Db, projectRoleTable, projectID, resourceowner, key)
}

func (v *View) ResourceOwnerProjectRoles(projectID, resourceowner string) ([]*model.ProjectRoleView, error) {
	return view.ResourceOwnerProjectRoles(v.Db, projectRoleTable, projectID, resourceowner)
}

func (v *View) SearchProjectRoles(request *proj_model.ProjectRoleSearchRequest) ([]*model.ProjectRoleView, int, error) {
	return view.SearchProjectRoles(v.Db, projectRoleTable, request)
}

func (v *View) PutProjectRole(project *model.ProjectRoleView) error {
	err := view.PutProjectRole(v.Db, projectRoleTable, project)
	if err != nil {
		return err
	}
	return v.ProcessedProjectRoleSequence(project.Sequence)
}

func (v *View) DeleteProjectRole(projectID, orgID, key string, eventSequence uint64) error {
	err := view.DeleteProjectRole(v.Db, projectRoleTable, projectID, orgID, key)
	if err != nil {
		return nil
	}
	return v.ProcessedProjectRoleSequence(eventSequence)
}

func (v *View) GetLatestProjectRoleSequence() (uint64, error) {
	return v.latestSequence(projectRoleTable)
}

func (v *View) ProcessedProjectRoleSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(projectRoleTable, eventSequence)
}

func (v *View) GetLatestProjectRoleFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(projectRoleTable, sequence)
}

func (v *View) ProcessedProjectRoleFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
