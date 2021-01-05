package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	machineKeyTable = "auth.machine_keys"
)

func (v *View) MachineKeyByIDs(userID, keyID string) (*model.MachineKeyView, error) {
	return view.MachineKeyByIDs(v.Db, machineKeyTable, userID, keyID)
}

func (v *View) MachineKeysByUserID(userID string) ([]*model.MachineKeyView, error) {
	return view.MachineKeysByUserID(v.Db, machineKeyTable, userID)
}

func (v *View) MachineKeyByID(keyID string) (*model.MachineKeyView, error) {
	return view.MachineKeyByID(v.Db, machineKeyTable, keyID)
}

func (v *View) SearchMachineKeys(request *usr_model.MachineKeySearchRequest) ([]*model.MachineKeyView, uint64, error) {
	return view.SearchMachineKeys(v.Db, machineKeyTable, request)
}

func (v *View) PutMachineKey(key *model.MachineKeyView, event *models.Event) error {
	err := view.PutMachineKey(v.Db, machineKeyTable, key)
	if err != nil {
		return err
	}
	return v.ProcessedMachineKeySequence(event)
}

func (v *View) DeleteMachineKey(keyID string, event *models.Event) error {
	err := view.DeleteMachineKey(v.Db, machineKeyTable, keyID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedMachineKeySequence(event)
}

func (v *View) DeleteMachineKeysByUserID(userID string, event *models.Event) error {
	err := view.DeleteMachineKey(v.Db, machineKeyTable, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedMachineKeySequence(event)
}

func (v *View) GetLatestMachineKeySequence(aggregateType string) (*repository.CurrentSequence, error) {
	return v.latestSequence(machineKeyTable, aggregateType)
}

func (v *View) ProcessedMachineKeySequence(event *models.Event) error {
	return v.saveCurrentSequence(machineKeyTable, event)
}

func (v *View) UpdateMachineKeySpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(machineKeyTable)
}

func (v *View) GetLatestMachineKeyFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(machineKeyTable, sequence)
}

func (v *View) ProcessedMachineKeyFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
