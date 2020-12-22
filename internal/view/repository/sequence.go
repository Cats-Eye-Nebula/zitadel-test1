package repository

import (
	"strings"
	"time"

	caos_errs "github.com/caos/zitadel/internal/errors"
	int_model "github.com/caos/zitadel/internal/model"
	"github.com/caos/zitadel/internal/view/model"
	"github.com/jinzhu/gorm"
)

type CurrentSequence struct {
	ViewName                 string    `gorm:"column:view_name;primary_key"`
	CurrentSequence          uint64    `gorm:"column:current_sequence"`
	EventTimestamp           time.Time `gorm:"column:event_timestamp"`
	LastSuccessfulSpoolerRun time.Time `gorm:"column:last_successful_spooler_run"`
	AggregateType            string    `gorm:"column:aggregate_type;primary_key"`
}

type currentSequenceViewWithSequence struct {
	ViewName                 string    `gorm:"column:view_name;primary_key"`
	CurrentSequence          uint64    `gorm:"column:current_sequence"`
	LastSuccessfulSpoolerRun time.Time `gorm:"column:last_successful_spooler_run"`
}

type currentSequenceView struct {
	ViewName                 string    `gorm:"column:view_name;primary_key"`
	LastSuccessfulSpoolerRun time.Time `gorm:"column:last_successful_spooler_run"`
}

type SequenceSearchKey int32

const (
	SequenceSearchKeyUndefined SequenceSearchKey = iota
	SequenceSearchKeyViewName
	SequenceSearchKeyAggregateType
)

type sequenceSearchKey SequenceSearchKey

func (key sequenceSearchKey) ToColumnName() string {
	switch SequenceSearchKey(key) {
	case SequenceSearchKeyViewName:
		return "view_name"
	case SequenceSearchKeyAggregateType:
		return "aggregate_type"
	default:
		return ""
	}
}

type sequenceSearchQuery struct {
	key   sequenceSearchKey
	value string
}

func (q *sequenceSearchQuery) GetKey() ColumnKey {
	return q.key
}

func (q *sequenceSearchQuery) GetMethod() int_model.SearchMethod {
	return int_model.SearchMethodEquals
}

func (q *sequenceSearchQuery) GetValue() interface{} {
	return q.value
}

func CurrentSequenceToModel(sequence *CurrentSequence) *model.View {
	dbView := strings.Split(sequence.ViewName, ".")
	return &model.View{
		Database:                 dbView[0],
		ViewName:                 dbView[1],
		CurrentSequence:          sequence.CurrentSequence,
		EventTimestamp:           sequence.EventTimestamp,
		LastSuccessfulSpoolerRun: sequence.LastSuccessfulSpoolerRun,
		AggregateType:            sequence.AggregateType,
	}
}

func SaveCurrentSequence(db *gorm.DB, table, viewName, aggregateType string, sequence uint64, eventTimestamp time.Time) error {
	return UpdateCurrentSequence(db, table, &CurrentSequence{viewName, sequence, eventTimestamp, time.Now(), aggregateType})
}

func UpdateCurrentSequence(db *gorm.DB, table string, currentSequence *CurrentSequence) (err error) {
	var seq interface{} = currentSequence
	if currentSequence.AggregateType == "" && currentSequence.CurrentSequence > 0 {
		//spooler run
		seq = &currentSequenceView{ViewName: currentSequence.ViewName, LastSuccessfulSpoolerRun: currentSequence.LastSuccessfulSpoolerRun}
	} else if currentSequence.AggregateType == "" {
		//reset current sequence on view
		seq = &currentSequenceViewWithSequence{ViewName: currentSequence.ViewName, LastSuccessfulSpoolerRun: currentSequence.LastSuccessfulSpoolerRun, CurrentSequence: currentSequence.CurrentSequence}
	}

	save := PrepareSave(table)
	err = save(db, seq)
	if err != nil {
		return caos_errs.ThrowInternal(err, "VIEW-5kOhP", "unable to updated processed sequence")
	}
	return nil
}

func LatestSequence(db *gorm.DB, table, viewName, aggregateType string) (*CurrentSequence, error) {
	searchQueries := make([]SearchQuery, 0, 2)
	searchQueries = append(searchQueries, &sequenceSearchQuery{key: sequenceSearchKey(SequenceSearchKeyViewName), value: viewName})
	if aggregateType != "" {
		searchQueries = append(searchQueries, &sequenceSearchQuery{key: sequenceSearchKey(SequenceSearchKeyAggregateType), value: aggregateType})
	} else {
		// ensure highest sequence of view
		db = db.Order("current_sequence DESC")
	}

	query := PrepareGetByQuery(table, searchQueries...)
	sequence := new(CurrentSequence)
	err := query(db, sequence)

	if err == nil {
		return sequence, nil
	}

	if caos_errs.IsNotFound(err) {
		return sequence, nil
	}
	return nil, caos_errs.ThrowInternalf(err, "VIEW-9LyCB", "unable to get latest sequence of %s", viewName)
}

func AllCurrentSequences(db *gorm.DB, table string) ([]*CurrentSequence, error) {
	sequences := make([]*CurrentSequence, 0)
	query := PrepareSearchQuery(table, GeneralSearchRequest{})
	_, err := query(db, &sequences)
	if err != nil {
		return nil, err
	}
	return sequences, nil
}

func ClearView(db *gorm.DB, truncateView, sequenceTable string) error {
	truncate := PrepareTruncate(truncateView)
	err := truncate(db)
	if err != nil {
		return err
	}
	return SaveCurrentSequence(db, sequenceTable, truncateView, "", 0, time.Now())
}
