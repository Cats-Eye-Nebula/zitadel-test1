package view

import (
	"time"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/key/repository/view"
	"github.com/caos/zitadel/internal/key/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	keyTable = "auth.keys"
)

func (v *View) KeyByIDAndType(keyID string, private bool) (*model.KeyView, error) {
	return view.KeyByIDAndType(v.Db, keyTable, keyID, private)
}

func (v *View) GetSigningKey(expiry time.Time) (*key_model.SigningKey, time.Time, error) {
	key, err := view.GetSigningKey(v.Db, keyTable, expiry)
	if err != nil {
		return nil, time.Time{}, err
	}
	signingKey, err := key_model.SigningKeyFromKeyView(model.KeyViewToModel(key), v.keyAlgorithm)
	return signingKey, key.Expiry, err
}

func (v *View) GetActiveKeySet() ([]*key_model.PublicKey, error) {
	keys, err := view.GetActivePublicKeys(v.Db, keyTable)
	if err != nil {
		return nil, err
	}
	return key_model.PublicKeysFromKeyView(model.KeyViewsToModel(keys), v.keyAlgorithm)
}

func (v *View) PutKeys(privateKey, publicKey *model.KeyView, event *models.Event) error {
	err := view.PutKeys(v.Db, keyTable, privateKey, publicKey)
	if err != nil {
		return err
	}
	return v.ProcessedKeySequence(event)
}

func (v *View) DeleteKey(keyID string, private bool, event *models.Event) error {
	err := view.DeleteKey(v.Db, keyTable, keyID, private)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedKeySequence(event)
}

func (v *View) DeleteKeyPair(keyID string, event *models.Event) error {
	err := view.DeleteKeyPair(v.Db, keyTable, keyID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedKeySequence(event)
}

func (v *View) GetLatestKeySequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(keyTable)
}

func (v *View) ProcessedKeySequence(event *models.Event) error {
	return v.saveCurrentSequence(keyTable, event)
}

func (v *View) UpdateKeySpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(keyTable)
}

func (v *View) GetLatestKeyFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(keyTable, sequence)
}

func (v *View) ProcessedKeyFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
