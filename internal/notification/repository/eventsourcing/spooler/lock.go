package spooler

import (
	"database/sql"
	"time"

	es_locker "github.com/zitadel/zitadel/internal/eventstore/v1/locker"
)

const (
	lockTable = "notification.locks"
)

type locker struct {
	dbClient *sql.DB
}

func (l *locker) Renew(lockerID, viewModel string, waitTime time.Duration) error {
	return es_locker.Renew(l.dbClient, lockTable, lockerID, viewModel, waitTime)
}
