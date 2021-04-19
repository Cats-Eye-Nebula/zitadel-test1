package handler

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/view"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	userTable = "adminapi.users"
)

type User struct {
	handler
	systemDefaults systemdefaults.SystemDefaults
	subscription   *v1.Subscription
}

func newUser(
	handler handler,
	systemDefaults systemdefaults.SystemDefaults,
) *User {
	h := &User{
		handler:        handler,
		systemDefaults: systemDefaults,
	}

	h.subscribe()

	return h
}

func (u *User) subscribe() {
	u.subscription = u.es.Subscribe(u.AggregateTypes()...)
	go func() {
		for event := range u.subscription.Events {
			query.ReduceEvent(u, event)
		}
	}()
}

func (u *User) ViewModel() string {
	return userTable
}

func (u *User) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{es_model.UserAggregate, org_es_model.OrgAggregate}
}

func (u *User) CurrentSequence() (uint64, error) {
	sequence, err := u.view.GetLatestUserSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (u *User) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := u.view.GetLatestUserSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(u.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (u *User) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case es_model.UserAggregate:
		return u.ProcessUser(event)
	case org_es_model.OrgAggregate:
		return u.ProcessOrg(event)
	default:
		return nil
	}
}

func (u *User) ProcessUser(event *es_models.Event) (err error) {
	user := new(view_model.UserView)
	switch event.Type {
	case es_model.UserAdded,
		es_model.UserRegistered,
		es_model.HumanRegistered,
		es_model.MachineAdded,
		es_model.HumanAdded:
		err = user.AppendEvent(event)
		if err != nil {
			return err
		}
		err = u.fillLoginNames(user)
	case es_model.UserProfileChanged,
		es_model.UserEmailChanged,
		es_model.UserEmailVerified,
		es_model.UserPhoneChanged,
		es_model.UserPhoneVerified,
		es_model.UserPhoneRemoved,
		es_model.UserAddressChanged,
		es_model.UserDeactivated,
		es_model.UserReactivated,
		es_model.UserLocked,
		es_model.UserUnlocked,
		es_model.MFAOTPAdded,
		es_model.MFAOTPVerified,
		es_model.MFAOTPRemoved,
		es_model.HumanProfileChanged,
		es_model.HumanEmailChanged,
		es_model.HumanEmailVerified,
		es_model.HumanPhoneChanged,
		es_model.HumanPhoneVerified,
		es_model.HumanPhoneRemoved,
		es_model.HumanAddressChanged,
		es_model.HumanMFAOTPAdded,
		es_model.HumanMFAOTPVerified,
		es_model.HumanMFAOTPRemoved,
		es_model.HumanMFAU2FTokenAdded,
		es_model.HumanMFAU2FTokenVerified,
		es_model.HumanMFAU2FTokenRemoved,
		es_model.HumanPasswordlessTokenAdded,
		es_model.HumanPasswordlessTokenVerified,
		es_model.HumanPasswordlessTokenRemoved,
		es_model.MachineChanged:
		user, err = u.view.UserByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = user.AppendEvent(event)
	case es_model.DomainClaimed,
		es_model.UserUserNameChanged:
		user, err = u.view.UserByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = user.AppendEvent(event)
		if err != nil {
			return err
		}
		err = u.fillLoginNames(user)
	case es_model.UserRemoved:
		return u.view.DeleteUser(event.AggregateID, event)
	default:
		return u.view.ProcessedUserSequence(event)
	}
	if err != nil {
		return err
	}
	return u.view.PutUser(user, event)
}

func (u *User) ProcessOrg(event *es_models.Event) (err error) {
	switch event.Type {
	case org_es_model.OrgDomainVerified,
		org_es_model.OrgDomainRemoved,
		org_es_model.OrgIAMPolicyAdded,
		org_es_model.OrgIAMPolicyChanged,
		org_es_model.OrgIAMPolicyRemoved:
		return u.fillLoginNamesOnOrgUsers(event)
	case org_es_model.OrgDomainPrimarySet:
		return u.fillPreferredLoginNamesOnOrgUsers(event)
	default:
		return u.view.ProcessedUserSequence(event)
	}
}

func (u *User) fillLoginNamesOnOrgUsers(event *es_models.Event) error {
	org, err := u.getOrgByID(context.Background(), event.ResourceOwner)
	if err != nil {
		return err
	}
	policy := org.OrgIamPolicy
	if policy == nil {
		policy, err = u.getDefaultOrgIAMPolicy(context.Background())
		if err != nil {
			return err
		}
	}
	users, err := u.view.UsersByOrgID(event.AggregateID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.SetLoginNames(policy, org.Domains)
	}
	return u.view.PutUsers(users, event)
}

func (u *User) fillPreferredLoginNamesOnOrgUsers(event *es_models.Event) error {
	org, err := u.getOrgByID(context.Background(), event.ResourceOwner)
	if err != nil {
		return err
	}
	policy := org.OrgIamPolicy
	if policy == nil {
		policy, err = u.getDefaultOrgIAMPolicy(context.Background())
		if err != nil {
			return err
		}
	}
	if !policy.UserLoginMustBeDomain {
		return nil
	}
	users, err := u.view.UsersByOrgID(event.AggregateID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.PreferredLoginName = user.GenerateLoginName(org.GetPrimaryDomain().Domain, policy.UserLoginMustBeDomain)
	}
	return u.view.PutUsers(users, event)
}

func (u *User) fillLoginNames(user *view_model.UserView) (err error) {
	org, err := u.getOrgByID(context.Background(), user.ResourceOwner)
	if err != nil {
		return err
	}
	policy := org.OrgIamPolicy
	if policy == nil {
		policy, err = u.getDefaultOrgIAMPolicy(context.Background())
		if err != nil {
			return err
		}
	}

	user.SetLoginNames(policy, org.Domains)
	user.PreferredLoginName = user.GenerateLoginName(org.GetPrimaryDomain().Domain, policy.UserLoginMustBeDomain)
	return nil
}

func (u *User) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-vLmwQ", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserFailedEvent, u.view.ProcessedUserFailedEvent, u.view.ProcessedUserSequence, u.errorCountUntilSkip)
}

func (u *User) OnSuccess() error {
	return spooler.HandleSuccess(u.view.UpdateUserSpoolerRunTimestamp)
}

func (u *User) getOrgByID(ctx context.Context, orgID string) (*org_model.Org, error) {
	query, err := view.OrgByIDQuery(orgID, 0)
	if err != nil {
		return nil, err
	}

	esOrg := &org_es_model.Org{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: orgID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, esOrg.AppendEvents, query)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-kVLb2", "Errors.Org.NotFound")
	}

	return org_es_model.OrgToModel(esOrg), nil
}

func (u *User) getIAMByID(ctx context.Context) (*iam_model.IAM, error) {
	query, err := iam_view.IAMByIDQuery(domain.IAMID, 0)
	if err != nil {
		return nil, err
	}
	iam := &model.IAM{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: domain.IAMID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, iam.AppendEvents, query)
	if err != nil && caos_errs.IsNotFound(err) && iam.Sequence == 0 {
		return nil, err
	}
	return model.IAMToModel(iam), nil
}

func (u *User) getDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicy, error) {
	existingIAM, err := u.getIAMByID(ctx)
	if err != nil {
		return nil, err
	}
	if existingIAM.DefaultOrgIAMPolicy == nil {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-2Fj8s", "Errors.IAM.OrgIAMPolicy.NotExisting")
	}
	return existingIAM.DefaultOrgIAMPolicy, nil
}
