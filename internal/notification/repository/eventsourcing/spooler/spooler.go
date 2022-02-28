package spooler

import (
	"database/sql"
	"net/http"

	"github.com/caos/zitadel/internal/command"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/query"
)

type SpoolerConfig struct {
	BulkLimit             uint64
	FailureCountUntilSkip uint64
	ConcurrentWorkers     int
	Handlers              handler.Configs
}

func StartSpooler(c SpoolerConfig,
	es v1.Eventstore,
	view *view.View,
	sql *sql.DB,
	command *command.Commands,
	queries *query.Queries,
	systemDefaults sd.SystemDefaults,
	dir http.FileSystem,
	assetsPrefix string,
	userEncryption crypto.EncryptionAlgorithm,
	smtpEncryption crypto.EncryptionAlgorithm,
) *spooler.Spooler {
	spoolerConfig := spooler.Config{
		Eventstore:        es,
		Locker:            &locker{dbClient: sql},
		ConcurrentWorkers: c.ConcurrentWorkers,
		ViewHandlers:      handler.Register(c.Handlers, c.BulkLimit, c.FailureCountUntilSkip, view, es, command, queries, systemDefaults, dir, assetsPrefix, userEncryption, smtpEncryption),
	}
	spool := spoolerConfig.New()
	spool.Start()
	return spool
}
