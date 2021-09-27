package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/errors"
)

const (
	currentSequencesTable = "zitadel.projections.current_sequences"
)

type LatestSequence struct {
	Sequence  uint64
	Timestamp time.Time
}

func prepareLatestSequence() (sq.SelectBuilder, func(*sql.Row) (*LatestSequence, error)) {
	return sq.Select(
			CurrentSequenceColCurrentSequence.toFullColumnName(),
			CurrentSequenceColTimestamp.toFullColumnName()).
			From(currentSequencesTable).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*LatestSequence, error) {
			seq := new(LatestSequence)
			err := row.Scan(
				&seq.Sequence,
				&seq.Timestamp,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-gmd9o", "errors.orgs.not_found")
				}
				return nil, errors.ThrowInternal(err, "QUERY-aAZ1D", "errors.internal")
			}
			return seq, nil
		}
}

func (q *Queries) latestSequence(ctx context.Context, projection table) (*LatestSequence, error) {
	query, scan := prepareLatestSequence()
	stmt, args, err := query.Where(sq.Eq{
		CurrentSequenceColProjectionName.toFullColumnName(): projection.identifier(),
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-5CfX9", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

type CurrentSequenceColumn int32

const (
	CurrentSequenceColProjectionName CurrentSequenceColumn = iota
	CurrentSequenceColAggregateType
	CurrentSequenceColCurrentSequence
	CurrentSequenceColTimestamp
)

func (c CurrentSequenceColumn) toFullColumnName() string {
	switch c {
	case CurrentSequenceColProjectionName:
		return "projection_name"
	case CurrentSequenceColAggregateType:
		return "aggregate_type"
	case CurrentSequenceColCurrentSequence:
		return "current_sequence"
	case CurrentSequenceColTimestamp:
		return "timestamp"
	default:
		return ""
	}
}
