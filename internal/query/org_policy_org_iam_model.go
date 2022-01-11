package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

type OrgOrgIAMPolicyReadModel struct{ OrgIAMPolicyReadModel }

func (rm *OrgOrgIAMPolicyReadModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.OrgIAMPolicyAddedEvent:
			rm.OrgIAMPolicyReadModel.AppendEvents(&e.OrgIAMPolicyAddedEvent)
		case *policy.OrgIAMPolicyAddedEvent:
			rm.OrgIAMPolicyReadModel.AppendEvents(e)
		}
	}
}
