package member

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	proj_model "github.com/caos/zitadel/internal/project/model"
	member_pb "github.com/caos/zitadel/pkg/grpc/member"
)

func ProjectGrantMembersToPb(members []*proj_model.ProjectGrantMemberView) []*member_pb.Member {
	m := make([]*member_pb.Member, len(members))
	for i, member := range members {
		m[i] = ProjectGrantMemberToPb(member)
	}
	return m
}

func ProjectGrantMemberToPb(m *proj_model.ProjectGrantMemberView) *member_pb.Member {
	return &member_pb.Member{
		UserId: m.UserID,
		Roles:  m.Roles,
		// PreferredLoginName: //TODO: not implemented in be
		Email:       m.Email,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		DisplayName: m.DisplayName,
		Details: object.ToDetailsPb(
			m.Sequence,
			m.ChangeDate,
			"m.ResourceOwner", //TODO: not returnd
		),
	}
}

func MemberQueriesToProjectGrantMember(queries []*member_pb.SearchQuery) []*proj_model.ProjectGrantMemberSearchQuery {
	q := make([]*proj_model.ProjectGrantMemberSearchQuery, len(queries))
	for i, query := range queries {
		q[i] = MemberQueryToProjectGrantMember(query)
	}
	return q
}

func MemberQueryToProjectGrantMember(query *member_pb.SearchQuery) *proj_model.ProjectGrantMemberSearchQuery {
	switch q := query.Query.(type) {
	case *member_pb.SearchQuery_Email:
		return EmailQueryToProjectGrantMemberQuery(q.Email)
	case *member_pb.SearchQuery_FirstName:
		return FirstNameQueryToProjectGrantMemberQuery(q.FirstName)
	case *member_pb.SearchQuery_LastName:
		return LastNameQueryToProjectGrantMemberQuery(q.LastName)
	case *member_pb.SearchQuery_UserId:
		return UserIDQueryToProjectGrantMemberQuery(q.UserId)
	default:
		return nil
	}
}

func FirstNameQueryToProjectGrantMemberQuery(query *member_pb.FirstNameQuery) *proj_model.ProjectGrantMemberSearchQuery {
	return &proj_model.ProjectGrantMemberSearchQuery{
		Key:    proj_model.ProjectGrantMemberSearchKeyFirstName,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.FirstName,
	}
}

func LastNameQueryToProjectGrantMemberQuery(query *member_pb.LastNameQuery) *proj_model.ProjectGrantMemberSearchQuery {
	return &proj_model.ProjectGrantMemberSearchQuery{
		Key:    proj_model.ProjectGrantMemberSearchKeyLastName,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.LastName,
	}
}

func EmailQueryToProjectGrantMemberQuery(query *member_pb.EmailQuery) *proj_model.ProjectGrantMemberSearchQuery {
	return &proj_model.ProjectGrantMemberSearchQuery{
		Key:    proj_model.ProjectGrantMemberSearchKeyEmail,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.Email,
	}
}

func UserIDQueryToProjectGrantMemberQuery(query *member_pb.UserIDQuery) *proj_model.ProjectGrantMemberSearchQuery {
	return &proj_model.ProjectGrantMemberSearchQuery{
		Key:    proj_model.ProjectGrantMemberSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  query.UserId,
	}
}
