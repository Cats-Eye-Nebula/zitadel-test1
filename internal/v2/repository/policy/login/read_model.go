package login

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/idp/provider"
)

const (
	LoginPolicyAddedEventType              = "policy.login.added"
	LoginPolicyChangedEventType            = "policy.login.changed"
	LoginPolicyRemovedEventType            = "policy.login.removed"
	LoginPolicyIDPProviderAddedEventType   = "policy.login." + provider.AddedEventType
	LoginPolicyIDPProviderRemovedEventType = "policy.login." + provider.RemovedEventType
)

type LoginPolicyReadModel struct {
	eventstore.ReadModel

	AllowUserNamePassword bool
	AllowRegister         bool
	AllowExternalIDP      bool
	ForceMFA              bool
	PasswordlessType      PasswordlessType
}

func (rm *LoginPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *LoginPolicyAddedEvent:
			rm.AllowUserNamePassword = e.AllowUserNamePassword
			rm.AllowExternalIDP = e.AllowExternalIDP
			rm.AllowRegister = e.AllowRegister
			rm.ForceMFA = e.ForceMFA
			rm.PasswordlessType = e.PasswordlessType
		case *LoginPolicyChangedEvent:
			rm.AllowUserNamePassword = e.AllowUserNamePassword
			rm.AllowExternalIDP = e.AllowExternalIDP
			rm.AllowRegister = e.AllowRegister
			rm.ForceMFA = e.ForceMFA
			rm.PasswordlessType = e.PasswordlessType
		}
	}
	return rm.ReadModel.Reduce()
}
