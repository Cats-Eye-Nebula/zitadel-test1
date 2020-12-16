package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

const (
	orgIAMPolicyTable = "adminapi.org_iam_policies"
)

type OrgIAMPolicy struct {
	handler
	subscription *eventstore.Subscription
}

func newOrgIAMPolicy(handler handler) *OrgIAMPolicy {
	h := &OrgIAMPolicy{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (p *OrgIAMPolicy) subscribe() {
	p.subscription = p.es.Subscribe(p.AggregateTypes()...)
	go func() {
		for event := range p.subscription.Events {
			query.ReduceEvent(p, event)
		}
	}()
}

func (p *OrgIAMPolicy) ViewModel() string {
	return orgIAMPolicyTable
}

func (p *OrgIAMPolicy) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *OrgIAMPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestOrgIAMPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *OrgIAMPolicy) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestOrgIAMPolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *OrgIAMPolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = p.processOrgIAMPolicy(event)
	}
	return err
}

func (p *OrgIAMPolicy) processOrgIAMPolicy(event *models.Event) (err error) {
	policy := new(iam_model.OrgIAMPolicyView)
	switch event.Type {
	case iam_es_model.OrgIAMPolicyAdded, model.OrgIAMPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.OrgIAMPolicyChanged, model.OrgIAMPolicyChanged:
		policy, err = p.view.OrgIAMPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case model.OrgIAMPolicyRemoved:
		return p.view.DeleteOrgIAMPolicy(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return p.view.ProcessedOrgIAMPolicySequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return p.view.PutOrgIAMPolicy(policy, policy.Sequence, event.CreationDate)
}

func (p *OrgIAMPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Wm8fs", "id", event.AggregateID).WithError(err).Warn("something went wrong in orgIAM policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestOrgIAMPolicyFailedEvent, p.view.ProcessedOrgIAMPolicyFailedEvent, p.view.ProcessedOrgIAMPolicySequence, p.errorCountUntilSkip)
}

func (p *OrgIAMPolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateOrgIAMPolicySpoolerRunTimestamp)
}
