package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_model "github.com/caos/zitadel/internal/org/repository/view/model"
)

const (
	orgTable = "management.orgs"
)

type Org struct {
	handler
	subscription *eventstore.Subscription
}

func newOrg(handler handler) *Org {
	h := &Org{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *Org) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (o *Org) ViewModel() string {
	return orgTable
}

func (_ *Org) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate}
}

func (o *Org) CurrentSequence() (uint64, error) {
	sequence, err := o.view.GetLatestOrgSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (o *Org) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := o.view.GetLatestOrgSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.OrgQuery(sequence.CurrentSequence), nil
}

func (o *Org) Reduce(event *es_models.Event) (err error) {
	org := new(org_model.OrgView)

	switch event.Type {
	case model.OrgAdded:
		err = org.AppendEvent(event)
	case model.OrgChanged:
		err = org.SetData(event)
		if err != nil {
			return err
		}
		org, err = o.view.OrgByID(org.ID)
		if err != nil {
			return err
		}
		err = org.AppendEvent(event)
	default:
		return o.view.ProcessedOrgSequence(event)
	}
	if err != nil {
		return err
	}
	return o.view.PutOrg(org, event)
}

func (o *Org) OnError(event *es_models.Event, spoolerErr error) error {
	logging.LogWithFields("SPOOL-ls9ew", "id", event.AggregateID).WithError(spoolerErr).Warn("something went wrong in project app handler")
	return spooler.HandleError(event, spoolerErr, o.view.GetLatestOrgFailedEvent, o.view.ProcessedOrgFailedEvent, o.view.ProcessedOrgSequence, o.errorCountUntilSkip)
}

func (o *Org) OnSuccess() error {
	return spooler.HandleSuccess(o.view.UpdateOrgSpoolerRunTimestamp)
}
