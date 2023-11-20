package org

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/crypto"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/repository/idpconfig"
)

const (
	IDPOIDCConfigAddedEventType   eventstore.EventType = "org.idp." + idpconfig.OIDCConfigAddedEventType
	IDPOIDCConfigChangedEventType eventstore.EventType = "org.idp." + idpconfig.OIDCConfigChangedEventType
)

type IDPOIDCConfigAddedEvent struct {
	idpconfig.OIDCConfigAddedEvent
}

func NewIDPOIDCConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	clientID,
	idpConfigID,
	issuer,
	authorizationEndpoint,
	tokenEndpoint string,
	clientSecret *crypto.CryptoValue,
	idpDisplayNameMapping,
	userNameMapping domain.OIDCMappingField,
	scopes ...string,
) *IDPOIDCConfigAddedEvent {

	return &IDPOIDCConfigAddedEvent{
		OIDCConfigAddedEvent: *idpconfig.NewOIDCConfigAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				IDPOIDCConfigAddedEventType,
			),
			clientID,
			idpConfigID,
			issuer,
			authorizationEndpoint,
			tokenEndpoint,
			clientSecret,
			idpDisplayNameMapping,
			userNameMapping,
			scopes...,
		),
	}
}

func IDPOIDCConfigAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idpconfig.OIDCConfigAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPOIDCConfigAddedEvent{OIDCConfigAddedEvent: *e.(*idpconfig.OIDCConfigAddedEvent)}, nil
}

type IDPOIDCConfigChangedEvent struct {
	idpconfig.OIDCConfigChangedEvent
}

func NewIDPOIDCConfigChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID string,
	changes []idpconfig.OIDCConfigChanges,
) (*IDPOIDCConfigChangedEvent, error) {
	changeEvent, err := idpconfig.NewOIDCConfigChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			IDPOIDCConfigChangedEventType),
		idpConfigID,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &IDPOIDCConfigChangedEvent{OIDCConfigChangedEvent: *changeEvent}, nil
}

func IDPOIDCConfigChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idpconfig.OIDCConfigChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPOIDCConfigChangedEvent{OIDCConfigChangedEvent: *e.(*idpconfig.OIDCConfigChangedEvent)}, nil
}
