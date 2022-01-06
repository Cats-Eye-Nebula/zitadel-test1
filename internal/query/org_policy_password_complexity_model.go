package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

type OrgPasswordComplexityPolicyReadModel struct {
	PasswordComplexityPolicyReadModel
}

func (rm *OrgPasswordComplexityPolicyReadModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.PasswordComplexityPolicyAddedEvent:
			rm.PasswordComplexityPolicyReadModel.AppendEvents(&e.PasswordComplexityPolicyAddedEvent)
		case *org.PasswordComplexityPolicyChangedEvent:
			rm.PasswordComplexityPolicyReadModel.AppendEvents(&e.PasswordComplexityPolicyChangedEvent)
		case *policy.PasswordComplexityPolicyAddedEvent, *policy.PasswordComplexityPolicyChangedEvent:
			rm.PasswordComplexityPolicyReadModel.AppendEvents(e)
		}
	}
}
