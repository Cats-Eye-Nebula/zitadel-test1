package handler

import (
	"time"

	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/static"
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

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, es v1.Eventstore, defaults systemdefaults.SystemDefaults, command *command.Commands, static static.Storage, localDevMode bool) []query.Handler {
	handlers := []query.Handler{
		newIAMMember(
			handler{view, bulkLimit, configs.cycleDuration("IamMember"), errorCount, es}),
		newIDPConfig(
			handler{view, bulkLimit, configs.cycleDuration("IDPConfig"), errorCount, es}),
		newLabelPolicy(
			handler{view, bulkLimit, configs.cycleDuration("LabelPolicy"), errorCount, es}),
		newIDPProvider(
			handler{view, bulkLimit, configs.cycleDuration("IDPProvider"), errorCount, es},
			defaults),
		newUser(
			handler{view, bulkLimit, configs.cycleDuration("User"), errorCount, es},
			defaults),
		newExternalIDP(
			handler{view, bulkLimit, configs.cycleDuration("ExternalIDP"), errorCount, es},
			defaults),
		newMailTemplate(
			handler{view, bulkLimit, configs.cycleDuration("MailTemplate"), errorCount, es}),
		newMessageText(
			handler{view, bulkLimit, configs.cycleDuration("MessageText"), errorCount, es}),
		newFeatures(
			handler{view, bulkLimit, configs.cycleDuration("Features"), errorCount, es}),
		newCustomText(
			handler{view, bulkLimit, configs.cycleDuration("CustomTexts"), errorCount, es}),
	}
	if static != nil {
		handlers = append(handlers, newStyling(
			handler{view, bulkLimit, configs.cycleDuration("Styling"), errorCount, es},
			static,
			localDevMode))
	}
	return handlers
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
