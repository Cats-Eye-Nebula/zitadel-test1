package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

type ProjectWriteModel struct {
	eventstore.WriteModel

	Name                 string
	ProjectRoleAssertion bool
	ProjectRoleCheck     bool
	State                domain.ProjectState
}

func NewProjectWriteModel(projectID string) *ProjectWriteModel {
	return &ProjectWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: projectID,
		},
	}
}

func (wm *ProjectWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *project.ProjectAddedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *ProjectWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *project.ProjectAddedEvent:
			wm.Name = e.Name
			wm.State = domain.ProjectStateActive
			//case *project.ProjectChangedEvent:
			//	wm.Name = e.Name
		}
	}
	return nil
}

func (wm *ProjectWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, project.AggregateType).
		AggregateIDs(wm.AggregateID)
}

func ProjectAggregateFromWriteModel(wm *eventstore.WriteModel) *project.Aggregate {
	return &project.Aggregate{
		Aggregate: *eventstore.AggregateFromWriteModel(wm, project.AggregateType, project.AggregateVersion),
	}
}
