package org

import (
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LabelPolicyAddedEventType   = orgEventTypePrefix + policy.LabelPolicyAddedEventType
	LabelPolicyChangedEventType = orgEventTypePrefix + policy.LabelPolicyChangedEventType
	LabelPolicyRemovedEventType = orgEventTypePrefix + policy.LabelPolicyRemovedEventType
)

type LabelPolicyAddedEvent struct {
	policy.LabelPolicyAddedEvent
}

type LabelPolicyChangedEvent struct {
	policy.LabelPolicyChangedEvent
}
