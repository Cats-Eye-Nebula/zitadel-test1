package project

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	uniqueProjectnameTable = "project_names"
	projectEventTypePrefix = eventstore.EventType("project.")
	ProjectAddedType       = projectEventTypePrefix + "added"
	ProjectChanged         = projectEventTypePrefix + "changed"
	ProjectDeactivated     = projectEventTypePrefix + "deactivated"
	ProjectReactivated     = projectEventTypePrefix + "reactivated"
	ProjectRemoved         = projectEventTypePrefix + "removed"
)

func NewAddProjectNameUniqueConstraint(projectName, resourceOwner string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueProjectnameTable,
		projectName+resourceOwner,
		"Errors.Project.AlreadyExists")
}

func NewRemoveProjectNameUniqueConstraint(projectName, resourceOwner string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		uniqueProjectnameTable,
		projectName+resourceOwner)
}

type ProjectAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name                 string `json:"name,omitempty"`
	ProjectRoleAssertion bool   `json:"projectRoleAssertion,omitempty"`
	ProjectRoleCheck     bool   `json:"projectRoleCheck,omitempty"`
}

func (e *ProjectAddedEvent) Data() interface{} {
	return e
}

func (e *ProjectAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddProjectNameUniqueConstraint(e.Name, e.ResourceOwner())}
}

func NewProjectAddedEvent(ctx context.Context, name, resourceOwner string) *ProjectAddedEvent {
	return &ProjectAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPushWithResourceOwner(
			ctx,
			ProjectAddedType,
			resourceOwner,
		),
		Name: name,
	}
}

func ProjectAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ProjectAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-Bfg2f", "unable to unmarshal project")
	}

	return e, nil
}
