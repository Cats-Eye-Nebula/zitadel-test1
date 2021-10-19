package query

import (
	"context"
	"database/sql"
	errs "errors"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
	"github.com/lib/pq"
)

type IDP struct {
	CreationDate  time.Time
	ChangeDate    time.Time
	Sequence      uint64
	ResourceOwner string
	ID            string
	State         domain.IDPConfigState
	Name          string
	StylingType   domain.IDPConfigStylingType
	OwnerType     domain.IdentityProviderType
	AutoRegister  bool
	*OIDCIDP
	*JWTIDP
}

type IDPs struct {
	SearchResponse
	IDPs []*IDP
}

type OIDCIDP struct {
	IDPID                 string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	Issuer                string
	Scopes                []string
	DisplayNameMapping    domain.OIDCMappingField
	UsernameMapping       domain.OIDCMappingField
	AuthorizationEndpoint string
	TokenEndpoint         string
}

type JWTIDP struct {
	IDPID        string
	Issuer       string
	KeysEndpoint string
	HeaderName   string
	Endpoint     string
}

//IDPByIDAndResourceOwner searches for the requested id in the context of the resource owner and IAM
func (q *Queries) IDPByIDAndResourceOwner(ctx context.Context, id, resourceOwner string) (*IDP, error) {
	stmt, scan := prepareIDPByIDQuery()
	query, args, err := stmt.Where(
		sq.And{
			sq.Eq{
				IDPIDCol.identifier(): id,
			},
			sq.Or{
				sq.Eq{
					IDPResourceOwnerCol.identifier(): resourceOwner,
				},
				sq.Eq{
					IDPResourceOwnerCol.identifier(): q.iamID,
				},
			},
		},
	).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-0gocI", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

//SearchIDPs searches executes the query in the context of the resource owner and IAM
func (q *Queries) SearchIDPs(ctx context.Context, resourceOwner string, queries *IDPSearchQueries) (idps *IDPs, err error) {
	query, scan := prepareIDPsQuery()
	query = queries.toQuery(query)
	query = query.Where(
		sq.Or{
			sq.Eq{
				IDPResourceOwnerCol.identifier(): resourceOwner,
				IDPResourceOwnerCol.identifier(): q.iamID,
			},
		},
	)
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-zC6gk", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		log.Println(err)
		return nil, errors.ThrowInternal(err, "QUERY-YTug9", "Errors.Internal")
	}
	idps, err = scan(rows)
	if err != nil {
		return nil, err
	}
	idps.LatestSequence, err = q.latestSequence(ctx, idpTable)
	return idps, err
}

type IDPSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func NewIDPIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(IDPIDCol, id, TextEquals)
}

func NewIDPOwnerTypeSearchQuery(ownerType domain.IdentityProviderType) (SearchQuery, error) {
	switch ownerType {
	case domain.IdentityProviderTypeOrg:
		return NewBoolQuery(LoginPolicyColumnIsDefault, false)
	case domain.IdentityProviderTypeSystem:
		return NewBoolQuery(LoginPolicyColumnIsDefault, true)
	default:
		return nil, errors.ThrowUnimplemented(nil, "QUERY-8yZAI", "Errors.Query.InvalidRequest")
	}
}

func NewIDPNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(IDPNameCol, value, method)
}

