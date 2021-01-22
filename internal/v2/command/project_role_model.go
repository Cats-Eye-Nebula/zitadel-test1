package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

type ProjectRoleWriteModel struct {
	eventstore.WriteModel

	Key         string
	DisplayName string
	Group       string
	State       domain.ProjectRoleState
}

func NewProjectRoleWriteModel(projectID string, resourceOwner string) *ProjectRoleWriteModel {
	return &ProjectRoleWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *ProjectRoleWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
}

func (wm *ProjectRoleWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *project.RoleAddedEvent:
			wm.Key = e.Key
			wm.DisplayName = e.DisplayName
			wm.Group = e.Group
			wm.State = domain.ProjectRoleStateActive
		case *project.RoleChangedEvent:
			if e.Key != nil {
				wm.Key = *e.Key
			}
			if e.DisplayName != nil {
				wm.DisplayName = *e.DisplayName
			}
			if e.Group != nil {
				wm.Group = *e.Group
			}
		case *project.RoleRemovedEvent:
			wm.State = domain.ProjectRoleStateRemoved
		case *project.ProjectRemovedEvent:
			wm.State = domain.ProjectRoleStateRemoved
		}
	}
	return nil
}

func (wm *ProjectRoleWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, project.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(project.RoleAddedType, project.RoleChangedType, project.RoleRemovedType, project.ProjectRemovedType)
}

func (wm *ProjectRoleWriteModel) NewProjectRoleChangedEvent(
	ctx context.Context,
	key,
	displayName,
	group string,
) (*project.RoleChangedEvent, bool, error) {
	changes := make([]project.RoleChanges, 0)
	var err error

	if wm.Key != key {
		changes = append(changes, project.ChangeKey(key))
	}
	if wm.DisplayName != displayName {
		changes = append(changes, project.ChangeDisplayName(displayName))
	}
	if wm.Group != group {
		changes = append(changes, project.ChangeGroup(group))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := project.NewRoleChangedEvent(ctx, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
