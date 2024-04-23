package projection

import (
	"time"

	ms "github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/instance"
	"github.com/zitadel/zitadel/internal/v2/user"
)

type milestone struct {
	reachedAt time.Time

	typ ms.Type
}

type InstanceCreatedMilestone struct {
	milestone
}

func NewInstanceCreatedMilestone() *InstanceCreatedMilestone {
	return &InstanceCreatedMilestone{
		milestone: milestone{
			typ: ms.InstanceCreated,
		},
	}
}

func (p *InstanceCreatedMilestone) Filter() []*eventstore.Filter {
	if !p.reachedAt.IsZero() {
		return nil
	}

	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.AppendAggregateFilter(
				instance.AggregateType,
				eventstore.AppendEvent(
					eventstore.EventType("instance.added"),
				),
			),
			eventstore.FilterPagination(
				eventstore.Limit(1),
			),
		),
	}
}

func (p *InstanceCreatedMilestone) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		if event.Type != "instance.added" {
			continue
		}
		p.reachedAt = event.CreatedAt
	}
	return nil
}

type InstanceRemovedMilestone struct {
	milestone
}

func NewInstanceRemovedMilestone() *InstanceRemovedMilestone {
	return &InstanceRemovedMilestone{
		milestone: milestone{
			typ: ms.InstanceDeleted,
		},
	}
}

func (p *InstanceRemovedMilestone) Filter() []*eventstore.Filter {
	if !p.reachedAt.IsZero() {
		return nil
	}

	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.AppendAggregateFilter(
				instance.AggregateType,
				eventstore.AppendEvent(
					eventstore.EventType("instance.removed"),
				),
			),
			eventstore.FilterPagination(
				eventstore.Limit(1),
			),
		),
	}
}

func (p *InstanceRemovedMilestone) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		if event.Type != "instance.removed" {
			continue
		}

		p.reachedAt = event.CreatedAt
	}
	return nil
}

type AuthOnInstanceMilestone struct {
	milestone
}

func NewAuthOnInstanceMilestone() *AuthOnInstanceMilestone {
	return &AuthOnInstanceMilestone{
		milestone: milestone{
			typ: ms.AuthenticationSucceededOnInstance,
		},
	}
}

func (p *AuthOnInstanceMilestone) Filter() []*eventstore.Filter {
	if !p.reachedAt.IsZero() {
		return nil
	}

	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.AppendAggregateFilter(
				"user",
				eventstore.AppendEvent(
					eventstore.EventType("user.token.added"),
				),
			),
			eventstore.FilterPagination(
				eventstore.Limit(1),
			),
		),
	}
}

func (p *AuthOnInstanceMilestone) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		if event.Type != "user.token.added" {
			continue
		}

		p.reachedAt = event.CreatedAt
	}
	return nil
}

type AuthOnAppMilestone struct {
	milestone

	consoleAppID string

	position eventstore.GlobalPosition
}

func NewAuthOnAppMilestone() *AuthOnAppMilestone {
	return &AuthOnAppMilestone{
		milestone: milestone{
			typ: ms.AuthenticationSucceededOnApplication,
		},
	}
}

func (p *AuthOnAppMilestone) Filter() []*eventstore.Filter {
	if !p.reachedAt.IsZero() {
		return nil
	}

	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.AppendAggregateFilter(
				"user",
				eventstore.AppendEvent(
					eventstore.EventType("user.token.added"),
				),
			),
			eventstore.FilterPagination(
				// used because we need to check for first login and an app which is not console
				eventstore.PositionGreater(p.position.Position, p.position.InPositionOrder),
			),
		),
	}
}

func (p *AuthOnAppMilestone) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		if event.Type != "user.token.added" {
			continue
		}
		tokenAdded, err := user.TokenAddedEventFromStorage(event)
		if err != nil {
			return err
		}

		p.position = event.Position

		if tokenAdded.Payload.ApplicationID == p.consoleAppID {
			continue
		}
		p.reachedAt = event.CreatedAt
	}
	return nil
}

type ProjectCreatedMilestone struct {
	milestone
}

func NewProjectCreatedMilestone() *ProjectCreatedMilestone {
	return &ProjectCreatedMilestone{
		milestone: milestone{
			typ: ms.InstanceDeleted,
		},
	}
}

