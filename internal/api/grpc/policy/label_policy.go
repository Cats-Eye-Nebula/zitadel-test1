package policy

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/iam/model"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
)

func ModelLabelPolicyToPb(policy *model.LabelPolicyView) *policy_pb.LabelPolicy {
	return &policy_pb.LabelPolicy{
		IsDefault:           policy.Default,
		PrimaryColor:        policy.PrimaryColor,
		BackgroundColor:     policy.BackgroundColor,
		FontColor:           policy.FontColor,
		WarnColor:           policy.WarnColor,
		PrimaryColorDark:    policy.PrimaryColorDark,
		BackgroundColorDark: policy.BackgroundColorDark,
		WarnColorDark:       policy.WarnColorDark,
		FontColorDark:       policy.FontColorDark,
		FontUrl:             policy.FontURL,
		LogoUrl:             policy.LogoURL,
		LogoUrlDark:         policy.LogoDarkURL,
		IconUrl:             policy.IconURL,
		IconUrlDark:         policy.IconDarkURL,

		DisableWatermark:    policy.DisableWatermark,
		HideLoginNameSuffix: policy.HideLoginNameSuffix,
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			"", //TODO: resourceowner
		),
	}
}
