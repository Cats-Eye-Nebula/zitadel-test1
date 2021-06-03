package model

import (
	"time"

	"github.com/caos/zitadel/internal/domain"
)

type FeaturesView struct {
	AggregateID  string
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
	Default      bool

	TierName                 string
	TierDescription          string
	State                    domain.FeaturesState
	StateDescription         string
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
}

func (f *FeaturesView) FeatureList() []string {
	list := make([]string, 0)
	if f.LoginPolicyFactors {
		list = append(list, domain.FeatureLoginPolicyFactors)
	}
	if f.LoginPolicyIDP {
		list = append(list, domain.FeatureLoginPolicyIDP)
	}
	if f.LoginPolicyPasswordless {
		list = append(list, domain.FeatureLoginPolicyPasswordless)
	}
	if f.LoginPolicyRegistration {
		list = append(list, domain.FeatureLoginPolicyRegistration)
	}
	if f.LoginPolicyUsernameLogin {
		list = append(list, domain.FeatureLoginPolicyUsernameLogin)
	}
	if f.LoginPolicyPasswordReset {
		list = append(list, domain.FeatureLoginPolicyPasswordReset)
	}
	if f.PasswordComplexityPolicy {
		list = append(list, domain.FeaturePasswordComplexityPolicy)
	}
	if f.LabelPolicyPrivateLabel {
		list = append(list, domain.FeatureLabelPolicyPrivateLabel)
	}
	if f.LabelPolicyWatermark {
		list = append(list, domain.FeatureLabelPolicyWatermark)
	}
	if f.CustomDomain {
		list = append(list, domain.FeatureCustomDomain)
	}
	return list
}

type FeaturesSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn FeaturesSearchKey
	Asc           bool
	Queries       []*FeaturesSearchQuery
}

type FeaturesSearchKey int32

const (
	FeaturesSearchKeyUnspecified FeaturesSearchKey = iota
	FeaturesSearchKeyAggregateID
	FeaturesSearchKeyDefault
)

type FeaturesSearchQuery struct {
	Key    FeaturesSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type FeaturesSearchResult struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*FeaturesView
	Sequence    uint64
	Timestamp   time.Time
}
