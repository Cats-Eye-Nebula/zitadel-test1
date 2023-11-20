package command

import (
	"net/url"

	"github.com/zitadel/zitadel/v2/internal/crypto"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/repository/idpintent"
)

type IDPIntentWriteModel struct {
	eventstore.WriteModel

	SuccessURL  *url.URL
	FailureURL  *url.URL
	IDPID       string
	IDPUser     []byte
	IDPUserID   string
	IDPUserName string
	UserID      string

	IDPAccessToken *crypto.CryptoValue
	IDPIDToken     string

	IDPEntryAttributes map[string][]string

	RequestID string
	Assertion *crypto.CryptoValue

	State     domain.IDPIntentState
	aggregate *eventstore.Aggregate
}

func NewIDPIntentWriteModel(id, resourceOwner string) *IDPIntentWriteModel {
	return &IDPIntentWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
		},
		aggregate: &idpintent.NewAggregate(id, resourceOwner).Aggregate,
	}
}

func (wm *IDPIntentWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idpintent.StartedEvent:
			wm.reduceStartedEvent(e)
		case *idpintent.SucceededEvent:
			wm.reduceOAuthSucceededEvent(e)
		case *idpintent.SAMLSucceededEvent:
			wm.reduceSAMLSucceededEvent(e)
		case *idpintent.SAMLRequestEvent:
			wm.reduceSAMLRequestEvent(e)
		case *idpintent.LDAPSucceededEvent:
			wm.reduceLDAPSucceededEvent(e)
		case *idpintent.FailedEvent:
			wm.reduceFailedEvent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *IDPIntentWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(idpintent.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			idpintent.StartedEventType,
			idpintent.SucceededEventType,
			idpintent.SAMLSucceededEventType,
			idpintent.SAMLRequestEventType,
			idpintent.LDAPSucceededEventType,
			idpintent.FailedEventType,
		).
		Builder()
}

func (wm *IDPIntentWriteModel) reduceStartedEvent(e *idpintent.StartedEvent) {
	wm.SuccessURL = e.SuccessURL
	wm.FailureURL = e.FailureURL
	wm.IDPID = e.IDPID
	wm.State = domain.IDPIntentStateStarted
}

func (wm *IDPIntentWriteModel) reduceSAMLSucceededEvent(e *idpintent.SAMLSucceededEvent) {
	wm.UserID = e.UserID
	wm.IDPUser = e.IDPUser
	wm.IDPUserID = e.IDPUserID
	wm.IDPUserName = e.IDPUserName
	wm.Assertion = e.Assertion
	wm.State = domain.IDPIntentStateSucceeded
}

func (wm *IDPIntentWriteModel) reduceLDAPSucceededEvent(e *idpintent.LDAPSucceededEvent) {
	wm.UserID = e.UserID
	wm.IDPUser = e.IDPUser
	wm.IDPUserID = e.IDPUserID
	wm.IDPUserName = e.IDPUserName
	wm.IDPEntryAttributes = e.EntryAttributes
	wm.State = domain.IDPIntentStateSucceeded
}

func (wm *IDPIntentWriteModel) reduceOAuthSucceededEvent(e *idpintent.SucceededEvent) {
	wm.UserID = e.UserID
	wm.IDPUser = e.IDPUser
	wm.IDPUserID = e.IDPUserID
	wm.IDPUserName = e.IDPUserName
	wm.IDPAccessToken = e.IDPAccessToken
	wm.IDPIDToken = e.IDPIDToken
	wm.State = domain.IDPIntentStateSucceeded
}

func (wm *IDPIntentWriteModel) reduceSAMLRequestEvent(e *idpintent.SAMLRequestEvent) {
	wm.RequestID = e.RequestID
}

func (wm *IDPIntentWriteModel) reduceFailedEvent(e *idpintent.FailedEvent) {
	wm.State = domain.IDPIntentStateFailed
}
