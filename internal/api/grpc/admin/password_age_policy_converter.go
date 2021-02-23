package admin

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func passwordAgePolicyToDomain(policy *admin.DefaultPasswordAgePolicyRequest) *domain.PasswordAgePolicy {
	return &domain.PasswordAgePolicy{
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
	}
}

func passwordAgePolicyFromDomain(policy *domain.PasswordAgePolicy) *admin.DefaultPasswordAgePolicy {
	return &admin.DefaultPasswordAgePolicy{
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
		ChangeDate:     timestamppb.New(policy.ChangeDate),
	}
}

func passwordAgePolicyViewFromModel(policy *iam_model.PasswordAgePolicyView) *admin.DefaultPasswordAgePolicyView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-2Gs9o").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-8Hjss").OnError(err).Debug("date parse failed")

	return &admin.DefaultPasswordAgePolicyView{
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
		CreationDate:   creationDate,
		ChangeDate:     changeDate,
	}
}
