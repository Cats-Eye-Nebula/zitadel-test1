package policy

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/iam/model"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
)

func ModelLabelPolicyToPb(policy *model.LabelPolicyView) *policy_pb.LabelPolicy {
	return &policy_pb.LabelPolicy{
		PrimaryColor:   policy.PrimaryColor,
		SecondaryColor: policy.SecondaryColor,
		Details: object.ToDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			"policy.ResourceOwner", //TODO: für da haui öppert
		),
	}
}
