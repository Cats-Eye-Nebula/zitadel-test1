package handler

import (
	"time"

	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/query"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
)

type Configs map[string]*Config

type Config struct {
	MinimumCycleDuration types.Duration
}

type handler struct {
	view                *view.View
	bulkLimit           uint64
	cycleDuration       time.Duration
	errorCountUntilSkip uint64

	es eventstore.Eventstore
}

func (h *handler) Eventstore() eventstore.Eventstore {
	return h.es
}

type EventstoreRepos struct {
	IamEvents *iam_event.IAMEventstore
	OrgEvents *org_event.OrgEventstore
}

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, es eventstore.Eventstore, repos EventstoreRepos, defaults systemdefaults.SystemDefaults) []query.Handler {
	return []query.Handler{
		newOrg(
			handler{view, bulkLimit, configs.cycleDuration("Org"), errorCount, es}),
		newIAMMember(
			handler{view, bulkLimit, configs.cycleDuration("IamMember"), errorCount, es}),
		newIDPConfig(
			handler{view, bulkLimit, configs.cycleDuration("IDPConfig"), errorCount, es}),
		newLabelPolicy(
			handler{view, bulkLimit, configs.cycleDuration("LabelPolicy"), errorCount, es}),
		newLoginPolicy(
			handler{view, bulkLimit, configs.cycleDuration("LoginPolicy"), errorCount, es}),
		newIDPProvider(
			handler{view, bulkLimit, configs.cycleDuration("IDPProvider"), errorCount, es},
			defaults,
			repos.IamEvents,
			repos.OrgEvents),
		newUser(
			handler{view, bulkLimit, configs.cycleDuration("User"), errorCount, es},
			repos.OrgEvents,
			repos.IamEvents,
			defaults),
		newPasswordComplexityPolicy(
			handler{view, bulkLimit, configs.cycleDuration("PasswordComplexityPolicy"), errorCount, es}),
		newPasswordAgePolicy(
			handler{view, bulkLimit, configs.cycleDuration("PasswordAgePolicy"), errorCount, es}),
		newPasswordLockoutPolicy(
			handler{view, bulkLimit, configs.cycleDuration("PasswordLockoutPolicy"), errorCount, es}),
		newOrgIAMPolicy(
			handler{view, bulkLimit, configs.cycleDuration("OrgIAMPolicy"), errorCount, es}),
		newExternalIDP(
			handler{view, bulkLimit, configs.cycleDuration("ExternalIDP"), errorCount, es},
			defaults,
			repos.IamEvents,
			repos.OrgEvents),
		newMailTemplate(
			handler{view, bulkLimit, configs.cycleDuration("MailTemplate"), errorCount, es}),
		newMailText(
			handler{view, bulkLimit, configs.cycleDuration("MailText"), errorCount, es}),
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 3 * time.Minute
	}
	return c.MinimumCycleDuration.Duration
}

func (h *handler) MinimumCycleDuration() time.Duration {
	return h.cycleDuration
}

func (h *handler) LockDuration() time.Duration {
	return h.cycleDuration / 3
}

func (h *handler) QueryLimit() uint64 {
	return h.bulkLimit
}
