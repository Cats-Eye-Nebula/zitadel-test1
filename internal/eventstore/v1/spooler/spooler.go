package spooler

import (
	"context"

	"github.com/getsentry/sentry-go"

	"strconv"
	"sync"
	"time"

	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/view/repository"
)

type Spooler struct {
	handlers   []query.Handler
	locker     Locker
	lockID     string
	eventstore v1.Eventstore
	workers    int
	queue      chan *spooledHandler
}

type Locker interface {
	Renew(lockerID, viewModel string, waitTime time.Duration) error
}

type spooledHandler struct {
	query.Handler
	locker     Locker
	queuedAt   time.Time
	eventstore v1.Eventstore
}

func (s *Spooler) Start() {
	defer logging.LogWithFields("SPOOL-N0V1g", "lockerID", s.lockID, "workers", s.workers).Info("spooler started")
	if s.workers < 1 {
		return
	}

	for i := 0; i < s.workers; i++ {
		go func(workerIdx int) {
			workerID := s.lockID + "--" + strconv.Itoa(workerIdx)
			for task := range s.queue {
				go requeueTask(task, s.queue)
				task.load(workerID)
			}
		}(i)
	}
	go func() {
		for _, handler := range s.handlers {
			s.queue <- &spooledHandler{Handler: handler, locker: s.locker, queuedAt: time.Now(), eventstore: s.eventstore}
		}
	}()
}

func requeueTask(task *spooledHandler, queue chan<- *spooledHandler) {
	time.Sleep(task.MinimumCycleDuration() - time.Since(task.queuedAt))
	task.queuedAt = time.Now()
	queue <- task
}

func (s *spooledHandler) load(workerID string) {
	errs := make(chan error)
	defer func() {
		close(errs)
		err := recover()

		if err != nil {
			sentry.CurrentHub().Recover(err)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	go s.awaitError(cancel, errs, workerID)
	hasLocked := s.lock(ctx, errs, workerID)

	if <-hasLocked {
		for {
			events, err := s.query(ctx)
			if err != nil {
				errs <- err
				break
			}
			err = s.process(ctx, events, workerID)
			if err != nil {
				errs <- err
				break
			}
			if uint64(len(events)) < s.QueryLimit() {
				// no more events to process
				// stop chan
				if ctx.Err() == nil {
					errs <- nil
				}
				break
			}
		}
	}
	<-ctx.Done()
}

func (s *spooledHandler) awaitError(cancel func(), errs chan error, workerID string) {
	select {
	case err := <-errs:
		cancel()
		logging.Log("SPOOL-OT8di").OnError(err).WithField("view", s.ViewModel()).WithField("worker", workerID).Debug("load canceled")
	}
}

func (s *spooledHandler) process(ctx context.Context, events []*models.Event, workerID string) error {
	for i, event := range events {
		select {
		case <-ctx.Done():
			logging.LogWithFields("SPOOL-FTKwH", "view", s.ViewModel(), "worker", workerID, "traceID", tracing.TraceIDFromCtx(ctx)).Debug("context canceled")
			return nil
		default:
			if err := s.Reduce(event); err != nil {
				err = s.OnError(event, err)
				if err == nil {
					continue
				}
				time.Sleep(100 * time.Millisecond)
				return s.process(ctx, events[i:], workerID)
			}
		}
	}
	err := s.OnSuccess()
	logging.LogWithFields("SPOOL-49ods", "view", s.ViewModel(), "worker", workerID, "traceID", tracing.TraceIDFromCtx(ctx)).OnError(err).Warn("could not process on success func")
	return err
}

func (s *spooledHandler) query(ctx context.Context) ([]*models.Event, error) {
	query, err := s.EventQuery()
	if err != nil {
		return nil, err
	}
	factory := models.FactoryFromSearchQuery(query)
	sequence, err := s.eventstore.LatestSequence(ctx, factory)
	logging.Log("SPOOL-7SciK").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to query latest sequence")
	var processedSequence uint64
	for _, filter := range query.Filters {
		if filter.GetField() == models.Field_LatestSequence {
			processedSequence = filter.GetValue().(uint64)
		}
	}
	if sequence != 0 && processedSequence == sequence {
		return nil, nil
	}

	query.Limit = s.QueryLimit()
	return s.eventstore.FilterEvents(ctx, query)
}

//lock ensures the lock on the database.
// the returned channel will be closed if ctx is done or an error occured durring lock
func (s *spooledHandler) lock(ctx context.Context, errs chan<- error, workerID string) chan bool {
	renewTimer := time.After(0)
	locked := make(chan bool)

	go func(locked chan bool) {
		var firstLock sync.Once
		defer close(locked)
		for {
			select {
			case <-ctx.Done():
				return
			case <-renewTimer:
				err := s.locker.Renew(workerID, s.ViewModel(), s.LockDuration())
				firstLock.Do(func() {
					locked <- err == nil
				})
				if err == nil {
					renewTimer = time.After(s.LockDuration())
					continue
				}

				if ctx.Err() == nil {
					errs <- err
				}
				return
			}
		}
	}(locked)

	return locked
}

func HandleError(event *models.Event, failedErr error,
	latestFailedEvent func(sequence uint64) (*repository.FailedEvent, error),
	processFailedEvent func(*repository.FailedEvent) error,
	processSequence func(*models.Event) error,
	errorCountUntilSkip uint64) error {
	failedEvent, err := latestFailedEvent(event.Sequence)
	if err != nil {
		return err
	}
	failedEvent.FailureCount++
	failedEvent.ErrMsg = failedErr.Error()
	err = processFailedEvent(failedEvent)
	if err != nil {
		return err
	}
	if errorCountUntilSkip <= failedEvent.FailureCount {
		return processSequence(event)
	}
	return failedErr
}

func HandleSuccess(updateSpoolerRunTimestamp func() error) error {
	return updateSpoolerRunTimestamp()
}
