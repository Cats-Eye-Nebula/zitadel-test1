package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	OrgAggregate       models.AggregateType = "org"
	OrgDomainAggregate models.AggregateType = "org.domain"
	OrgNameAggregate   models.AggregateType = "org.name"

	OrgAdded                    models.EventType = "org.added"
	OrgChanged                  models.EventType = "org.changed"
	OrgDeactivated              models.EventType = "org.deactivated"
	OrgReactivated              models.EventType = "org.reactivated"
	OrgRemoved                  models.EventType = "org.removed"
	OrgDomainAdded              models.EventType = "org.domain.added"
	OrgDomainVerificationAdded  models.EventType = "org.domain.verification.added"
	OrgDomainVerificationFailed models.EventType = "org.domain.verification.failed"
	OrgDomainVerified           models.EventType = "org.domain.verified"
	OrgDomainRemoved            models.EventType = "org.domain.removed"
	OrgDomainPrimarySet         models.EventType = "org.domain.primary.set"

	OrgNameReserved models.EventType = "org.name.reserved"
	OrgNameReleased models.EventType = "org.name.released"

	OrgDomainReserved models.EventType = "org.domain.reserved"
	OrgDomainReleased models.EventType = "org.domain.released"

	OrgMemberAdded   models.EventType = "org.member.added"
	OrgMemberChanged models.EventType = "org.member.changed"
	OrgMemberRemoved models.EventType = "org.member.removed"

	OrgIAMPolicyAdded   models.EventType = "org.iam.policy.added"
	OrgIAMPolicyChanged models.EventType = "org.iam.policy.changed"
	OrgIAMPolicyRemoved models.EventType = "org.iam.policy.removed"

	IDPConfigAdded       models.EventType = "org.idp.config.added"
	IDPConfigChanged     models.EventType = "org.idp.config.changed"
	IDPConfigRemoved     models.EventType = "org.idp.config.removed"
	IDPConfigDeactivated models.EventType = "org.idp.config.deactivated"
	IDPConfigReactivated models.EventType = "org.idp.config.reactivated"

	OIDCIDPConfigAdded   models.EventType = "org.idp.oidc.config.added"
	OIDCIDPConfigChanged models.EventType = "org.idp.oidc.config.changed"

	SAMLIDPConfigAdded   models.EventType = "org.idp.saml.config.added"
	SAMLIDPConfigChanged models.EventType = "org.idp.saml.config.changed"

	LoginPolicyAdded                     models.EventType = "org.policy.login.added"
	LoginPolicyChanged                   models.EventType = "org.policy.login.changed"
	LoginPolicyRemoved                   models.EventType = "org.policy.login.removed"
	LoginPolicyIDPProviderAdded          models.EventType = "org.policy.login.idpprovider.added"
	LoginPolicyIDPProviderRemoved        models.EventType = "org.policy.login.idpprovider.removed"
	LoginPolicyIDPProviderCascadeRemoved models.EventType = "org.policy.login.idpprovider.cascade.removed"
	LabelPolicyAdded                     models.EventType = "org.policy.label.added"
	LabelPolicyChanged                   models.EventType = "org.policy.label.changed"
	LabelPolicyRemoved                   models.EventType = "org.policy.label.removed"
	LabelPolicyIDPProviderAdded          models.EventType = "org.policy.label.idpprovider.added"
	LabelPolicyIDPProviderRemoved        models.EventType = "org.policy.label.idpprovider.removed"
	LabelPolicyIDPProviderCascadeRemoved models.EventType = "org.policy.label.idpprovider.cascade.removed"
)