func (p *ProjectCreatedMilestone) Filter() []*eventstore.Filter {
	if !p.reachedAt.IsZero() {
		return nil
	}

	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.AppendAggregateFilter(
				"project",
				eventstore.AppendEvent(
					eventstore.EventType("project.added"),
					eventstore.EventCreatorsNotContains("", "SYSTEM"),
				),
			),
			eventstore.FilterPagination(
				eventstore.Limit(1),
			),
		),
	}
}

func (p *ProjectCreatedMilestone) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		if event.Type != "project.added" {
			continue
		}

		p.reachedAt = event.CreatedAt
	}
	return nil
}

type AppCreatedMilestone struct {
	milestone
}

func NewAppCreatedMilestone() *AppCreatedMilestone {
	return &AppCreatedMilestone{
		milestone: milestone{
			typ: ms.InstanceDeleted,
		},
	}
}

func (p *AppCreatedMilestone) Filter() []*eventstore.Filter {
	if !p.reachedAt.IsZero() {
		return nil
	}

	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.AppendAggregateFilter(
				"project",
				eventstore.AppendEvent(
					eventstore.EventCreatorsNotContains("", "SYSTEM"),
					eventstore.EventType("project.application.added"),
				),
			),
			eventstore.FilterPagination(
				eventstore.Limit(1),
			),
		),
	}
}

func (p *AppCreatedMilestone) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		if event.Type != "project.application.added" {
			continue
		}

		p.reachedAt = event.CreatedAt
	}
	return nil
}

/*
filter:
	instance: string
	owner: string
	pagination
	aggregates: []aggregateFilter

eventFilter:
	type textFilter[string]
	revision numberFilter[uint16]
	createdAt: timeFilter
	sequence: numberFilter[uint32]
	creator: textFilter[string] => not system
	pagination

aggregateFilter:
	type string
	id ListFilter
	events []eventFilter
	pagination

pagination:
	limit uint32
	offset uint32
	positionFilter

positionFilter:
	position float64
	in_tx_order uint32
*/

/*
milestonesFilter:
  instance: asdf
  aggregates:
    - type: instance
      events:
  	    - type: instance.added
		  pagination:
		    limit: 1
  	    - type: instance.removed
		  pagination:
		    limit: 1
  	    - type: instance.domain.primary.set
  	      creator: NotInList("", "system")
		  pagination:
		    limit: 1
    - type: project
      id: NotInFilter: (zitadel)
      events:
  	    - type: project.added
		  pagination:
		    limit: 1
  	    - type: project.application.added
		  pagination:
		    limit: 1
    - type: user
      events:
	    - user.token.added # query as long as e.applicationId is not empty and not console
		  pagination:
		    position:
			  position >= {{current position}}
			  in_tx_order > {{current in_tx_order}}
	- type: instance
	  pagination:
	    limit: 1
	  events:
	    - type: instance.idp.config.added
	    - type: instance.idp.oauth.added
	    - type: instance.idp.oidc.added
	    - type: instance.idp.jwt.added
	    - type: instance.idp.azure.added
	    - type: instance.idp.github.added
	    - type: instance.idp.github.enterprise.added
	    - type: instance.idp.gitlab.added
	    - type: instance.idp.gitlab.selfhosted.added
	    - type: instance.idp.google.added
	    - type: instance.idp.ldap.added
	    - type: instance.idp.config.apple.added
	    - type: instance.idp.saml.added
	- type: org
	  pagination:
	    limit: 1
	  events:
	    - type: org.idp.config.added
	    - type: org.idp.oauth.added
	    - type: org.idp.oidc.added
	    - type: org.idp.jwt.added
	    - type: org.idp.azure.added
	    - type: org.idp.github.added
	    - type: org.idp.github.enterprise.added
	    - type: org.idp.gitlab.added
	    - type: org.idp.gitlab.selfhosted.added
	    - type: org.idp.google.added
	    - type: org.idp.ldap.added
	    - type: org.idp.config.apple.added
	    - type: org.idp.saml.added
	- type: instance
	  pagination:
	    limit: 1
	  events:
	    - type: instance.login.policy.idp.added
	- type: org
	  pagination:
	    limit: 1
	  events:
	    - type: org.login.policy.idp.added
	- type: instance
	  events:
	  - type: instance.label.policy.added
	    pagination:
	      limit: 1
	  - type: instance.label.policy.activated
	    pagination:
	      limit: 1
	- type: org
	  events:
	  - type: org.label.policy.added
	    pagination:
	      limit: 1
	  - type: org.label.policy.activated
	    pagination:
	      limit: 1
	- type: instance
	  pagination:
	    limit: 1
	  events:
	    - type: instance.smtp.config.added
*/
