package query

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/query/projection"
	"github.com/caos/zitadel/internal/telemetry/tracing"

	sq "github.com/Masterminds/squirrel"
)

type MembersQuery struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *MembersQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func NewMemberEmailSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(HumanEmailCol, value, method)
}

func NewMemberFirstNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(HumanFirstNameCol, value, method)
}

func NewMemberLastNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(HumanLastNameCol, value, method)
}

func NewMemberUserIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(memberUserID, value, TextEquals)
}
func NewMemberResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(memberResourceOwner, value, TextEquals)
}

type Members struct {
	SearchResponse
	Members []*Member
}

type Member struct {
	CreationDate  time.Time
	ChangeDate    time.Time
	Sequence      uint64
	ResourceOwner string

	UserID             string
	Roles              []string
	PreferredLoginName string
	Email              string
	FirstName          string
	LastName           string
	DisplayName        string
	AvatarURL          string
}

func (r *Queries) IAMMemberByID(ctx context.Context, iamID, userID string) (member *IAMMemberReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	member = NewIAMMemberReadModel(iamID, userID)
	err = r.eventstore.FilterToQueryReducer(ctx, member)
	if err != nil {
		return nil, err
	}

	return member, nil
}

var (
	memberTableAlias = table{
		name:  "members",
		alias: "members",
	}
	memberUserID = Column{
		name:  projection.MemberUserIDCol,
		table: memberTableAlias,
	}
	memberResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: memberTableAlias,
	}
)
