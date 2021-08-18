package crdb

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/errors"
)

const (
	lockStmtFormat = "INSERT INTO %[1]s" +
		" (locker_id, locked_until, projection_name) VALUES ($1, now()+$2::INTERVAL, $3)" +
		" ON CONFLICT (projection_name)" +
		" DO UPDATE SET locker_id = $1, locked_until = now()+$2::INTERVAL" +
		" WHERE %[1]s.projection_name = $3 AND (%[1]s.locker_id = $1 OR %[1]s.locked_until < now())"
)

func (h *StatementHandler) Lock(ctx context.Context, lockDuration time.Duration) <-chan error {
	errs := make(chan error)
	go h.handleLock(ctx, errs, lockDuration)
	return errs
}

func (h *StatementHandler) handleLock(ctx context.Context, errs chan error, lockDuration time.Duration) {
	renewLock := time.NewTimer(0)
	for {
		select {
		case <-renewLock.C:
			errs <- h.renewLock(ctx, lockDuration)
			//refresh the lock 500ms before it times out. 500ms should be enough for one transaction
			renewLock.Reset(lockDuration - (500 * time.Millisecond))
		case <-ctx.Done():
			close(errs)
			renewLock.Stop()
			return
		}
	}
}

func (h *StatementHandler) renewLock(ctx context.Context, lockDuration time.Duration) error {
	//the unit of crdb interval is seconds (https://www.cockroachlabs.com/docs/stable/interval.html).
	//the var is float for not losing ms of the lock
	res, err := h.client.Exec(h.lockStmt, h.workerName, lockDuration.Seconds(), h.ProjectionName)
	if err != nil {
		return errors.ThrowInternal(err, "CRDB-uaDoR", "unable to execute lock")
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.ThrowAlreadyExists(nil, "CRDB-mmi4J", "projection already locked")
	}

	return nil
}

func (h *StatementHandler) Unlock() error {
	_, err := h.client.Exec(h.lockStmt, h.workerName, float64(0), h.ProjectionName)
	if err != nil {
		return errors.ThrowUnknown(err, "CRDB-JjfwO", "unlock failed")
	}
	return nil
}
