package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"testing"
)

func TestOrgIamPolicyAddedAggregates(t *testing.T) {
	type res struct {
		eventsCount int
		eventType   es_models.EventType
		isErr       func(error) bool
	}
	type args struct {
		ctx        context.Context
		aggCreator *es_models.AggregateCreator
		org        *model.Org
		policy     *model.OrgIamPolicy
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no policy error",
			args: args{
				ctx:        auth.NewMockContext("org", "user"),
				aggCreator: es_models.NewAggregateCreator("test"),
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "policy successful",
			args: args{
				ctx:        auth.NewMockContext("org", "user"),
				aggCreator: es_models.NewAggregateCreator("test"),
				org: &model.Org{
					ObjectRoot: es_models.ObjectRoot{
						AggregateID: "sdaf",
						Sequence:    5,
					},
				},
				policy: &model.OrgIamPolicy{
					Description:           "description",
					UserLoginMustBeDomain: true,
				},
			},
			res: res{
				eventsCount: 1,
				eventType:   model.OrgIamPolicyAdded,
				isErr:       nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg := OrgIamPolicyAddedAggregate(tt.args.aggCreator, tt.args.org, tt.args.policy)
			got, err := agg(tt.args.ctx)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got %T: %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if tt.res.isErr == nil && got.Events[0].Type != tt.res.eventType {
				t.Errorf("OrgIamPolicyAddedAggregate() event type = %v, wanted count %v", got.Events[0].Type, tt.res.eventType)
			}
			if tt.res.isErr == nil && len(got.Events) != tt.res.eventsCount {
				t.Errorf("OrgIamPolicyAddedAggregate() event count = %d, wanted count %d", len(got.Events), tt.res.eventsCount)
			}
		})
	}
}

func TestOrgIamPolicyChangedAggregates(t *testing.T) {
	type res struct {
		eventsCount int
		eventType   es_models.EventType
		isErr       func(error) bool
	}
	type args struct {
		ctx        context.Context
		aggCreator *es_models.AggregateCreator
		org        *model.Org
		policy     *model.OrgIamPolicy
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no policy error",
			args: args{
				ctx:        auth.NewMockContext("org", "user"),
				aggCreator: es_models.NewAggregateCreator("test"),
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "policy successful",
			args: args{
				ctx:        auth.NewMockContext("org", "user"),
				aggCreator: es_models.NewAggregateCreator("test"),
				org: &model.Org{
					ObjectRoot: es_models.ObjectRoot{
						AggregateID: "sdaf",
						Sequence:    5,
					},
					OrgIamPolicy: &model.OrgIamPolicy{
						Description:           "description",
						UserLoginMustBeDomain: true,
					},
				},
				policy: &model.OrgIamPolicy{
					Description:           "description",
					UserLoginMustBeDomain: false,
				},
			},
			res: res{
				eventsCount: 1,
				eventType:   model.OrgIamPolicyChanged,
				isErr:       nil,
			},
		},
		{
			name: "policy no changes",
			args: args{
				ctx:        auth.NewMockContext("org", "user"),
				aggCreator: es_models.NewAggregateCreator("test"),
				org: &model.Org{
					ObjectRoot: es_models.ObjectRoot{
						AggregateID: "sdaf",
						Sequence:    5,
					},
					OrgIamPolicy: &model.OrgIamPolicy{
						Description:           "description",
						UserLoginMustBeDomain: true,
					},
				},
				policy: &model.OrgIamPolicy{
					Description:           "description",
					UserLoginMustBeDomain: true,
				},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg := OrgIamPolicyChangedAggregate(tt.args.aggCreator, tt.args.org, tt.args.policy)
			got, err := agg(tt.args.ctx)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got %T: %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if tt.res.isErr == nil && got.Events[0].Type != tt.res.eventType {
				t.Errorf("OrgIamPolicyChangedAggregate() event type = %v, wanted count %v", got.Events[0].Type, tt.res.eventType)
			}
			if tt.res.isErr == nil && len(got.Events) != tt.res.eventsCount {
				t.Errorf("OrgIamPolicyChangedAggregate() event count = %d, wanted count %d", len(got.Events), tt.res.eventsCount)
			}
		})
	}
}

func TestOrgIamPolicyRemovedAggregates(t *testing.T) {
	type res struct {
		eventsCount int
		eventType   es_models.EventType
		isErr       func(error) bool
	}
	type args struct {
		ctx        context.Context
		aggCreator *es_models.AggregateCreator
		org        *model.Org
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "policy successful",
			args: args{
				ctx:        auth.NewMockContext("org", "user"),
				aggCreator: es_models.NewAggregateCreator("test"),
				org: &model.Org{
					ObjectRoot: es_models.ObjectRoot{
						AggregateID: "sdaf",
						Sequence:    5,
					},
					OrgIamPolicy: &model.OrgIamPolicy{
						Description:           "description",
						UserLoginMustBeDomain: true,
					},
				},
			},
			res: res{
				eventsCount: 1,
				eventType:   model.OrgIamPolicyRemoved,
				isErr:       nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg := OrgIamPolicyRemovedAggregate(tt.args.aggCreator, tt.args.org)
			got, err := agg(tt.args.ctx)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got %T: %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if tt.res.isErr == nil && got.Events[0].Type != tt.res.eventType {
				t.Errorf("OrgIamPolicyChangedAggregate() event type = %v, wanted count %v", got.Events[0].Type, tt.res.eventType)
			}
			if tt.res.isErr == nil && len(got.Events) != tt.res.eventsCount {
				t.Errorf("OrgIamPolicyChangedAggregate() event count = %d, wanted count %d", len(got.Events), tt.res.eventsCount)
			}
		})
	}
}
