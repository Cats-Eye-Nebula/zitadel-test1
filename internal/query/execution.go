package query

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	exec "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	executionTable = table{
		name:          projection.ExecutionTable,
		instanceIDCol: projection.ExecutionInstanceIDCol,
	}
	ExecutionColumnID = Column{
		name:  projection.ExecutionIDCol,
		table: executionTable,
	}
	ExecutionColumnChangeDate = Column{
		name:  projection.ExecutionChangeDateCol,
		table: executionTable,
	}
	ExecutionColumnInstanceID = Column{
		name:  projection.ExecutionInstanceIDCol,
		table: executionTable,
	}
	ExecutionColumnSequence = Column{
		name:  projection.ExecutionSequenceCol,
		table: executionTable,
	}

	executionTargetsTable = table{
		name:          projection.ExecutionTable + "_" + projection.ExecutionTargetSuffix,
		instanceIDCol: projection.ExecutionTargetInstanceIDCol,
	}
	executionTargetsTableAlias       = executionTargetsTable.setAlias("execution_targets")
	ExecutionTargetsColumnInstanceID = Column{
		name:  projection.ExecutionTargetInstanceIDCol,
		table: executionTargetsTableAlias,
	}
	ExecutionTargetsColumnExecutionID = Column{
		name:  projection.ExecutionTargetExecutionIDCol,
		table: executionTargetsTableAlias,
	}
	executionTargetsListCol = Column{
		name:  "targets",
		table: executionTargetsTableAlias,
	}
)

var (
	//go:embed execution_targets_join.sql
	executionTargetsJoinQuery string
	//go:embed execution_targets.sql
	executionTargetsQuery string
	//go:embed execution_targets_combined.sql
	executionTargetsCombinedQuery string
)

type Executions struct {
	SearchResponse
	Executions []*Execution
}

func (e *Executions) SetState(s *State) {
	e.State = s
}

type Execution struct {
	ID string
	domain.ObjectDetails

	Targets []*exec.Target
}

type ExecutionSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *ExecutionSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) SearchExecutions(ctx context.Context, queries *ExecutionSearchQueries) (executions *Executions, err error) {
	eq := sq.Eq{
		ExecutionColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, scan := prepareExecutionsQuery(ctx, q.client)
	return genericRowsQueryWithState[*Executions](ctx, q.client, executionTable, combineToWhereStmt(query, queries.toQuery, eq), scan)
}

func (q *Queries) GetExecutionByID(ctx context.Context, id string) (execution *Execution, err error) {
	eq := sq.Eq{
		ExecutionColumnID.identifier():         id,
		ExecutionColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, scan := prepareExecutionQuery(ctx, q.client)
	return genericRowQuery[*Execution](ctx, q.client, query.Where(eq), scan)
}

func NewExecutionInIDsSearchQuery(values []string) (SearchQuery, error) {
	return NewInTextQuery(ExecutionColumnID, values)
}

func NewExecutionTypeSearchQuery(t domain.ExecutionType) (SearchQuery, error) {
	return NewTextQuery(ExecutionColumnID, t.String(), TextStartsWith)
}

func NewTargetSearchQuery(target string) (SearchQuery, error) {
	data, err := targetItemJSONB(domain.ExecutionTargetTypeTarget, target)
	if err != nil {
		return nil, err
	}
	return NewListContains(executionTargetsListCol, data)
}

func NewIncludeSearchQuery(include string) (SearchQuery, error) {
	data, err := targetItemJSONB(domain.ExecutionTargetTypeInclude, include)
	if err != nil {
		return nil, err
	}
	return NewListContains(executionTargetsListCol, data)
}

// marshall executionTargets into the same JSONB structure as in the SQL queries
func targetItemJSONB(t domain.ExecutionTargetType, target string) ([]byte, error) {
	targets := make([]*executionTarget, 0)
	switch t {
	case domain.ExecutionTargetTypeTarget:
		targets = append(targets, &executionTarget{Target: target})
	case domain.ExecutionTargetTypeInclude:
		targets = append(targets, &executionTarget{Include: target})
	}
	return json.Marshal(targets)
}

// ExecutionTargets: provide IDs to select all target information,
func (q *Queries) ExecutionTargets(ctx context.Context, ids []string) (execution []*ExecutionTarget, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	instanceID := authz.GetInstance(ctx).InstanceID()
	if instanceID == "" {
		return nil, nil
	}

	err = q.client.QueryContext(ctx,
		func(rows *sql.Rows) error {
			execution, err = scanExecutionTargets(rows)
			return err
		},
		executionTargetsQuery,
		instanceID,
		strings.Join(ids, ","),
	)
	return execution, err
}

func (q *Queries) ExecutionTargetsCombined(ctx context.Context, ids1, ids2 []string) (execution []*ExecutionTarget, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	instanceID := authz.GetInstance(ctx).InstanceID()
	if instanceID == "" {
		return nil, nil
	}

	err = q.client.QueryContext(ctx,
		func(rows *sql.Rows) error {
			execution, err = scanExecutionTargets(rows)
			return err
		},
		executionTargetsCombinedQuery,
		instanceID,
		strings.Join(ids1, ","),
		strings.Join(ids2, ","),
	)
	return execution, err
}

func prepareExecutionQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(row *sql.Row) (*Execution, error)) {
	return sq.Select(
			ExecutionColumnInstanceID.identifier(),
			ExecutionColumnID.identifier(),
			ExecutionColumnChangeDate.identifier(),
			ExecutionColumnSequence.identifier(),
			executionTargetsListCol.identifier(),
		).From(executionTable.identifier()).
			Join("(" + executionTargetsJoinQuery + ") AS " + executionTargetsTableAlias.alias + " ON " +
				ExecutionTargetsColumnInstanceID.identifier() + " = " + ExecutionColumnInstanceID.identifier() + " AND " +
				ExecutionTargetsColumnExecutionID.identifier() + " = " + ExecutionColumnID.identifier(),
			).
			PlaceholderFormat(sq.Dollar),
		scanExecution
}

func prepareExecutionsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(rows *sql.Rows) (*Executions, error)) {
	return sq.Select(
			ExecutionColumnInstanceID.identifier(),
			ExecutionColumnID.identifier(),
			ExecutionColumnChangeDate.identifier(),
			ExecutionColumnSequence.identifier(),
			executionTargetsListCol.identifier(),
			countColumn.identifier(),
		).From(executionTable.identifier()).
			Join("(" + executionTargetsJoinQuery + ") AS " + executionTargetsTableAlias.alias + " ON " +
				ExecutionTargetsColumnInstanceID.identifier() + " = " + ExecutionColumnInstanceID.identifier() + " AND " +
				ExecutionTargetsColumnExecutionID.identifier() + " = " + ExecutionColumnID.identifier(),
			).
			PlaceholderFormat(sq.Dollar),
		scanExecutions
}

type executionTarget struct {
	Position int    `json:"position,omitempty"`
	Include  string `json:"include,omitempty"`
	Target   string `json:"target,omitempty"`
}

func scanExecution(row *sql.Row) (*Execution, error) {
	execution := new(Execution)

	var (
		targets = make([]byte, 0)
	)

	err := row.Scan(
		&execution.ResourceOwner,
		&execution.ID,
		&execution.EventDate,
		&execution.Sequence,
		&targets,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, zerrors.ThrowNotFound(err, "QUERY-qzn1xycesh", "Errors.Execution.NotFound")
		}
		return nil, zerrors.ThrowInternal(err, "QUERY-f8sjvm4tb8", "Errors.Internal")
	}

	ets := make([]*executionTarget, 0)
	if err := json.Unmarshal(targets, &ets); err != nil {
		return nil, err
	}

	execution.Targets = make([]*exec.Target, len(ets))
	for i := range ets {
		if ets[i].Target != "" {
			execution.Targets[i] = &exec.Target{Type: domain.ExecutionTargetTypeTarget, Target: ets[i].Target}
		}
		if ets[i].Include != "" {
			execution.Targets[i] = &exec.Target{Type: domain.ExecutionTargetTypeInclude, Target: ets[i].Include}
		}
	}

	return execution, nil
}

