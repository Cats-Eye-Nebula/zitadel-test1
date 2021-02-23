package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/v1"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

const (
	passwordAgePolicyTable = "adminapi.password_age_policies"
)

type PasswordAgePolicy struct {
	handler
	subscription *v1.Subscription
}

func newPasswordAgePolicy(handler handler) *PasswordAgePolicy {
	h := &PasswordAgePolicy{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (p *PasswordAgePolicy) subscribe() {
	p.subscription = p.es.Subscribe(p.AggregateTypes()...)
	go func() {
		for event := range p.subscription.Events {
			query.ReduceEvent(p, event)
		}
	}()
}

func (p *PasswordAgePolicy) ViewModel() string {
	return passwordAgePolicyTable
}

func (p *PasswordAgePolicy) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *PasswordAgePolicy) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestPasswordAgePolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *PasswordAgePolicy) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := p.view.GetLatestPasswordAgePolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *PasswordAgePolicy) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = p.processPasswordAgePolicy(event)
	}
	return err
}

func (p *PasswordAgePolicy) processPasswordAgePolicy(event *es_models.Event) (err error) {
	policy := new(iam_model.PasswordAgePolicyView)
	switch event.Type {
	case iam_es_model.PasswordAgePolicyAdded, model.PasswordAgePolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.PasswordAgePolicyChanged, model.PasswordAgePolicyChanged:
		policy, err = p.view.PasswordAgePolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case model.PasswordAgePolicyRemoved:
		return p.view.DeletePasswordAgePolicy(event.AggregateID, event)
	default:
		return p.view.ProcessedPasswordAgePolicySequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutPasswordAgePolicy(policy, event)
}

func (p *PasswordAgePolicy) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-nD8sie", "id", event.AggregateID).WithError(err).Warn("something went wrong in passwordAge policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestPasswordAgePolicyFailedEvent, p.view.ProcessedPasswordAgePolicyFailedEvent, p.view.ProcessedPasswordAgePolicySequence, p.errorCountUntilSkip)
}

func (p *PasswordAgePolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateProcessedPasswordAgePolicySpoolerRunTimestamp)
}
