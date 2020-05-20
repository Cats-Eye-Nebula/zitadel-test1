package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	ProjectAggregate models.AggregateType = "project"

	ProjectAdded       models.EventType = "project.added"
	ProjectChanged     models.EventType = "project.changed"
	ProjectDeactivated models.EventType = "project.deactivated"
	ProjectReactivated models.EventType = "project.reactivated"

	ProjectMemberAdded   models.EventType = "project.member.added"
	ProjectMemberChanged models.EventType = "project.member.changed"
	ProjectMemberRemoved models.EventType = "project.member.removed"

	ProjectRoleAdded   models.EventType = "project.role.added"
	ProjectRoleChanged models.EventType = "project.role.changed"
	ProjectRoleRemoved models.EventType = "project.role.removed"

	ProjectGrantAdded       models.EventType = "project.grant.added"
	ProjectGrantChanged     models.EventType = "project.grant.changed"
	ProjectGrantRemoved     models.EventType = "project.grant.removed"
	ProjectGrantDeactivated models.EventType = "project.grant.deactivated"
	ProjectGrantReactivated models.EventType = "project.grant.reactivated"

	ProjectGrantMemberAdded   models.EventType = "project.grant.member.added"
	ProjectGrantMemberChanged models.EventType = "project.grant.member.changed"
	ProjectGrantMemberRemoved models.EventType = "project.grant.member.removed"

	ApplicationAdded       models.EventType = "project.application.added"
	ApplicationChanged     models.EventType = "project.application.changed"
	ApplicationRemoved     models.EventType = "project.application.removed"
	ApplicationDeactivated models.EventType = "project.application.deactivated"
	ApplicationReactivated models.EventType = "project.application.reactivated"

	OIDCConfigAdded                models.EventType = "project.application.config.oidc.added"
	OIDCConfigChanged              models.EventType = "project.application.config.oidc.changed"
	OIDCConfigSecretChanged        models.EventType = "project.application.config.oidc.secret.changed"
	OIDCClientSecretCheckSucceeded models.EventType = "project.application.oidc.secret.check.succeeded"
	OIDCClientSecretCheckFailed    models.EventType = "project.application.oidc.secret.check.failed"
)