var (
	idpTable = table{
		name: projection.IDPTable,
	}
	IDPIDCol = Column{
		name:  projection.IDPIDCol,
		table: idpTable,
	}
	IDPCreationDateCol = Column{
		name:  projection.IDPCreationDateCol,
		table: idpTable,
	}
	IDPChangeDateCol = Column{
		name:  projection.IDPChangeDateCol,
		table: idpTable,
	}
	IDPSequenceCol = Column{
		name:  projection.IDPSequenceCol,
		table: idpTable,
	}
	IDPResourceOwnerCol = Column{
		name:  projection.IDPResourceOwnerCol,
		table: idpTable,
	}
	IDPStateCol = Column{
		name:  projection.IDPStateCol,
		table: idpTable,
	}
	IDPNameCol = Column{
		name:  projection.IDPNameCol,
		table: idpTable,
	}
	IDPStylingTypeCol = Column{
		name:  projection.IDPStylingTypeCol,
		table: idpTable,
	}
	IDPOwnerCol = Column{
		name:  projection.IDPOwnerTypeCol,
		table: idpTable,
	}
	IDPAutoRegisterCol = Column{
		name:  projection.IDPAutoRegisterCol,
		table: idpTable,
	}
)

var (
	oidcIDPTable = table{
		name: projection.IDPOIDCTable,
	}
	OIDCIDPColIDPID = Column{
		name:  projection.OIDCConfigIDPIDCol,
		table: oidcIDPTable,
	}
	OIDCIDPColClientID = Column{
		name:  projection.OIDCConfigClientIDCol,
		table: oidcIDPTable,
	}
	OIDCIDPColClientSecret = Column{
		name:  projection.OIDCConfigClientSecretCol,
		table: oidcIDPTable,
	}
	OIDCIDPColIssuer = Column{
		name:  projection.OIDCConfigIssuerCol,
		table: oidcIDPTable,
	}
	OIDCIDPColScopes = Column{
		name:  projection.OIDCConfigScopesCol,
		table: oidcIDPTable,
	}
	OIDCIDPColDisplayNameMapping = Column{
		name:  projection.OIDCConfigDisplayNameMappingCol,
		table: oidcIDPTable,
	}
	OIDCIDPColUsernameMapping = Column{
		name:  projection.OIDCConfigUsernameMappingCol,
		table: oidcIDPTable,
	}
	OIDCIDPColAuthorizationEndpoint = Column{
		name:  projection.OIDCConfigAuthorizationEndpointCol,
		table: oidcIDPTable,
	}
	OIDCIDPColTokenEndpoint = Column{
		name:  projection.OIDCConfigTokenEndpointCol,
		table: oidcIDPTable,
	}
)

var (
	jwtIDPTable = table{
		name: projection.IDPJWTTable,
	}
	JWTIDPColIDPID = Column{
		name:  projection.JWTConfigIDPIDCol,
		table: jwtIDPTable,
	}
	JWTIDPColIssuer = Column{
		name:  projection.JWTConfigIssuerCol,
		table: jwtIDPTable,
	}
	JWTIDPColKeysEndpoint = Column{
		name:  projection.JWTConfigKeysEndpointCol,
		table: jwtIDPTable,
	}
	JWTIDPColHeaderName = Column{
		name:  projection.JWTConfigHeaderNameCol,
		table: jwtIDPTable,
	}
	JWTIDPColEndpoint = Column{
		name:  projection.JWTConfigEndpointCol,
		table: jwtIDPTable,
	}
)

func (q *IDPSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.ToQuery(query)
	}
	return query
}

