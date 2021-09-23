package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/v1"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

const (
	lockoutPolicyTable = "auth.lockout_policies"
)

type LockoutPolicy struct {
	handler
	subscription *v1.Subscription
}

func newLockoutPolicy(handler handler) *LockoutPolicy {
	h := &LockoutPolicy{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (p *LockoutPolicy) subscribe() {
	p.subscription = p.es.Subscribe(p.AggregateTypes()...)
	go func() {
		for event := range p.subscription.Events {
			query.ReduceEvent(p, event)
		}
	}()
}

func (p *LockoutPolicy) ViewModel() string {
	return lockoutPolicyTable
}

func (p *LockoutPolicy) Subscription() *v1.Subscription {
	return p.subscription
}

func (_ *LockoutPolicy) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{org_es_model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *LockoutPolicy) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestLockoutPolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *LockoutPolicy) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := p.view.GetLatestLockoutPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *LockoutPolicy) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case org_es_model.OrgAggregate, iam_es_model.IAMAggregate:
		err = p.processLockoutPolicy(event)
	}
	return err
}

func (p *LockoutPolicy) processLockoutPolicy(event *es_models.Event) (err error) {
	policy := new(iam_model.LockoutPolicyView)
	switch event.Type {
	case iam_es_model.LockoutPolicyAdded, org_es_model.LockoutPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.LockoutPolicyChanged, org_es_model.LockoutPolicyChanged:
		policy, err = p.view.LockoutPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case org_es_model.LockoutPolicyRemoved:
		return p.view.DeleteLockoutPolicy(event.AggregateID, event)
	default:
		return p.view.ProcessedLockoutPolicySequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutLockoutPolicy(policy, event)
}

func (p *LockoutPolicy) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-0pos2", "id", event.AggregateID).WithError(err).Warn("something went wrong in passwordLockout policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestLockoutPolicyFailedEvent, p.view.ProcessedLockoutPolicyFailedEvent, p.view.ProcessedLockoutPolicySequence, p.errorCountUntilSkip)
}

func (p *LockoutPolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateLockoutPolicySpoolerRunTimestamp)
}
