package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"testing"
)

func TestAppendAddMemberEvent(t *testing.T) {
	type args struct {
		iam    *Iam
		member *IamMember
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append add member event",
			args: args{
				iam:    &Iam{},
				member: &IamMember{UserID: "UserID", Roles: []string{"Role"}},
				event:  &es_models.Event{},
			},
			result: &Iam{Members: []*IamMember{&IamMember{UserID: "UserID", Roles: []string{"Role"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.member != nil {
				data, _ := json.Marshal(tt.args.member)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddMemberEvent(tt.args.event)
			if len(tt.args.iam.Members) != 1 {
				t.Errorf("got wrong result should have one member actual: %v ", len(tt.args.iam.Members))
			}
			if tt.args.iam.Members[0] == tt.result.Members[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Members[0], tt.args.iam.Members[0])
			}
		})
	}
}

func TestAppendChangeMemberEvent(t *testing.T) {
	type args struct {
		iam    *Iam
		member *IamMember
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append change member event",
			args: args{
				iam:    &Iam{Members: []*IamMember{&IamMember{UserID: "UserID", Roles: []string{"Role"}}}},
				member: &IamMember{UserID: "UserID", Roles: []string{"ChangedRole"}},
				event:  &es_models.Event{},
			},
			result: &Iam{Members: []*IamMember{&IamMember{UserID: "UserID", Roles: []string{"ChangedRole"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.member != nil {
				data, _ := json.Marshal(tt.args.member)
				tt.args.event.Data = data
			}
			tt.args.iam.appendChangeMemberEvent(tt.args.event)
			if len(tt.args.iam.Members) != 1 {
				t.Errorf("got wrong result should have one member actual: %v ", len(tt.args.iam.Members))
			}
			if tt.args.iam.Members[0] == tt.result.Members[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Members[0], tt.args.iam.Members[0])
			}
		})
	}
}

func TestAppendRemoveMemberEvent(t *testing.T) {
	type args struct {
		iam    *Iam
		member *IamMember
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append remove member event",
			args: args{
				iam:    &Iam{Members: []*IamMember{&IamMember{UserID: "UserID", Roles: []string{"Role"}}}},
				member: &IamMember{UserID: "UserID"},
				event:  &es_models.Event{},
			},
			result: &Iam{Members: []*IamMember{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.member != nil {
				data, _ := json.Marshal(tt.args.member)
				tt.args.event.Data = data
			}
			tt.args.iam.appendRemoveMemberEvent(tt.args.event)
			if len(tt.args.iam.Members) != 0 {
				t.Errorf("got wrong result should have no member actual: %v ", len(tt.args.iam.Members))
			}
		})
	}
}
