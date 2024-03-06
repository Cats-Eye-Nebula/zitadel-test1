package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/target"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	TargetTable               = "projections.targets"
	TargetIDCol               = "id"
	TargetCreationDateCol     = "creation_date"
	TargetChangeDateCol       = "change_date"
	TargetResourceOwnerCol    = "resource_owner"
	TargetInstanceIDCol       = "instance_id"
	TargetSequenceCol         = "sequence"
	TargetNameCol             = "name"
	TargetTargetType          = "target_type"
	TargetURLCol              = "url"
	TargetTimeoutCol          = "timeout"
	TargetAsyncCol            = "async"
	TargetInterruptOnErrorCol = "interrupt_on_error"
)

type targetProjection struct{}

func newTargetProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(targetProjection))
}

func (*targetProjection) Name() string {
	return TargetTable
}

func (*targetProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(TargetIDCol, handler.ColumnTypeText),
			handler.NewColumn(TargetCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(TargetChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(TargetResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(TargetInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(TargetTargetType, handler.ColumnTypeEnum),
			handler.NewColumn(TargetSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(TargetNameCol, handler.ColumnTypeText),
			handler.NewColumn(TargetURLCol, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(TargetTimeoutCol, handler.ColumnTypeInt64, handler.Default(0)),
			handler.NewColumn(TargetAsyncCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(TargetInterruptOnErrorCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(TargetInstanceIDCol, TargetResourceOwnerCol, TargetIDCol),
		),
	)
}

func (p *targetProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: target.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  target.AddedEventType,
					Reduce: p.reduceTargetAdded,
				},
				{
					Event:  target.ChangedEventType,
					Reduce: p.reduceTargetChanged,
				},
				{
					Event:  target.RemovedEventType,
					Reduce: p.reduceTargetRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(TargetInstanceIDCol),
				},
			},
		},
	}
}

func (p *targetProjection) reduceTargetAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*target.AddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-2i7tb34fbv", "reduce.wrong.event.type% s", target.AddedEventType)
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(TargetInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(TargetResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(TargetIDCol, e.Aggregate().ID),
			handler.NewCol(TargetCreationDateCol, e.CreationDate()),
			handler.NewCol(TargetChangeDateCol, e.CreationDate()),
			handler.NewCol(TargetSequenceCol, e.Sequence()),
			handler.NewCol(TargetNameCol, e.Name),
			handler.NewCol(TargetURLCol, e.URL),
			handler.NewCol(TargetTargetType, e.TargetType),
			handler.NewCol(TargetTimeoutCol, e.Timeout),
			handler.NewCol(TargetAsyncCol, e.Async),
			handler.NewCol(TargetInterruptOnErrorCol, e.InterruptOnError),
		},
	), nil
}

func (p *targetProjection) reduceTargetChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*target.ChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-tyr9irah4l", "reduce.wrong.event.type %s", target.ChangedEventType)
	}
	values := []handler.Column{
		handler.NewCol(TargetChangeDateCol, e.CreationDate()),
		handler.NewCol(TargetSequenceCol, e.Sequence()),
	}
	if e.Name != nil {
		values = append(values, handler.NewCol(TargetNameCol, *e.Name))
	}
	if e.TargetType != nil {
		values = append(values, handler.NewCol(TargetTargetType, *e.TargetType))
	}
	if e.URL != nil {
		values = append(values, handler.NewCol(TargetURLCol, *e.URL))
	}
	if e.Timeout != nil {
		values = append(values, handler.NewCol(TargetTimeoutCol, *e.Timeout))
	}
	if e.Async != nil {
		values = append(values, handler.NewCol(TargetAsyncCol, *e.Async))
	}
	if e.InterruptOnError != nil {
		values = append(values, handler.NewCol(TargetInterruptOnErrorCol, *e.InterruptOnError))
	}
	return handler.NewUpdateStatement(
		e,
		values,
		[]handler.Condition{
			handler.NewCond(TargetInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(TargetResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCond(TargetIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *targetProjection) reduceTargetRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*target.RemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-5xy8ssr3ic", "reduce.wrong.event.type %s", target.RemovedEventType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(TargetInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(TargetResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCond(TargetIDCol, e.Aggregate().ID),
		},
	), nil
}