func prepareIDPByIDQuery() (sq.SelectBuilder, func(*sql.Row) (*IDP, error)) {
	return sq.Select(
			IDPIDCol.identifier(),
			IDPResourceOwnerCol.identifier(),
			IDPCreationDateCol.identifier(),
			IDPChangeDateCol.identifier(),
			IDPSequenceCol.identifier(),
			IDPStateCol.identifier(),
			IDPNameCol.identifier(),
			IDPStylingTypeCol.identifier(),
			IDPOwnerCol.identifier(),
			IDPAutoRegisterCol.identifier(),
			OIDCIDPColIDPID.identifier(),
			OIDCIDPColClientID.identifier(),
			OIDCIDPColClientSecret.identifier(),
			OIDCIDPColIssuer.identifier(),
			OIDCIDPColScopes.identifier(),
			OIDCIDPColDisplayNameMapping.identifier(),
			OIDCIDPColUsernameMapping.identifier(),
			OIDCIDPColAuthorizationEndpoint.identifier(),
			OIDCIDPColTokenEndpoint.identifier(),
			JWTIDPColIDPID.identifier(),
			JWTIDPColIssuer.identifier(),
			JWTIDPColKeysEndpoint.identifier(),
			JWTIDPColHeaderName.identifier(),
			JWTIDPColEndpoint.identifier(),
		).From(idpTable.identifier()).
			LeftJoin(join(OIDCIDPColIDPID, IDPIDCol)).
			LeftJoin(join(JWTIDPColIDPID, IDPIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*IDP, error) {
			idp := new(IDP)

			oidcIDPID := sql.NullString{}
			oidcClientID := sql.NullString{}
			oidcClientSecret := new(crypto.CryptoValue)
			oidcIssuer := sql.NullString{}
			oidcScopes := pq.StringArray{}
			oidcDisplayNameMapping := sql.NullInt32{}
			oidcUsernameMapping := sql.NullInt32{}
			oidcAuthorizationEndpoint := sql.NullString{}
			oidcTokenEndpoint := sql.NullString{}

			jwtIDPID := sql.NullString{}
			jwtIssuer := sql.NullString{}
			jwtKeysEndpoint := sql.NullString{}
			jwtHeaderName := sql.NullString{}
			jwtEndpoint := sql.NullString{}

			err := row.Scan(
				&idp.ID,
				&idp.ResourceOwner,
				&idp.CreationDate,
				&idp.ChangeDate,
				&idp.Sequence,
				&idp.State,
				&idp.Name,
				&idp.StylingType,
				&idp.OwnerType,
				&idp.AutoRegister,
				&oidcIDPID,
				&oidcClientID,
				oidcClientSecret,
				&oidcIssuer,
				&oidcScopes,
				&oidcDisplayNameMapping,
				&oidcUsernameMapping,
				&oidcAuthorizationEndpoint,
				&oidcTokenEndpoint,
				&jwtIDPID,
				&jwtIssuer,
				&jwtKeysEndpoint,
				&jwtHeaderName,
				&jwtEndpoint,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-rhR2o", "Errors.IDPConfig.NotExisting")
				}
				return nil, errors.ThrowInternal(err, "QUERY-zE3Ro", "Errors.Internal")
			}

			if oidcIDPID.Valid {
				idp.OIDCIDP = &OIDCIDP{
					IDPID:                 oidcIDPID.String,
					ClientID:              oidcClientID.String,
					ClientSecret:          oidcClientSecret,
					Issuer:                oidcIssuer.String,
					Scopes:                oidcScopes,
					DisplayNameMapping:    domain.OIDCMappingField(oidcDisplayNameMapping.Int32),
					UsernameMapping:       domain.OIDCMappingField(oidcUsernameMapping.Int32),
					AuthorizationEndpoint: oidcAuthorizationEndpoint.String,
					TokenEndpoint:         oidcTokenEndpoint.String,
				}
			} else if jwtIDPID.Valid {
				idp.JWTIDP = &JWTIDP{
					IDPID:        jwtIDPID.String,
					Issuer:       jwtIssuer.String,
					KeysEndpoint: jwtKeysEndpoint.String,
					HeaderName:   jwtHeaderName.String,
					Endpoint:     jwtEndpoint.String,
				}
			}

			return idp, nil
		}
}

