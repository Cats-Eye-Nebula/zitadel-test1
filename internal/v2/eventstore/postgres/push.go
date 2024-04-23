package postgres

import (
	"context"
	"database/sql"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/v2/database"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// TODO: option to pass tx
// Push implements eventstore.Pusher.
func (s *Storage) Push(ctx context.Context, pushIntents ...eventstore.PushIntent) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	tx, err := s.client.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		return err
	}
	defer func() {
		err = database.CloseTx(tx, err)
	}()

	var previousAppName string
	if err := tx.QueryRowContext(ctx, "select current_setting('application_name')").Scan(&previousAppName); err != nil {
		logging.WithError(err).Debug("getting app name failed")
		return zerrors.ThrowInternal(err, "POSTG-qkGKk", "Errors.Internal")
	}
	// allows smaller wait times on query side for instances which are not actively writing
	if err := setAppName(ctx, tx, "es_pusher_"+authz.GetInstance(ctx).InstanceID()); err != nil {
		return err
	}
	defer func() {
		setErr := setAppName(ctx, tx, previousAppName)
		logging.OnError(setErr).Warn("resetting the app name failed")
	}()

	intents, err := lockAggregates(ctx, tx, pushIntents)
	if err != nil {
		return err
	}

	if !checkSequences(intents) {
		return zerrors.ThrowInvalidArgument(nil, "POSTG-KOM6E", "Errors.Internal.Eventstore.SequenceNotMatched")
	}

	var commands []*command
	for _, intent := range intents {
		additionalCommands, err := intentToCommands(intent)
		if err != nil {
			return err
		}
		commands = append(commands, additionalCommands...)
	}

	err = uniqueConstraints(ctx, tx, commands)
	if err != nil {
		return err
	}

	return push(ctx, tx, commands)
}

func setAppName(ctx context.Context, tx *sql.Tx, name string) error {
	_, err := tx.ExecContext(ctx, "select current_setting('application_name')")
	if err != nil {
		logging.WithFields("name", name).WithError(err).Debug("setting app name failed")
		return zerrors.ThrowInternal(err, "POSTG-G3OmZ", "Errors.Internal")
	}

	return nil
}

