package spooler

import (
	"context"
	"strconv"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/view/repository"

	"time"
)

type Spooler struct {
	handlers   []query.Handler
	locker     Locker
	lockID     string
	eventstore eventstore.Eventstore
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
	eventstore eventstore.Eventstore
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
				go func(handler *spooledHandler, queue chan<- *spooledHandler) {
					time.Sleep(handler.MinimumCycleDuration() - time.Since(handler.queuedAt))
					handler.queuedAt = time.Now()
					queue <- handler
				}(task, s.queue)

				task.load(workerID)
			}
		}(i)
	}
	for _, handler := range s.handlers {
		handler := &spooledHandler{Handler: handler, locker: s.locker, queuedAt: time.Now(), eventstore: s.eventstore}
		s.queue <- handler
	}
}

func (s *spooledHandler) load(workerID string) {
	errs := make(chan error)
	defer close(errs)
	ctx, cancel := context.WithCancel(context.Background())
	go s.awaitError(cancel, errs, workerID)
	hasLocked := s.lock(ctx, errs, workerID)

	if <-hasLocked {
		go func() {
			for l := range hasLocked {
				if !l {
					// we only need to break. An error is already written by the lock-routine to the errs channel
					break
				}
			}
		}()
		events, err := s.query(ctx)
		if err != nil {
			errs <- err
		} else {
			errs <- s.process(ctx, events, workerID)
			logging.Log("SPOOL-0pV8o").WithField("view", s.ViewModel()).WithField("worker", workerID).Debug("process done")
		}
	}
	<-ctx.Done()
}

func (s *spooledHandler) awaitError(cancel func(), errs chan error, workerID string) {
	select {
	case err := <-errs:
		cancel()
		logging.Log("SPOOL-K2lst").OnError(err).WithField("view", s.ViewModel()).WithField("worker", workerID).Debug("load canceled")
	}
}

func (s *spooledHandler) process(ctx context.Context, events []*models.Event, workerID string) error {
	for _, event := range events {
		select {
		case <-ctx.Done():
			logging.LogWithFields("SPOOL-FTKwH", "view", s.ViewModel(), "worker", workerID).Debug("context canceled")
			return nil
		default:
			if err := s.Reduce(event); err != nil {
				return s.OnError(event, err)
			}
		}
	}
	return nil
}

func HandleError(event *models.Event, failedErr error,
	latestFailedEvent func(sequence uint64) (*repository.FailedEvent, error),
	processFailedEvent func(*repository.FailedEvent) error,
	processSequence func(uint64) error, errorCountUntilSkip uint64) error {
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
		return processSequence(event.Sequence)
	}
	return nil
}

func (s *spooledHandler) query(ctx context.Context) ([]*models.Event, error) {
	query, err := s.EventQuery()
	if err != nil {
		return nil, err
	}
	factory := models.FactoryFromSearchQuery(query)
	sequence, err := s.eventstore.LatestSequence(ctx, factory)
	logging.Log("SPOOL-7SciK").OnError(err).Debug("unable to query latest sequence")
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

func (s *spooledHandler) lock(ctx context.Context, errs chan<- error, workerID string) chan bool {
	renewTimer := time.After(0)
	renewDuration := s.MinimumCycleDuration()
	locked := make(chan bool)

	go func(locked chan bool) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-renewTimer:
				logging.Log("SPOOL-K2lst").WithField("view", s.ViewModel()).WithField("worker", workerID).Debug("renew")
				err := s.locker.Renew(workerID, s.ViewModel(), s.MinimumCycleDuration()*2)
				logging.Log("SPOOL-K2lst").WithField("view", s.ViewModel()).WithField("worker", workerID).WithError(err).Debug("renew done")
				if err == nil {
					locked <- true
					renewTimer = time.After(renewDuration)
					continue
				}

				if ctx.Err() == nil {
					errs <- err
				}

				locked <- false
				return
			}
		}
	}(locked)

	return locked
}
