package view

import (
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/jinzhu/gorm"
)

func ProjectRoleByIDs(db *gorm.DB, table, projectID, orgID, key string) (*model.ProjectRoleView, error) {
	role := new(model.ProjectRoleView)

	projectIDQuery := model.ProjectRoleSearchQuery{Key: proj_model.ProjectRoleSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals}
	grantIDQuery := model.ProjectRoleSearchQuery{Key: proj_model.ProjectRoleSearchKeyOrgID, Value: orgID, Method: global_model.SearchMethodEquals}
	keyQuery := model.ProjectRoleSearchQuery{Key: proj_model.ProjectRoleSearchKeyKey, Value: key, Method: global_model.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, projectIDQuery, grantIDQuery, keyQuery)
	err := query(db, role)
	return role, err
}

func ResourceOwnerProjectRolesByKey(db *gorm.DB, table, projectID, resourceOwner, key string) ([]*model.ProjectRoleView, error) {
	roles := make([]*model.ProjectRoleView, 0)
	queries := []*proj_model.ProjectRoleSearchQuery{
		&proj_model.ProjectRoleSearchQuery{Key: proj_model.ProjectRoleSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals},
		&proj_model.ProjectRoleSearchQuery{Key: proj_model.ProjectRoleSearchKeyResourceOwner, Value: resourceOwner, Method: global_model.SearchMethodEquals},
		&proj_model.ProjectRoleSearchQuery{Key: proj_model.ProjectRoleSearchKeyKey, Value: key, Method: global_model.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.ProjectRoleSearchRequest{Queries: queries})
	_, err := query(db, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func ResourceOwnerProjectRoles(db *gorm.DB, table, projectID, resourceOwner string) ([]*model.ProjectRoleView, error) {
	roles := make([]*model.ProjectRoleView, 0)
	queries := []*proj_model.ProjectRoleSearchQuery{
		&proj_model.ProjectRoleSearchQuery{Key: proj_model.ProjectRoleSearchKeyProjectID, Value: projectID, Method: global_model.SearchMethodEquals},
		&proj_model.ProjectRoleSearchQuery{Key: proj_model.ProjectRoleSearchKeyResourceOwner, Value: resourceOwner, Method: global_model.SearchMethodEquals},
	}
	query := repository.PrepareSearchQuery(table, model.ProjectRoleSearchRequest{Queries: queries})
	_, err := query(db, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func SearchProjectRoles(db *gorm.DB, table string, req *proj_model.ProjectRoleSearchRequest) ([]*model.ProjectRoleView, int, error) {
	roles := make([]*model.ProjectRoleView, 0)
	query := repository.PrepareSearchQuery(table, model.ProjectRoleSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &roles)
	if err != nil {
		return nil, 0, err
	}
	return roles, count, nil
}

func PutProjectRole(db *gorm.DB, table string, role *model.ProjectRoleView) error {
	save := repository.PrepareSave(table)
	return save(db, role)
}

func DeleteProjectRole(db *gorm.DB, table, projectID, orgID, key string) error {
	role, err := ProjectRoleByIDs(db, table, projectID, orgID, key)
	if err != nil {
		return err
	}
	delete := repository.PrepareDeleteByObject(table, role)
	return delete(db)
}
