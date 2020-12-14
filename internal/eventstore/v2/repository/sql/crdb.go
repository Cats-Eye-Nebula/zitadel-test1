package sql

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strconv"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/cockroachdb/cockroach-go/v2/crdb"

	//sql import for cockroach
	_ "github.com/lib/pq"
)

const (
	//as soon as stored procedures are possible in crdb
	// we could move the code to migrations and coll the procedure
	// traking issue: https://github.com/cockroachdb/cockroach/issues/17511
	crdbInsert = "WITH data ( " +
		"    event_type, " +
		"    aggregate_type, " +
		"    aggregate_id, " +
		"    aggregate_version, " +
		"    creation_date, " +
		"    event_data, " +
		"    editor_user, " +
		"    editor_service, " +
		"    resource_owner, " +
		// variables below are calculated
		"    previous_sequence" +
		") AS (" +
		//previous_data selects the needed data of the latest event of the aggregate
		// and buffers it (crdb inmemory)
		"    WITH previous_data AS (" +
		"        SELECT COALESCE($9, MAX(event_sequence)) AS seq, resource_owner " +
		"        FROM eventstore.events " +
		//TODO: remove LIMIT 1 as soon as data cleaned up (only 1 resource_owner per aggregate)
		"        WHERE aggregate_type = $2 AND aggregate_id = $3 GROUP BY resource_owner LIMIT 1" +
		"    )" +
		// defines the data to be inserted
		"    SELECT " +
		"        $1::VARCHAR AS event_type, " +
		"        $2::VARCHAR AS aggregate_type, " +
		"        $3::VARCHAR AS aggregate_id, " +
		"        $4::VARCHAR AS aggregate_version, " +
		"        NOW() AS creation_date, " +
		"        $5::JSONB AS event_data, " +
		"        $6::VARCHAR AS editor_user, " +
		"        $7::VARCHAR AS editor_service, " +
		"        CASE WHEN EXISTS (SELECT * FROM previous_data) " +
		"            THEN (SELECT resource_owner FROM previous_data) " +
		"            ELSE $8::VARCHAR " +
		"        end AS resource_owner, " +
		"        CASE WHEN EXISTS (SELECT * FROM previous_data) " +
		"            THEN (SELECT seq FROM previous_data) " +
		"            ELSE NULL " +
		"        end AS previous_sequence" +
		") " +
		"INSERT INTO eventstore.events " +
		"	( " +
		"		event_type, " +
		"		aggregate_type," +
		"		aggregate_id, " +
		"		aggregate_version, " +
		"		creation_date, " +
		"		event_data, " +
		"		editor_user, " +
		"		editor_service, " +
		"		resource_owner, " +
		"		previous_sequence " +
		"	) " +
		"	( " +
		"		SELECT " +
		"			event_type, " +
		"			aggregate_type," +
		"			aggregate_id, " +
		"			aggregate_version, " +
		"			COALESCE(creation_date, NOW()), " +
		"			event_data, " +
		"			editor_user, " +
		"			editor_service, " +
		"			resource_owner, " +
		"			previous_sequence " +
		"		FROM data " +
		"	) " +
		"RETURNING id, event_sequence, previous_sequence, creation_date, resource_owner"
)

type CRDB struct {
	client *sql.DB
}

func NewCRDB(client *sql.DB) *CRDB {
	return &CRDB{client}
}

func (db *CRDB) Health(ctx context.Context) error { return db.client.Ping() }