func lockAggregates(ctx context.Context, tx *sql.Tx, pushIntents []eventstore.PushIntent) (_ []*intent, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var stmt database.Statement

	stmt.WriteString("WITH existing AS (")
	for i, intent := range pushIntents {
		stmt.WriteString(`(SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = `)
		stmt.WriteArgs(intent.Aggregate().Instance)
		stmt.WriteString(` AND aggregate_type = `)
		stmt.WriteArgs(intent.Aggregate().Type)
		stmt.WriteString(` AND aggregate_id = `)
		stmt.WriteArgs(intent.Aggregate().ID)
		stmt.WriteString(` AND owner = `)
		stmt.WriteArgs(intent.Aggregate().Owner)
		stmt.WriteString(` ORDER BY "sequence" DESC LIMIT 1)`)

		if i < len(pushIntents)-1 {
			stmt.WriteString(" UNION ALL ")
		}
	}
	stmt.WriteString(") SELECT e.instance_id, e.owner, e.aggregate_type, e.aggregate_id, e.sequence FROM eventstore.events2 e JOIN existing ON e.instance_id = existing.instance_id AND e.aggregate_type = existing.aggregate_type AND e.aggregate_id = existing.aggregate_id AND e.sequence = existing.sequence FOR UPDATE")

	rows, err := tx.QueryContext(ctx, stmt.String(), stmt.Args()...)
	if err != nil {
		return nil, err
	}

	res := makeIntents(pushIntents)

	err = database.MapRowsToObject(rows, func(scan func(dest ...any) error) error {
		var sequence sql.Null[uint32]
		agg := new(eventstore.Aggregate)

		err := scan(
			&agg.Instance,
			&agg.Owner,
			&agg.Type,
			&agg.ID,
			&sequence,
		)
		if err != nil {
			return err
		}

		intentByAggregate(res, agg).sequence = sequence.V

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func push(ctx context.Context, tx *sql.Tx, commands []*command) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var stmt database.Statement

	stmt.WriteString(`INSERT INTO eventstore.events2 (instance_id, "owner", aggregate_type, aggregate_id, revision, creator, event_type, payload, "sequence", in_tx_order, created_at, "position") VALUES `)
	for i, cmd := range commands {
		cmd.Position.InPositionOrder = uint32(i)
		stmt.WriteString(`(`)
		stmt.WriteArgs(
			cmd.Aggregate.Instance,
			cmd.Aggregate.Owner,
			cmd.Aggregate.Type,
			cmd.Aggregate.ID,
			cmd.Revision,
			cmd.Creator,
			cmd.Type,
			cmd.Payload,
			cmd.Sequence,
			i,
		)
		stmt.WriteString(", statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())")
		stmt.WriteString(`)`)
		if i < len(commands)-1 {
			stmt.WriteString(", ")
		}
	}
	stmt.WriteString(` RETURNING created_at, "position"`)

	rows, err := tx.QueryContext(ctx, stmt.String(), stmt.Args()...)
	if err != nil {
		return err
	}

	var i int
	return database.MapRowsToObject(rows, func(scan func(dest ...any) error) error {
		defer func() { i++ }()

		err := scan(
			&commands[i].CreatedAt,
			&commands[i].Position.Position,
		)
		if err != nil {
			return err
		}
		reducer, ok := commands[i].intent.PushIntent.(eventstore.PushIntentReducer)
		if !ok {
			return nil
		}
		return reducer.Reduce(commands[i].Event)
	})
}

func uniqueConstraints(ctx context.Context, tx *sql.Tx, commands []*command) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var stmt database.Statement

	for _, cmd := range commands {
		if len(cmd.uniqueConstraints) == 0 {
			continue
		}
		for _, constraint := range cmd.uniqueConstraints {
			stmt.Reset()
			instance := cmd.Aggregate.Instance
			if constraint.IsGlobal {
				instance = ""
			}
			switch constraint.Action {
			case eventstore.UniqueConstraintAdd:
				stmt.WriteString(`INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field) VALUES (`)
				stmt.WriteArgs(instance, constraint.UniqueType, constraint.UniqueField)
				stmt.WriteString(`)`)
			case eventstore.UniqueConstraintInstanceRemove:
				stmt.WriteString(`DELETE FROM eventstore.unique_constraints WHERE instance_id = `)
				stmt.WriteArgs(instance)
			case eventstore.UniqueConstraintRemove:
				stmt.WriteString(`DELETE FROM eventstore.unique_constraints WHERE `)
				stmt.AppendArgs(
					sql.Named("@instanceId", instance),
					sql.Named("@uniqueType", constraint.UniqueType),
					sql.Named("@uniqueField", constraint.UniqueField),
				)
				stmt.WriteString(deleteUniqueConstraintClause)
			}
			_, err := tx.ExecContext(ctx, stmt.String(), stmt.Args()...)
			if err != nil {
				logging.WithFields("action", constraint.Action).Warn("handling of unique constraint failed")
				errMessage := constraint.ErrorMessage
				if errMessage == "" {
					errMessage = "Errors.Internal"
				}
				return zerrors.ThrowAlreadyExists(err, "POSTG-QzjyP", errMessage)
			}
		}
	}

	return nil
}

// the query is so complex because we accidentally stored unique constraint case sensitive
// the query checks first if there is a case sensitive match and afterwards if there is a case insensitive match
var deleteUniqueConstraintClause = `
(instance_id = @instanceId AND unique_type = @uniqueType AND unique_field = (
    SELECT unique_field from (
        SELECT instance_id, unique_type, unique_field
        FROM eventstore.unique_constraints
        WHERE instance_id = @instanceId AND unique_type = @uniqueType AND unique_field = @uniqueField
    UNION ALL
        SELECT instance_id, unique_type, unique_field
        FROM eventstore.unique_constraints
        WHERE instance_id = @instanceId AND unique_type = @uniqueType AND unique_field = LOWER(@uniqueField)
    ) AS case_insensitive_constraints LIMIT 1)
)`
