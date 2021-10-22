package crdb

import (
	"strconv"
	"strings"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
)

type execOption func(*execConfig)
type execConfig struct {
	tableName string

	args []interface{}
	err  error
}

func WithTableSuffix(name string) func(*execConfig) {
	return func(o *execConfig) {
		o.tableName += "_" + name
	}
}

func NewCreateStatement(event eventstore.EventReader, values []handler.Column, opts ...execOption) *handler.Statement {
	cols, params, args := columnsToQuery(values)
	columnNames := strings.Join(cols, ", ")
	valuesPlaceholder := strings.Join(params, ", ")

	config := execConfig{
		args: args,
	}

	if len(values) == 0 {
		config.err = handler.ErrNoValues
	}

	q := func(config execConfig) string {
		return "INSERT INTO " + config.tableName + " (" + columnNames + ") VALUES (" + valuesPlaceholder + ")"
	}

	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		Execute:          exec(config, q, opts),
	}
}

func NewUpsertStatement(event eventstore.EventReader, values []handler.Column, opts ...execOption) *handler.Statement {
	cols, params, args := columnsToQuery(values)
	columnNames := strings.Join(cols, ", ")
	valuesPlaceholder := strings.Join(params, ", ")

	config := execConfig{
		args: args,
	}

	if len(values) == 0 {
		config.err = handler.ErrNoValues
	}

	q := func(config execConfig) string {
		return "UPSERT INTO " + config.tableName + " (" + columnNames + ") VALUES (" + valuesPlaceholder + ")"
	}

	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		Execute:          exec(config, q, opts),
	}
}

func NewUpdateStatement(event eventstore.EventReader, values []handler.Column, conditions []handler.Condition, opts ...execOption) *handler.Statement {
	cols, params, args := columnsToQuery(values)
	wheres, whereArgs := conditionsToWhere(conditions, len(params))
	args = append(args, whereArgs...)

	columnNames := strings.Join(cols, ", ")
	valuesPlaceholder := strings.Join(params, ", ")
	wheresPlaceholders := strings.Join(wheres, " AND ")

	config := execConfig{
		args: args,
	}

	if len(values) == 0 {
		config.err = handler.ErrNoValues
	}

	if len(conditions) == 0 {
		config.err = handler.ErrNoCondition
	}

	q := func(config execConfig) string {
		return "UPDATE " + config.tableName + " SET (" + columnNames + ") = (" + valuesPlaceholder + ") WHERE " + wheresPlaceholders
	}

	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		Execute:          exec(config, q, opts),
	}
}

func NewDeleteStatement(event eventstore.EventReader, conditions []handler.Condition, opts ...execOption) *handler.Statement {
	wheres, args := conditionsToWhere(conditions, 0)

	wheresPlaceholders := strings.Join(wheres, " AND ")

	config := execConfig{
		args: args,
	}

	if len(conditions) == 0 {
		config.err = handler.ErrNoCondition
	}

	q := func(config execConfig) string {
		return "DELETE FROM " + config.tableName + " WHERE " + wheresPlaceholders
	}

	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		Execute:          exec(config, q, opts),
	}
}

func NewNoOpStatement(event eventstore.EventReader) *handler.Statement {
	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
	}
}

func NewMultiStatement(event eventstore.EventReader, opts ...func(eventstore.EventReader) Exec) *handler.Statement {
	if len(opts) == 0 {
		return NewNoOpStatement(event)
	}
	execs := make([]Exec, len(opts))
	for i, opt := range opts {
		execs[i] = opt(event)
	}
	return &handler.Statement{
		AggregateType:    event.Aggregate().Type,
		Sequence:         event.Sequence(),
		PreviousSequence: event.PreviousAggregateTypeSequence(),
		Execute:          multiExec(execs),
	}
}

type Exec func(ex handler.Executer, projectionName string) error

func AddCreateStatement(columns []handler.Column, opts ...execOption) func(eventstore.EventReader) Exec {
	return func(event eventstore.EventReader) Exec {
		return NewCreateStatement(event, columns, opts...).Execute
	}
}

func AddUpsertStatement(values []handler.Column, opts ...execOption) func(eventstore.EventReader) Exec {
	return func(event eventstore.EventReader) Exec {
		return NewUpsertStatement(event, values, opts...).Execute
	}
}

func AddUpdateStatement(values []handler.Column, conditions []handler.Condition, opts ...execOption) func(eventstore.EventReader) Exec {
	return func(event eventstore.EventReader) Exec {
		return NewUpdateStatement(event, values, conditions, opts...).Execute
	}
}

func AddDeleteStatement(conditions []handler.Condition, opts ...execOption) func(eventstore.EventReader) Exec {
	return func(event eventstore.EventReader) Exec {
		return NewDeleteStatement(event, conditions, opts...).Execute
	}
}

func NewArrayAppendCol(column string, value interface{}) handler.Column {
	return handler.Column{
		Name:  column,
		Value: value,
		ParameterOpt: func(placeholder string) string {
			return "array_append(" + column + ", " + placeholder + ")"
		},
	}
}

func NewArrayRemoveCol(column string, value interface{}) handler.Column {
	return handler.Column{
		Name:  column,
		Value: value,
		ParameterOpt: func(placeholder string) string {
			return "array_remove(" + column + ", " + placeholder + ")"
		},
	}
}

func columnsToQuery(cols []handler.Column) (names []string, parameters []string, values []interface{}) {
	names = make([]string, len(cols))
	values = make([]interface{}, len(cols))
	parameters = make([]string, len(cols))
	for i, col := range cols {
		names[i] = col.Name
		values[i] = col.Value
		parameters[i] = "$" + strconv.Itoa(i+1)
		if col.ParameterOpt != nil {
			parameters[i] = col.ParameterOpt(parameters[i])
		}
	}
	return names, parameters, values
}

func conditionsToWhere(cols []handler.Condition, paramOffset int) (wheres []string, values []interface{}) {
	wheres = make([]string, len(cols))
	values = make([]interface{}, len(cols))

	for i, col := range cols {
		wheres[i] = "(" + col.Name + " = $" + strconv.Itoa(i+1+paramOffset) + ")"
		values[i] = col.Value
	}

	return wheres, values
}

type query func(config execConfig) string

func exec(config execConfig, q query, opts []execOption) Exec {
	return func(ex handler.Executer, projectionName string) error {
		if projectionName == "" {
			return handler.ErrNoProjection
		}

		if config.err != nil {
			return config.err
		}

		config.tableName = projectionName
		for _, opt := range opts {
			opt(&config)
		}

		if _, err := ex.Exec(q(config), config.args...); err != nil {
			return errors.ThrowInternal(err, "CRDB-pKtsr", "exec failed")
		}

		return nil
	}
}

func multiExec(execList []Exec) Exec {
	return func(ex handler.Executer, projectionName string) error {
		for _, exec := range execList {
			if err := exec(ex, projectionName); err != nil {
				return err
			}
		}
		return nil
	}
}
