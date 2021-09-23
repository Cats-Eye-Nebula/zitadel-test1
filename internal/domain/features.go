package domain

import (
	"time"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	FeatureLoginPolicy              = "login_policy"
	FeatureLoginPolicyFactors       = FeatureLoginPolicy + ".factors"
	FeatureLoginPolicyIDP           = FeatureLoginPolicy + ".idp"
	FeatureLoginPolicyPasswordless  = FeatureLoginPolicy + ".passwordless"
	FeatureLoginPolicyRegistration  = FeatureLoginPolicy + ".registration"
	FeatureLoginPolicyUsernameLogin = FeatureLoginPolicy + ".username_login"
	FeatureLoginPolicyPasswordReset = FeatureLoginPolicy + ".password_reset"
	FeaturePasswordComplexityPolicy = "password_complexity_policy"
	FeatureLabelPolicy              = "label_policy"
	FeatureLabelPolicyPrivateLabel  = FeatureLabelPolicy + ".private_label"
	FeatureLabelPolicyWatermark     = FeatureLabelPolicy + ".watermark"
	FeatureCustomDomain             = "custom_domain"
	FeaturePrivacyPolicy            = "privacy_policy"
	FeatureLockoutPolicy            = "lockout_policy"
	FeatureMetadata                 = "metadata"
	FeatureCustomText               = "custom_text"
	FeatureCustomTextMessage        = FeatureCustomText + ".message"
	FeatureCustomTextLogin          = FeatureCustomText + ".login"
	FeatureMetadataUser             = FeatureMetadata + ".user"
)

type Features struct {
	es_models.ObjectRoot

	TierName         string
	TierDescription  string
	State            FeaturesState
	StateDescription string
	IsDefault        bool

	AuditLogRetention        time.Duration
	LoginPolicyFactors       bool
	LoginPolicyIDP           bool
	LoginPolicyPasswordless  bool
	LoginPolicyRegistration  bool
	LoginPolicyUsernameLogin bool
	LoginPolicyPasswordReset bool
	PasswordComplexityPolicy bool
	LabelPolicyPrivateLabel  bool
	LabelPolicyWatermark     bool
	CustomDomain             bool
	CustomTextMessage        bool
	CustomTextLogin          bool
	PrivacyPolicy            bool
	MetadataUser             bool
	LockoutPolicy            bool
}

type FeaturesState int32

const (
	FeaturesStateUnspecified FeaturesState = iota
	FeaturesStateActive
	FeaturesStateActionRequired
	FeaturesStateCanceled
	FeaturesStateGrandfathered
	FeaturesStateRemoved

	featuresStateCount
)

func (f FeaturesState) Valid() bool {
	return f >= 0 && f < featuresStateCount
}
