package idp

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
)

type GitHubIDPAddedEvent struct {
	OAuthIDPAddedEvent
}

func NewGitHubIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options Options,
) *GitHubIDPAddedEvent {
	return &GitHubIDPAddedEvent{
		OAuthIDPAddedEvent: OAuthIDPAddedEvent{
			BaseEvent:    *base,
			ID:           id,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       scopes,
			Options:      options,
		},
	}
}

func (e *GitHubIDPAddedEvent) Data() interface{} {
	return e
}

func (e *GitHubIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
	//return []*eventstore.EventUniqueConstraint{idpconfig.NewAddIDPConfigNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)}
}

func GitHubIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := OAuthIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubIDPAddedEvent{OAuthIDPAddedEvent: *e.(*OAuthIDPAddedEvent)}, nil
}

type GitHubIDPChangedEvent struct {
	OAuthIDPChangedEvent
}

func NewGitHubIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []OAuthIDPChanges,
) (*GitHubIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-BH3dl", "Errors.NoChangesFound")
	}
	changedEvent := &GitHubIDPChangedEvent{
		OAuthIDPChangedEvent: OAuthIDPChangedEvent{
			BaseEvent: *base,
			ID:        id,
		},
	}
	for _, change := range changes {
		change(&changedEvent.OAuthIDPChangedEvent)
	}
	return changedEvent, nil
}

//
//type OAuthIDPChanges func(*OAuthIDPChangedEvent)
//
//func ChangeOAuthName(name string) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.Name = &name
//	}
//}
//func ChangeOAuthClientID(clientID string) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.ClientID = &clientID
//	}
//}
//
//func ChangeOAuthClientSecret(clientSecret *crypto.CryptoValue) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.ClientSecret = clientSecret
//	}
//}
//
//func ChangeOAuthOptions(options OptionChanges) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.OptionChanges = options
//	}
//}
//
//func ChangeOAuthAuthorizationEndpoint(authorizationEndpoint string) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.AuthorizationEndpoint = &authorizationEndpoint
//	}
//}
//
//func ChangeOAuthTokenEndpoint(tokenEndpoint string) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.TokenEndpoint = &tokenEndpoint
//	}
//}
//
//func ChangeOAuthUserEndpoint(userEndpoint string) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.UserEndpoint = &userEndpoint
//	}
//}
//
//func ChangeOAuthScopes(scopes []string) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.Scopes = scopes
//	}
//}

func (e *GitHubIDPChangedEvent) Data() interface{} {
	return e
}

func (e *GitHubIDPChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
	//if e.Name == nil || e.oldName == *e.Name { // TODO: nil check should be enough
	//	return nil
	//}
	//return []*eventstore.EventUniqueConstraint{
	//	idpconfig.NewRemoveIDPConfigNameUniqueConstraint(e.oldName, e.Aggregate().ResourceOwner),
	//	idpconfig.NewAddIDPConfigNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner),
	//}
}

func GitHubIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := OAuthIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubIDPChangedEvent{OAuthIDPChangedEvent: *e.(*OAuthIDPChangedEvent)}, nil
}

//
//func OAuthIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
//	e := &OAuthIDPChangedEvent{
//		BaseEvent: *eventstore.BaseEventFromRepo(event),
//	}
//
//	err := json.Unmarshal(event.Data, e)
//	if err != nil {
//		return nil, errors.ThrowInternal(err, "IDP-D3gjzh", "unable to unmarshal event")
//	}
//
//	return e, nil
//}

type GitHubEnterpriseIDPAddedEvent struct {
	OAuthIDPAddedEvent
}

func NewGitHubEnterpriseIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint string,
	scopes []string,
	options Options,
) *GitHubEnterpriseIDPAddedEvent {
	return &GitHubEnterpriseIDPAddedEvent{
		OAuthIDPAddedEvent: *NewOAuthIDPAddedEvent(
			base,
			id,
			name,
			clientID,
			clientSecret,
			authorizationEndpoint,
			tokenEndpoint,
			userEndpoint,
			scopes,
			options,
		),
	}
}

func (e *GitHubEnterpriseIDPAddedEvent) Data() interface{} {
	return e
}

func (e *GitHubEnterpriseIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{idpconfig.NewAddIDPConfigNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)}
}

func GitHubEnterpriseIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := OAuthIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubEnterpriseIDPAddedEvent{OAuthIDPAddedEvent: *e.(*OAuthIDPAddedEvent)}, nil
}

type GitHubEnterpriseIDPChangedEvent struct {
	OAuthIDPChangedEvent
}

func NewGitHubEnterpriseIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []OAuthIDPChanges,
) (*GitHubEnterpriseIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-JHKs9", "Errors.NoChangesFound")
	}
	changedEvent := &GitHubEnterpriseIDPChangedEvent{
		OAuthIDPChangedEvent: OAuthIDPChangedEvent{
			BaseEvent: *base,
			ID:        id,
		},
	}
	for _, change := range changes {
		change(&changedEvent.OAuthIDPChangedEvent)
	}
	return changedEvent, nil
}

func (e *GitHubEnterpriseIDPChangedEvent) Data() interface{} {
	return e
}

func (e *GitHubEnterpriseIDPChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	if e.Name == nil || e.oldName == *e.Name { // TODO: nil check should be enough
		return nil
	}
	return []*eventstore.EventUniqueConstraint{
		idpconfig.NewRemoveIDPConfigNameUniqueConstraint(e.oldName, e.Aggregate().ResourceOwner),
		idpconfig.NewAddIDPConfigNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner),
	}
}

func GitHubEnterpriseIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := OAuthIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubEnterpriseIDPChangedEvent{OAuthIDPChangedEvent: *e.(*OAuthIDPChangedEvent)}, nil
}