func executionTargetsUnmarshal(executionTargets []byte) ([]*exec.Target, error) {
	ets := make([]*executionTarget, 0)
	if err := json.Unmarshal(executionTargets, &ets); err != nil {
		return nil, err
	}

	targets := make([]*exec.Target, len(ets))
	// position starts with 1
	for _, et := range ets {
		if et.Target != "" {
			targets[et.Position-1] = &exec.Target{Type: domain.ExecutionTargetTypeTarget, Target: et.Target}
		}
		if et.Include != "" {
			targets[et.Position-1] = &exec.Target{Type: domain.ExecutionTargetTypeInclude, Target: et.Include}
		}
	}
	return targets, nil
}

func scanExecutions(rows *sql.Rows) (*Executions, error) {
	executions := make([]*Execution, 0)
	var count uint64

	for rows.Next() {
		execution := new(Execution)
		targets := make([]byte, 0)

		err := rows.Scan(
			&execution.ResourceOwner,
			&execution.ID,
			&execution.EventDate,
			&execution.Sequence,
			&targets,
			&count,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, zerrors.ThrowNotFound(err, "QUERY-tbrmno85vp", "Errors.Execution.NotFound")
			}
			return nil, zerrors.ThrowInternal(err, "QUERY-tyw2ydsj84", "Errors.Internal")
		}

		execution.Targets, err = executionTargetsUnmarshal(targets)
		if err != nil {
			return nil, err
		}
		executions = append(executions, execution)
	}

	if err := rows.Close(); err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-yhka3fs3mw", "Errors.Query.CloseRows")
	}

	return &Executions{
		Executions: executions,
		SearchResponse: SearchResponse{
			Count: count,
		},
	}, nil
}

type ExecutionTarget struct {
	InstanceID       string
	ExecutionID      string
	TargetID         string
	TargetType       domain.TargetType
	Endpoint         string
	Timeout          time.Duration
	InterruptOnError bool
}

func (e *ExecutionTarget) GetExecutionID() string {
	return e.ExecutionID
}
func (e *ExecutionTarget) GetTargetID() string {
	return e.TargetID
}
func (e *ExecutionTarget) IsInterruptOnError() bool {
	return e.InterruptOnError
}
func (e *ExecutionTarget) GetEndpoint() string {
	return e.Endpoint
}
func (e *ExecutionTarget) GetTargetType() domain.TargetType {
	return e.TargetType
}
func (e *ExecutionTarget) GetTimeout() time.Duration {
	return e.Timeout
}

func scanExecutionTargets(rows *sql.Rows) ([]*ExecutionTarget, error) {
	targets := make([]*ExecutionTarget, 0)
	for rows.Next() {
		target := new(ExecutionTarget)

		var (
			instanceID       = &sql.NullString{}
			executionID      = &sql.NullString{}
			targetID         = &sql.NullString{}
			targetType       = &sql.NullInt32{}
			endpoint         = &sql.NullString{}
			timeout          = &sql.NullInt64{}
			interruptOnError = &sql.NullBool{}
		)

		err := rows.Scan(
			instanceID,
			executionID,
			targetID,
			targetType,
			endpoint,
			timeout,
			interruptOnError,
		)

		if err != nil {
			return nil, err
		}

		target.InstanceID = instanceID.String
		target.ExecutionID = executionID.String
		target.TargetID = targetID.String
		target.TargetType = domain.TargetType(targetType.Int32)
		target.Endpoint = endpoint.String
		target.Timeout = time.Duration(timeout.Int64)
		target.InterruptOnError = interruptOnError.Bool

		targets = append(targets, target)
	}

	if err := rows.Close(); err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-37ardr0pki", "Errors.Query.CloseRows")
	}

	return targets, nil
}
