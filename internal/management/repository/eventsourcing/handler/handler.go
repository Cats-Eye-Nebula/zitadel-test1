package handler

import (
	"time"

	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/static"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
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

	es v1.Eventstore
}

func (h *handler) Eventstore() v1.Eventstore {
	return h.es
}

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, es v1.Eventstore, defaults systemdefaults.SystemDefaults, staticStorage static.Storage) []query.Handler {
	return []query.Handler{
		newProjectMember(handler{view, bulkLimit, configs.cycleDuration("ProjectMember"), errorCount, es}),
		newProjectGrantMember(handler{view, bulkLimit, configs.cycleDuration("ProjectGrantMember"), errorCount, es}),
		newUser(handler{view, bulkLimit, configs.cycleDuration("User"), errorCount, es},
			defaults.IamID),
		newUserGrant(handler{view, bulkLimit, configs.cycleDuration("UserGrant"), errorCount, es}),
		newOrgMember(
			handler{view, bulkLimit, configs.cycleDuration("OrgMember"), errorCount, es}),
		newUserMembership(
			handler{view, bulkLimit, configs.cycleDuration("UserMembership"), errorCount, es}),
		newAuthNKeys(
			handler{view, bulkLimit, configs.cycleDuration("MachineKeys"), errorCount, es}),
		newIDPConfig(
			handler{view, bulkLimit, configs.cycleDuration("IDPConfig"), errorCount, es}),
		newIDPProvider(
			handler{view, bulkLimit, configs.cycleDuration("IDPProvider"), errorCount, es},
			defaults),
		newExternalIDP(
			handler{view, bulkLimit, configs.cycleDuration("ExternalIDP"), errorCount, es},
			defaults),
		newMailTemplate(
			handler{view, bulkLimit, configs.cycleDuration("MailTemplate"), errorCount, es}),
		newMessageText(
			handler{view, bulkLimit, configs.cycleDuration("MessageText"), errorCount, es}),
		newCustomText(
			handler{view, bulkLimit, configs.cycleDuration("CustomText"), errorCount, es}),
		newMetadata(
			handler{view, bulkLimit, configs.cycleDuration("Metadata"), errorCount, es}),
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
