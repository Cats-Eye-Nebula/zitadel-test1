package model

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type OrgIAMPolicy struct {
	es_models.ObjectRoot

	State                 int32 `json:"-"`
	UserLoginMustBeDomain bool  `json:"userLoginMustBeDomain"`
}

func OrgIAMPolicyToModel(policy *OrgIAMPolicy) *iam_model.OrgIAMPolicy {
	return &iam_model.OrgIAMPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 iam_model.PolicyState(policy.State),
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
	}
}

func (p *OrgIAMPolicy) Changes(changed *OrgIAMPolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 1)

	if p.UserLoginMustBeDomain != changed.UserLoginMustBeDomain {
		changes["userLoginMustBeDomain"] = changed.UserLoginMustBeDomain
	}
	return changes
}

func (p *OrgIAMPolicy) SetData(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "EVENT-7JS9d", "unable to unmarshal data")
	}
	return nil
}
