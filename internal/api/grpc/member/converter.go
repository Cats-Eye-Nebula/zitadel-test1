package member

import (
	"github.com/caos/zitadel/internal/v2/domain"
	member_pb "github.com/caos/zitadel/pkg/grpc/member"
)

func MemberToDomain(member *member_pb.Member) *domain.Member {
	return &domain.Member{
		UserID: member.UserId,
		Roles:  member.Roles,
	}
}
