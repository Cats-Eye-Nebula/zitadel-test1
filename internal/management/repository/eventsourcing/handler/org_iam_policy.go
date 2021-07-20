package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

const (
	orgIAMPolicyTable = "management.org_iam_policies"
)

type OrgIAMPolicy struct {
	handler
	subscription *v1.Subscription
}

func newOrgIAMPolicy(handler handler) *OrgIAMPolicy {
	h := &OrgIAMPolicy{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *OrgIAMPolicy) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (m *OrgIAMPolicy) ViewModel() string {
	return orgIAMPolicyTable
}

func (m *OrgIAMPolicy) Subscription() *v1.Subscription {
	return m.subscription
}

func (_ *OrgIAMPolicy) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *OrgIAMPolicy) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestOrgIAMPolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *OrgIAMPolicy) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := m.view.GetLatestOrgIAMPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *OrgIAMPolicy) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processOrgIAMPolicy(event)
	}
	return err
}

func (m *OrgIAMPolicy) processOrgIAMPolicy(event *es_models.Event) (err error) {
	policy := new(iam_model.OrgIAMPolicyView)
	switch event.Type {
	case iam_es_model.OrgIAMPolicyAdded, model.OrgIAMPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.OrgIAMPolicyChanged, model.OrgIAMPolicyChanged:
		policy, err = m.view.OrgIAMPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case model.OrgIAMPolicyRemoved:
		return m.view.DeleteOrgIAMPolicy(event.AggregateID, event)
	default:
		return m.view.ProcessedOrgIAMPolicySequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutOrgIAMPolicy(policy, event)
}

func (m *OrgIAMPolicy) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-3Gf9s", "id", event.AggregateID).WithError(err).Warn("something went wrong in orgIAM policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestOrgIAMPolicyFailedEvent, m.view.ProcessedOrgIAMPolicyFailedEvent, m.view.ProcessedOrgIAMPolicySequence, m.errorCountUntilSkip)
}

func (o *OrgIAMPolicy) OnSuccess() error {
	return spooler.HandleSuccess(o.view.UpdateOrgIAMPolicySpoolerRunTimestamp)
}