func prepareIDPsQuery() (sq.SelectBuilder, func(*sql.Rows) (*IDPs, error)) {
	return sq.Select(
			IDPIDCol.identifier(),
			IDPResourceOwnerCol.identifier(),
			IDPCreationDateCol.identifier(),
			IDPChangeDateCol.identifier(),
			IDPSequenceCol.identifier(),
			IDPStateCol.identifier(),
			IDPNameCol.identifier(),
			IDPStylingTypeCol.identifier(),
			IDPOwnerCol.identifier(),
			IDPAutoRegisterCol.identifier(),
			OIDCIDPColIDPID.identifier(),
			OIDCIDPColClientID.identifier(),
			OIDCIDPColClientSecret.identifier(),
			OIDCIDPColIssuer.identifier(),
			OIDCIDPColScopes.identifier(),
			OIDCIDPColDisplayNameMapping.identifier(),
			OIDCIDPColUsernameMapping.identifier(),
			OIDCIDPColAuthorizationEndpoint.identifier(),
			OIDCIDPColTokenEndpoint.identifier(),
			JWTIDPColIDPID.identifier(),
			JWTIDPColIssuer.identifier(),
			JWTIDPColKeysEndpoint.identifier(),
			JWTIDPColHeaderName.identifier(),
			JWTIDPColEndpoint.identifier(),
			countColumn.identifier(),
		).From(idpTable.identifier()).
			LeftJoin(join(OIDCIDPColIDPID, IDPIDCol)).
			LeftJoin(join(JWTIDPColIDPID, IDPIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*IDPs, error) {
			idps := make([]*IDP, 0)
			var count uint64
			for rows.Next() {
				idp := new(IDP)

				oidcIDPID := sql.NullString{}
				oidcClientID := sql.NullString{}
				oidcClientSecret := new(crypto.CryptoValue)
				oidcIssuer := sql.NullString{}
				oidcScopes := pq.StringArray{}
				oidcDisplayNameMapping := sql.NullInt32{}
				oidcUsernameMapping := sql.NullInt32{}
				oidcAuthorizationEndpoint := sql.NullString{}
				oidcTokenEndpoint := sql.NullString{}

				jwtIDPID := sql.NullString{}
				jwtIssuer := sql.NullString{}
				jwtKeysEndpoint := sql.NullString{}
				jwtHeaderName := sql.NullString{}
				jwtEndpoint := sql.NullString{}

				err := rows.Scan(
					&idp.ID,
					&idp.ResourceOwner,
					&idp.CreationDate,
					&idp.ChangeDate,
					&idp.Sequence,
					&idp.State,
					&idp.Name,
					&idp.StylingType,
					&idp.OwnerType,
					&idp.AutoRegister,
					&oidcIDPID,
					&oidcClientID,
					oidcClientSecret,
					&oidcIssuer,
					&oidcScopes,
					&oidcDisplayNameMapping,
					&oidcUsernameMapping,
					&oidcAuthorizationEndpoint,
					&oidcTokenEndpoint,
					&jwtIDPID,
					&jwtIssuer,
					&jwtKeysEndpoint,
					&jwtHeaderName,
					&jwtEndpoint,
					&count,
				)

				if err != nil {
					return nil, err
				}

				if oidcIDPID.Valid {
					idp.OIDCIDP = &OIDCIDP{
						IDPID:                 oidcIDPID.String,
						ClientID:              oidcClientID.String,
						ClientSecret:          oidcClientSecret,
						Issuer:                oidcIssuer.String,
						Scopes:                oidcScopes,
						DisplayNameMapping:    domain.OIDCMappingField(oidcDisplayNameMapping.Int32),
						UsernameMapping:       domain.OIDCMappingField(oidcUsernameMapping.Int32),
						AuthorizationEndpoint: oidcAuthorizationEndpoint.String,
						TokenEndpoint:         oidcTokenEndpoint.String,
					}
				} else if jwtIDPID.Valid {
					idp.JWTIDP = &JWTIDP{
						IDPID:        jwtIDPID.String,
						Issuer:       jwtIssuer.String,
						KeysEndpoint: jwtKeysEndpoint.String,
						HeaderName:   jwtHeaderName.String,
						Endpoint:     jwtEndpoint.String,
					}
				}

				idps = append(idps, idp)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-iiBgK", "Errors.Query.CloseRows")
			}

			return &IDPs{
				IDPs: idps,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
