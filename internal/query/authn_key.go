package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/caos/zitadel/internal/domain"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

var (
	authNKeyTable = table{
		name: projection.AuthNKeyTable,
	}
	AuthNKeyColumnID = Column{
		name:  projection.AuthNKeyIDCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnCreationDate = Column{
		name:  projection.AuthNKeyCreationDateCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnResourceOwner = Column{
		name:  projection.AuthNKeyResourceOwnerCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnAggregateID = Column{
		name:  projection.AuthNKeyAggregateIDCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnSequence = Column{
		name:  projection.AuthNKeySequenceCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnObjectID = Column{
		name:  projection.AuthNKeyObjectIDCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnExpiration = Column{
		name:  projection.AuthNKeyExpirationCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnIdentifier = Column{
		name:  projection.AuthNKeyIdentifierCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnPublicKey = Column{
		name:  projection.AuthNKeyPublicKeyCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnType = Column{
		name:  projection.AuthNKeyTypeCol,
		table: authNKeyTable,
	}
	AuthNKeyColumnEnabled = Column{
		name:  projection.AuthNKeyEnabledCol,
		table: authNKeyTable,
	}
)

type AuthNKeys struct {
	SearchResponse
	AuthNKeys []*AuthNKey
}

type AuthNKey struct {
	ID            string
	CreationDate  time.Time
	ResourceOwner string
	AggregateID   string
	Sequence      uint64

	ObjectID   string
	Expiration time.Time
	Identifier string
	PublicKey  []byte
	Type       domain.AuthNKeyType
	Enabled    bool
}

type AuthNKeySearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *AuthNKeySearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) SearchAuthNKeys(ctx context.Context, queries *AuthNKeySearchQueries) (authNKeys *AuthNKeys, err error) {
	query, scan := prepareAuthNKeysQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-SAf3f", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Dbg53", "Errors.Internal")
	}
	authNKeys, err = scan(rows)
	if err != nil {
		return nil, err
	}
	authNKeys.LatestSequence, err = q.latestSequence(ctx, authNKeyTable)
	return authNKeys, err
}

func (q *Queries) GetAuthNKeyByID(ctx context.Context, id string, orgID string) (*AuthNKey, error) {
	stmt, scan := prepareAuthNKeyQuery()
	query, args, err := stmt.Where(
		sq.Eq{
			AuthNKeyColumnID.identifier():            id,
			AuthNKeyColumnResourceOwner.identifier(): orgID,
		}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-DAb32", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func NewAuthNKeyResourceOwnerQuery(id string) (SearchQuery, error) {
	return NewTextQuery(AuthNKeyColumnResourceOwner, id, TextEquals)
}

func NewAuthNKeyAggregateIDQuery(id string) (SearchQuery, error) {
	return NewTextQuery(AuthNKeyColumnAggregateID, id, TextEquals)
}

func NewAuthNKeyObjectIDQuery(id string) (SearchQuery, error) {
	return NewTextQuery(AuthNKeyColumnObjectID, id, TextEquals)
}

//
//func NewAuthNKeyNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
//	return NewTextQuery(AuthNKeyColumnName, value, method)
//}
//
//func NewAuthNKeyStateSearchQuery(value domain.AuthNKeyState) (SearchQuery, error) {
//	return NewNumberQuery(AuthNKeyColumnState, int(value), NumberEquals)
//}

func prepareAuthNKeysQuery() (sq.SelectBuilder, func(rows *sql.Rows) (*AuthNKeys, error)) {
	return sq.Select(
			AuthNKeyColumnID.identifier(),
			AuthNKeyColumnCreationDate.identifier(),
			AuthNKeyColumnResourceOwner.identifier(),
			AuthNKeyColumnAggregateID.identifier(),
			AuthNKeyColumnSequence.identifier(),
			AuthNKeyColumnObjectID.identifier(),
			AuthNKeyColumnExpiration.identifier(),
			AuthNKeyColumnIdentifier.identifier(),
			AuthNKeyColumnPublicKey.identifier(),
			AuthNKeyColumnType.identifier(),
			AuthNKeyColumnEnabled.identifier(),
			countColumn.identifier(),
		).From(authNKeyTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*AuthNKeys, error) {
			authNKeys := make([]*AuthNKey, 0)
			var count uint64
			for rows.Next() {
				authNKey := new(AuthNKey)
				err := rows.Scan(
					&authNKey.ID,
					&authNKey.CreationDate,
					&authNKey.ResourceOwner,
					&authNKey.AggregateID,
					&authNKey.Sequence,
					&authNKey.ObjectID,
					&authNKey.Expiration,
					&authNKey.Identifier,
					&authNKey.PublicKey,
					&authNKey.Type,
					&authNKey.Enabled,
					&count,
				)
				if err != nil {
					return nil, err
				}
				authNKeys = append(authNKeys, authNKey)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-Dgfn3", "Errors.Query.CloseRows")
			}

			return &AuthNKeys{
				AuthNKeys: authNKeys,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func prepareAuthNKeyQuery() (sq.SelectBuilder, func(row *sql.Row) (*AuthNKey, error)) {
	return sq.Select(
			AuthNKeyColumnID.identifier(),
			AuthNKeyColumnCreationDate.identifier(),
			AuthNKeyColumnResourceOwner.identifier(),
			AuthNKeyColumnAggregateID.identifier(),
			AuthNKeyColumnSequence.identifier(),
			AuthNKeyColumnObjectID.identifier(),
			AuthNKeyColumnExpiration.identifier(),
			AuthNKeyColumnIdentifier.identifier(),
			AuthNKeyColumnPublicKey.identifier(),
			AuthNKeyColumnType.identifier(),
			AuthNKeyColumnEnabled.identifier(),
		).From(authNKeyTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*AuthNKey, error) {
			authNKey := new(AuthNKey)
			err := row.Scan(
				&authNKey.ID,
				&authNKey.CreationDate,
				&authNKey.ResourceOwner,
				&authNKey.AggregateID,
				&authNKey.Sequence,
				&authNKey.ObjectID,
				&authNKey.Expiration,
				&authNKey.Identifier,
				&authNKey.PublicKey,
				&authNKey.Type,
				&authNKey.Enabled,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-Dgr3g", "Errors.AuthNKey.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-BGnbr", "Errors.Internal")
			}
			return authNKey, nil
		}
}
