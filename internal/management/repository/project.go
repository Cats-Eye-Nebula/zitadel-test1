package repository

import (
	"context"

	"github.com/caos/zitadel/internal/project/model"
)

type ProjectRepository interface {
	ProjectByID(ctx context.Context, id string) (*model.ProjectView, error)
	SearchProjects(ctx context.Context, request *model.ProjectViewSearchRequest) (*model.ProjectViewSearchResponse, error)
	SearchProjectGrants(ctx context.Context, request *model.ProjectGrantViewSearchRequest) (*model.ProjectGrantViewSearchResponse, error)
	SearchGrantedProjects(ctx context.Context, request *model.ProjectGrantViewSearchRequest) (*model.ProjectGrantViewSearchResponse, error)
	ProjectGrantViewByID(ctx context.Context, grantID string) (*model.ProjectGrantView, error)

	ProjectMemberByID(ctx context.Context, projectID, userID string) (*model.ProjectMemberView, error)
	SearchProjectMembers(ctx context.Context, request *model.ProjectMemberSearchRequest) (*model.ProjectMemberSearchResponse, error)
	GetProjectMemberRoles(ctx context.Context) ([]string, error)

	SearchProjectRoles(ctx context.Context, projectId string, request *model.ProjectRoleSearchRequest) (*model.ProjectRoleSearchResponse, error)
	ProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool) (*model.ProjectChanges, error)

	ApplicationByID(ctx context.Context, projectID, appID string) (*model.ApplicationView, error)
	ChangeOIDCConfig(ctx context.Context, config *model.OIDCConfig) (*model.OIDCConfig, error)
	ChangeOIDConfigSecret(ctx context.Context, projectID, appID string) (*model.OIDCConfig, error)
	SearchApplications(ctx context.Context, request *model.ApplicationSearchRequest) (*model.ApplicationSearchResponse, error)
	ApplicationChanges(ctx context.Context, id string, secId string, lastSequence uint64, limit uint64, sortAscending bool) (*model.ApplicationChanges, error)

	ProjectGrantByID(ctx context.Context, grantID string) (*model.ProjectGrantView, error)
	SearchProjectGrantMembers(ctx context.Context, request *model.ProjectGrantMemberSearchRequest) (*model.ProjectGrantMemberSearchResponse, error)

	ProjectGrantMemberByID(ctx context.Context, projectID, userID string) (*model.ProjectGrantMemberView, error)
	GetProjectGrantMemberRoles() []string
}