// Push adds all events to the eventstreams of the aggregates.
// This call is transaction save. The transaction will be rolled back if one event fails
func (db *CRDB) Push(ctx context.Context, events ...*repository.Event) error {
	err := crdb.ExecuteTx(ctx, db.client, nil, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, crdbInsert)
		if err != nil {
			logging.Log("SQL-3to5p").WithError(err).Warn("prepare failed")
			return caos_errs.ThrowInternal(err, "SQL-OdXRE", "prepare failed")
		}

		var previousSequence Sequence
		for _, event := range events {
			if previousSequence == 0 && event.
			err = stmt.QueryRowContext(ctx,
				event.Type,
				event.AggregateType,
				event.AggregateID,
				event.Version,
				Data(event.Data),
				event.EditorUser,
				event.EditorService,
				event.ResourceOwner,
			).Scan(&event.ID, &event.Sequence, &previousSequence, &event.CreationDate, &event.ResourceOwner)

			event.PreviousSequence = uint64(previousSequence)

			if event.CheckPreviousSequence && event.PreviousSequence != uint64(previousSequence) {
				return caos_errs.ThrowAlreadyExists(nil, "SQL-k0sNg", "wrong previous sequence")
			}

			if err != nil {
				logging.LogWithFields("SQL-IP3js",
					"aggregate", event.AggregateType,
					"aggregateId", event.AggregateID,
					"aggregateType", event.AggregateType,
					"eventType", event.Type).WithError(err).Info("query failed",
					"seq", event.PreviousSequence)
				return caos_errs.ThrowInternal(err, "SQL-SBP37", "unable to create event")
			}
		}

		return nil
	})
	if err != nil && !errors.Is(err, &caos_errs.CaosError{}) {
		err = caos_errs.ThrowInternal(err, "SQL-DjgtG", "unable to store events")
	}

	return err
}

// Filter returns all events matching the given search query
func (db *CRDB) Filter(ctx context.Context, searchQuery *repository.SearchQuery) (events []*repository.Event, err error) {
	events = []*repository.Event{}
	err = query(ctx, db, searchQuery, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

//LatestSequence returns the latests sequence found by the the search query
func (db *CRDB) LatestSequence(ctx context.Context, searchQuery *repository.SearchQuery) (uint64, error) {
	var seq Sequence
	err := query(ctx, db, searchQuery, &seq)
	if err != nil {
		return 0, err
	}
	return uint64(seq), nil
}

func (db *CRDB) db() *sql.DB {
	return db.client
}

func (db *CRDB) orderByEventSequence(desc bool) string {
	if desc {
		return " ORDER BY event_sequence DESC"
	}

	return " ORDER BY event_sequence"
}

func (db *CRDB) eventQuery() string {
	return "SELECT" +
		" creation_date" +
		", event_type" +
		", event_sequence" +
		", previous_sequence" +
		", event_data" +
		", editor_service" +
		", editor_user" +
		", resource_owner" +
		", aggregate_type" +
		", aggregate_id" +
		", aggregate_version" +
		" FROM eventstore.events"
}
func (db *CRDB) maxSequenceQuery() string {
	return "SELECT MAX(event_sequence) FROM eventstore.events"
}

func (db *CRDB) columnName(col repository.Field) string {
	switch col {
	case repository.FieldAggregateID:
		return "aggregate_id"
	case repository.FieldAggregateType:
		return "aggregate_type"
	case repository.FieldSequence:
		return "event_sequence"
	case repository.FieldResourceOwner:
		return "resource_owner"
	case repository.FieldEditorService:
		return "editor_service"
	case repository.FieldEditorUser:
		return "editor_user"
	case repository.FieldEventType:
		return "event_type"
	case repository.FieldEventData:
		return "event_data"
	default:
		return ""
	}
}

func (db *CRDB) conditionFormat(operation repository.Operation) string {
	if operation == repository.OperationIn {
		return "%s %s ANY(?)"
	}
	return "%s %s ?"
}

func (db *CRDB) operation(operation repository.Operation) string {
	switch operation {
	case repository.OperationEquals, repository.OperationIn:
		return "="
	case repository.OperationGreater:
		return ">"
	case repository.OperationLess:
		return "<"
	case repository.OperationJSONContains:
		return "@>"
	}
	return ""
}

var (
	placeholder = regexp.MustCompile(`\?`)
)

//placeholder replaces all "?" with postgres placeholders ($<NUMBER>)
func (db *CRDB) placeholder(query string) string {
	occurances := placeholder.FindAllStringIndex(query, -1)
	if len(occurances) == 0 {
		return query
	}
	replaced := query[:occurances[0][0]]

	for i, l := range occurances {
		nextIDX := len(query)
		if i < len(occurances)-1 {
			nextIDX = occurances[i+1][0]
		}
		replaced = replaced + "$" + strconv.Itoa(i+1) + query[l[1]:nextIDX]
	}
	return replaced
}
