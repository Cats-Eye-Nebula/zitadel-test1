package admin

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func passwordLockoutPolicyToDomain(policy *admin.DefaultPasswordLockoutPolicyRequest) *domain.PasswordLockoutPolicy {
	return &domain.PasswordLockoutPolicy{
		MaxAttempts:         policy.MaxAttempts,
		ShowLockOutFailures: policy.ShowLockoutFailure,
	}
}

func passwordLockoutPolicyFromDomain(policy *domain.PasswordLockoutPolicy) *admin.DefaultPasswordLockoutPolicy {
	return &admin.DefaultPasswordLockoutPolicy{
		MaxAttempts:        policy.MaxAttempts,
		ShowLockoutFailure: policy.ShowLockOutFailures,
		ChangeDate:         timestamppb.New(policy.ChangeDate),
	}
}

func passwordLockoutPolicyViewFromModel(policy *iam_model.PasswordLockoutPolicyView) *admin.DefaultPasswordLockoutPolicyView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-7Hmlo").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-0oLgs").OnError(err).Debug("date parse failed")

	return &admin.DefaultPasswordLockoutPolicyView{
		MaxAttempts:        policy.MaxAttempts,
		ShowLockoutFailure: policy.ShowLockOutFailures,
		CreationDate:       creationDate,
		ChangeDate:         changeDate,
	}
}
