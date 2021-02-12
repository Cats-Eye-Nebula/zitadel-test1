package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

const (
	loginPolicyTable = "management.login_policies"
)

type LoginPolicy struct {
	handler
	subscription *eventstore.Subscription
}

func newLoginPolicy(handler handler) *LoginPolicy {
	h := &LoginPolicy{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *LoginPolicy) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (m *LoginPolicy) ViewModel() string {
	return loginPolicyTable
}

func (_ *LoginPolicy) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (m *LoginPolicy) CurrentSequence() (uint64, error) {
	sequence, err := m.view.GetLatestLoginPolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *LoginPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestLoginPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *LoginPolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processLoginPolicy(event)
	}
	return err
}

func (m *LoginPolicy) processLoginPolicy(event *models.Event) (err error) {
	policy := new(iam_model.LoginPolicyView)
	switch event.Type {
	case iam_es_model.LoginPolicyAdded, model.LoginPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.LoginPolicyChanged, model.LoginPolicyChanged,
		iam_es_model.LoginPolicySecondFactorAdded, model.LoginPolicySecondFactorAdded,
		iam_es_model.LoginPolicySecondFactorRemoved, model.LoginPolicySecondFactorRemoved,
		iam_es_model.LoginPolicyMultiFactorAdded, model.LoginPolicyMultiFactorAdded,
		iam_es_model.LoginPolicyMultiFactorRemoved, model.LoginPolicyMultiFactorRemoved:
		policy, err = m.view.LoginPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case model.LoginPolicyRemoved:
		return m.view.DeleteLoginPolicy(event.AggregateID, event)
	default:
		return m.view.ProcessedLoginPolicySequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutLoginPolicy(policy, event)
}

func (m *LoginPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-4Djo9", "id", event.AggregateID).WithError(err).Warn("something went wrong in login policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestLoginPolicyFailedEvent, m.view.ProcessedLoginPolicyFailedEvent, m.view.ProcessedLoginPolicySequence, m.errorCountUntilSkip)
}

func (m *LoginPolicy) OnSuccess() error {
	return spooler.HandleSuccess(m.view.UpdateLoginPolicySpoolerRunTimestamp)
}
