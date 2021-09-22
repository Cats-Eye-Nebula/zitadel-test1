package flow

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/action"
	"github.com/caos/zitadel/internal/repository/org"
)

type FlowProjection struct {
	crdb.StatementHandler
}

func NewFlowProjection(ctx context.Context, config crdb.StatementHandlerConfig) *FlowProjection {
	p := &FlowProjection{}
	config.ProjectionName = "projections.flows"
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *FlowProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.TriggerActionsSetEventType,
					Reduce: p.reduceTriggerActionsSetEventType,
				},
				{
					Event:  org.FlowClearedEventType,
					Reduce: p.reduceFlowClearedEventType,
				},
			},
		},
		{
			Aggregate: action.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  action.AddedEventType,
					Reduce: p.reduceFlowActionAdded,
				},
				{
					Event:  action.ChangedEventType,
					Reduce: p.reduceFlowActionChanged,
				},
				{
					Event:  action.RemovedEventType,
					Reduce: p.reduceFlowActionRemoved,
				},
			},
		},
	}
}

const (
	triggerTableSuffix           = "triggers"
	flowTypeCol                  = "flow_type"
	flowTriggerTypeCol           = "trigger_type"
	flowResourceOwnerCol         = "resource_owner"
	flowActionTriggerSequenceCol = "trigger_sequence"
	flowActionIDCol              = "action_id"

	actionTableSuffix      = "actions"
	actionIDCol            = "id"
	actionCreationDateCol  = "creation_date"
	actionChangeDateCol    = "change_date"
	actionResourceOwnerCol = "resource_owner"
	actionSequenceCol      = "sequence"
	actionNameCol          = "name"
	actionScriptCol        = "script"
)

func (p *FlowProjection) reduceTriggerActionsSetEventType(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.TriggerActionsSetEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence, "expectedType", action.AddedEventType).Error("was not an trigger actions set event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-uYq4r", "reduce.wrong.event.type")
	}
	stmts := make([]func(reader eventstore.EventReader) crdb.Exec, len(e.ActionIDs)+1)
	stmts[0] = crdb.AddDeleteStatement(
		[]handler.Condition{
			handler.NewCond(flowTypeCol, e.FlowType),
			handler.NewCond(flowTriggerTypeCol, e.TriggerType),
		},
		crdb.WithTableSuffix(triggerTableSuffix),
	)
	for i, id := range e.ActionIDs {
		stmts[i+1] = crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(flowResourceOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCol(flowTypeCol, e.FlowType),
				handler.NewCol(flowTriggerTypeCol, e.TriggerType),
				handler.NewCol(flowActionIDCol, id),
				handler.NewCol(flowActionTriggerSequenceCol, i),
			},
			crdb.WithTableSuffix(triggerTableSuffix),
		)
	}
	return crdb.NewMultiStatement(e, stmts...), nil
}

func (p *FlowProjection) reduceFlowClearedEventType(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.FlowClearedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence, "expectedType", action.AddedEventType).Error("was not an trigger actions set event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-uYq4r", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(flowTypeCol, e.FlowType),
		},
		crdb.WithTableSuffix(triggerTableSuffix),
	), nil
}

func (p *FlowProjection) reduceFlowActionAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*action.AddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence, "expectedType", action.AddedEventType).Error("was not an flow action added event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-uYq4r", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(actionIDCol, e.Aggregate().ID),
			handler.NewCol(actionCreationDateCol, e.CreationDate()),
			handler.NewCol(actionChangeDateCol, e.CreationDate()),
			handler.NewCol(actionResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(actionSequenceCol, e.Sequence()),
			handler.NewCol(actionNameCol, e.Name),
			handler.NewCol(actionScriptCol, e.Script),
		},
		crdb.WithTableSuffix(actionTableSuffix),
	), nil
}

func (p *FlowProjection) reduceFlowActionChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*action.ChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-q4oq8", "seq", event.Sequence, "expected", action.ChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Bg8oM", "reduce.wrong.event.type")
	}
	values := []handler.Column{
		handler.NewCol(actionChangeDateCol, e.CreationDate()),
		handler.NewCol(actionSequenceCol, e.Sequence()),
	}
	if e.Name != nil {
		values = append(values, handler.NewCol(actionNameCol, *e.Name))
	}
	if e.Script != nil {
		values = append(values, handler.NewCol(actionScriptCol, *e.Script))
	}
	return crdb.NewUpdateStatement(
		e,
		values,
		[]handler.Condition{
			handler.NewCond(actionIDCol, e.Aggregate().ID),
		},
		crdb.WithTableSuffix(actionTableSuffix),
	), nil
}

func (p *FlowProjection) reduceFlowActionRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*action.RemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-79OhB", "seq", event.Sequence, "expectedType", action.RemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-4TbKT", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(actionIDCol, e.Aggregate().ID),
		},
		crdb.WithTableSuffix(actionTableSuffix),
	), nil
}
