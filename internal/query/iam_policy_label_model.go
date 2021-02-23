package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/policy"
)

type IAMLabelPolicyReadModel struct{ LabelPolicyReadModel }

func (rm *IAMLabelPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.LabelPolicyAddedEvent:
			rm.LabelPolicyReadModel.AppendEvents(&e.LabelPolicyAddedEvent)
		case *iam.LabelPolicyChangedEvent:
			rm.LabelPolicyReadModel.AppendEvents(&e.LabelPolicyChangedEvent)
		case *policy.LabelPolicyAddedEvent, *policy.LabelPolicyChangedEvent:
			rm.LabelPolicyReadModel.AppendEvents(e)
		}
	}
}
