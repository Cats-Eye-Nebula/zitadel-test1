package notification

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/command"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing"
	"github.com/rakyll/statik/fs"

	_ "github.com/caos/zitadel/internal/notification/statik"
)

type Config struct {
	APIDomain  string
	Repository eventsourcing.Config
}

func Start(ctx context.Context, config Config, systemDefaults sd.SystemDefaults, command *command.Commands, hasStatics bool) {
	statikFS, err := fs.NewWithNamespace("notification")
	logging.Log("CONFI-7usEW").OnError(err).Panic("unable to start listener")

	apiDomain := config.APIDomain
	if !hasStatics {
		apiDomain = ""
	}
	_, err = eventsourcing.Start(config.Repository, statikFS, systemDefaults, command, apiDomain)
	logging.Log("MAIN-9uBxp").OnError(err).Panic("unable to start app")
}
