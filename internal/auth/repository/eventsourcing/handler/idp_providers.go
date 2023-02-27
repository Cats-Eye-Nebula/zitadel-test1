package handler

//
//import (
//	"context"
//
//	"github.com/zitadel/logging"
//
//	"github.com/zitadel/zitadel/internal/config/systemdefaults"
//	"github.com/zitadel/zitadel/internal/domain"
//	"github.com/zitadel/zitadel/internal/eventstore"
//	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
//	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
//	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
//	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
//	"github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
//	iam_model "github.com/zitadel/zitadel/internal/iam/model"
//	iam_view_model "github.com/zitadel/zitadel/internal/iam/repository/view/model"
//	query2 "github.com/zitadel/zitadel/internal/query"
//	"github.com/zitadel/zitadel/internal/repository/instance"
//	"github.com/zitadel/zitadel/internal/repository/org"
//)
//
//const (
//	idpProviderTable = "auth.idp_providers2"
//)
//
//type IDPProvider struct {
//	handler
//	systemDefaults systemdefaults.SystemDefaults
//	subscription   *v1.Subscription
//	queries        *query2.Queries
//}
//
//func newIDPProvider(
//	ctx context.Context,
//	h handler,
//	defaults systemdefaults.SystemDefaults,
//	queries *query2.Queries,
//) *IDPProvider {
//	idpProvider := &IDPProvider{
//		handler:        h,
//		systemDefaults: defaults,
//		queries:        queries,
//	}
//
//	idpProvider.subscribe(ctx)
//
//	return idpProvider
//}
//
//func (i *IDPProvider) subscribe(ctx context.Context) {
//	i.subscription = i.es.Subscribe(i.AggregateTypes()...)
//	go func() {
//		for event := range i.subscription.Events {
//			query.ReduceEvent(ctx, i, event)
//		}
//	}()
//}
//
//func (i *IDPProvider) ViewModel() string {
//	return idpProviderTable
//}
//
//func (i *IDPProvider) Subscription() *v1.Subscription {
//	return i.subscription
//}
//
//func (_ *IDPProvider) AggregateTypes() []models.AggregateType {
//	return []es_models.AggregateType{instance.AggregateType, org.AggregateType}
//}
//
//func (i *IDPProvider) CurrentSequence(instanceID string) (uint64, error) {
//	sequence, err := i.view.GetLatestIDPProviderSequence(instanceID)
//	if err != nil {
//		return 0, err
//	}
//	return sequence.CurrentSequence, nil
//}
//
//func (i *IDPProvider) EventQuery(instanceIDs []string) (*es_models.SearchQuery, error) {
//	sequences, err := i.view.GetLatestIDPProviderSequences(instanceIDs)
//	if err != nil {
//		return nil, err
//	}
//
//	return newSearchQuery(sequences, i.AggregateTypes(), instanceIDs), nil
//}
//
//func (i *IDPProvider) Reduce(event *models.Event) (err error) {
//	switch event.AggregateType {
//	case instance.AggregateType, org.AggregateType:
//		err = i.processIdpProvider(event)
//	}
//	return err
//}
//
//func (i *IDPProvider) processIdpProvider(event *models.Event) (err error) {
//	provider := new(iam_view_model.IDPProviderView)
//	switch eventstore.EventType(event.Type) {
//	case instance.LoginPolicyIDPProviderAddedEventType, org.LoginPolicyIDPProviderAddedEventType:
//		err = provider.AppendEvent(event)
//		if err != nil {
//			return err
//		}
//		err = i.fillData(provider)
//	case instance.LoginPolicyIDPProviderRemovedEventType, instance.LoginPolicyIDPProviderCascadeRemovedEventType,
//		org.LoginPolicyIDPProviderRemovedEventType, org.LoginPolicyIDPProviderCascadeRemovedEventType:
//		err = provider.SetData(event)
//		if err != nil {
//			return err
//		}
//		return i.view.DeleteIDPProvider(event.AggregateID, provider.IDPConfigID, event.InstanceID, event)
//	case instance.IDPConfigChangedEventType, org.IDPConfigChangedEventType:
//		esConfig := new(iam_view_model.IDPConfigView)
//		providerType := iam_model.IDPProviderTypeSystem
//		if event.AggregateID != event.InstanceID {
//			providerType = iam_model.IDPProviderTypeOrg
//		}
//		err = esConfig.AppendEvent(providerType, event)
//		if err != nil {
//			return err
//		}
//		providers, err := i.view.IDPProvidersByIDPConfigID(esConfig.IDPConfigID, event.InstanceID)
//		if err != nil {
//			return err
//		}
//		config := new(query2.IDP)
//		if event.AggregateID == event.InstanceID {
//			config, err = i.getDefaultIDPConfig(event.InstanceID, esConfig.IDPConfigID)
//		} else {
//			config, err = i.getOrgIDPConfig(event.InstanceID, event.AggregateID, esConfig.IDPConfigID)
//		}
//		if err != nil {
//			return err
//		}
//		for _, provider := range providers {
//			i.fillConfigData(provider, config)
//		}
//		return i.view.PutIDPProviders(event, providers...)
//	case org.LoginPolicyRemovedEventType:
//		return i.view.DeleteIDPProvidersByAggregateID(event.AggregateID, event.InstanceID, event)
//	case instance.InstanceRemovedEventType:
//		return i.view.DeleteInstanceIDPProviders(event)
//	case org.OrgRemovedEventType:
//		return i.view.UpdateOrgOwnerRemovedIDPProviders(event)
//	default:
//		return i.view.ProcessedIDPProviderSequence(event)
//	}
//	if err != nil {
//		return err
//	}
//	return i.view.PutIDPProvider(provider, event)
//}
//
//func (i *IDPProvider) fillData(provider *iam_view_model.IDPProviderView) (err error) {
//	var config *query2.IDP
//	if provider.IDPProviderType == int32(iam_model.IDPProviderTypeSystem) {
//		config, err = i.getDefaultIDPConfig(provider.InstanceID, provider.IDPConfigID)
//	} else {
//		config, err = i.getOrgIDPConfig(provider.InstanceID, provider.AggregateID, provider.IDPConfigID)
//	}
//	if err != nil {
//		return err
//	}
//	i.fillConfigData(provider, config)
//	return nil
//}
//
//func (i *IDPProvider) fillConfigData(provider *iam_view_model.IDPProviderView, config *query2.IDP) {
//	provider.Name = config.Name
//	provider.StylingType = int32(config.StylingType)
//	if config.OIDCIDP != nil {
//		provider.IDPConfigType = int32(domain.IDPConfigTypeOIDC)
//	} else if config.JWTIDP != nil {
//		provider.IDPConfigType = int32(domain.IDPConfigTypeJWT)
//	}
//	switch config.State {
//	case domain.IDPConfigStateActive:
//		provider.IDPState = int32(iam_model.IDPConfigStateActive)
//	case domain.IDPConfigStateInactive:
//		provider.IDPState = int32(iam_model.IDPConfigStateActive)
//	case domain.IDPConfigStateRemoved:
//		provider.IDPState = int32(iam_model.IDPConfigStateRemoved)
//	default:
//		provider.IDPState = int32(iam_model.IDPConfigStateActive)
//	}
//}
//
//func (i *IDPProvider) OnError(event *es_models.Event, err error) error {
//	logging.WithFields("id", event.AggregateID).WithError(err).Warn("something went wrong in idp provider handler")
//	return spooler.HandleError(event, err, i.view.GetLatestIDPProviderFailedEvent, i.view.ProcessedIDPProviderFailedEvent, i.view.ProcessedIDPProviderSequence, i.errorCountUntilSkip)
//}
//
//func (i *IDPProvider) OnSuccess(instanceIDs []string) error {
//	return spooler.HandleSuccess(i.view.UpdateIDPProviderSpoolerRunTimestamp, instanceIDs)
//}
//
//func (i *IDPProvider) getOrgIDPConfig(instanceID, aggregateID, idpConfigID string) (*query2.IDP, error) {
//	return i.queries.IDPByIDAndResourceOwner(withInstanceID(context.Background(), instanceID), false, idpConfigID, aggregateID, false)
//}
//
//func (i *IDPProvider) getDefaultIDPConfig(instanceID, idpConfigID string) (*query2.IDP, error) {
//	return i.queries.IDPByIDAndResourceOwner(withInstanceID(context.Background(), instanceID), false, idpConfigID, instanceID, false)
//}
