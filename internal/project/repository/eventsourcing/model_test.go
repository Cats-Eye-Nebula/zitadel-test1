package eventsourcing

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	"testing"
)

func TestProjectFromEvents(t *testing.T) {
	type args struct {
		event   []*es_models.Event
		project *Project
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "project from events, ok",
			args: args{
				event: []*es_models.Event{
					&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded},
				},
				project: &Project{Name: "ProjectName"},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{ID: "ID"}, State: int32(model.Active), Name: "ProjectName"},
		},
		{
			name: "project from events, nil project",
			args: args{
				event: []*es_models.Event{
					&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded},
				},
				project: nil,
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{ID: "ID"}, State: int32(model.Active)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.project != nil {
				data, _ := json.Marshal(tt.args.project)
				tt.args.event[0].Data = data
			}
			result, _ := ProjectFromEvents(tt.args.project, tt.args.event...)
			if result.Name != tt.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.result.Name, result.Name)
			}
		})
	}
}

func TestAppendEvent(t *testing.T) {
	type args struct {
		event   *es_models.Event
		project *Project
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append added event",
			args: args{
				event:   &es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded},
				project: &Project{Name: "ProjectName"},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{ID: "ID"}, State: int32(model.Active), Name: "ProjectName"},
		},
		{
			name: "append change event",
			args: args{
				event:   &es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectChanged},
				project: &Project{Name: "ProjectName"},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{ID: "ID"}, State: int32(model.Active), Name: "ProjectName"},
		},
		{
			name: "append deactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectDeactivated},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{ID: "ID"}, State: int32(model.Inactive)},
		},
		{
			name: "append reactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectReactivated},
			},
			result: &Project{ObjectRoot: es_models.ObjectRoot{ID: "ID"}, State: int32(model.Active)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.project != nil {
				data, _ := json.Marshal(tt.args.project)
				tt.args.event.Data = data
			}
			result := &Project{}
			result.AppendEvent(tt.args.event)
			if result.State != tt.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.result.State, result.State)
			}
			if result.Name != tt.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.result.Name, result.Name)
			}
			if result.ObjectRoot.ID != tt.result.ObjectRoot.ID {
				t.Errorf("got wrong result id: expected: %v, actual: %v ", tt.result.ObjectRoot.ID, result.ObjectRoot.ID)
			}
		})
	}
}

func TestAppendDeactivatedEvent(t *testing.T) {
	type args struct {
		project *Project
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append reactivate event",
			args: args{
				project: &Project{},
			},
			result: &Project{State: int32(model.Inactive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.project.appendDeactivatedEvent()
			if tt.args.project.State != tt.result.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.project)
			}
		})
	}
}

func TestAppendReactivatedEvent(t *testing.T) {
	type args struct {
		project *Project
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append reactivate event",
			args: args{
				project: &Project{},
			},
			result: &Project{State: int32(model.Active)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.project.appendReactivatedEvent()
			if tt.args.project.State != tt.result.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.project)
			}
		})
	}
}
