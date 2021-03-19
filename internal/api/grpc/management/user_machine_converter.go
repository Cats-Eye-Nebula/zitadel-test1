package management

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
)

func machineCreateToModel(machine *management.CreateMachineRequest) *usr_model.Machine {
	return &usr_model.Machine{
		Name:        machine.Name,
		Description: machine.Description,
	}
}

func updateMachineToModel(machine *management.UpdateMachineRequest) *usr_model.Machine {
	return &usr_model.Machine{
		ObjectRoot:  models.ObjectRoot{AggregateID: machine.Id},
		Description: machine.Description,
	}
}

func machineFromModel(account *usr_model.Machine) *management.MachineResponse {
	return &management.MachineResponse{
		Name:        account.Name,
		Description: account.Description,
	}
}

func machineViewFromModel(machine *usr_model.MachineView) *management.MachineView {
	lastKeyAdded, err := ptypes.TimestampProto(machine.LastKeyAdded)
	logging.Log("MANAG-wGcAQ").OnError(err).Debug("unable to parse date")
	return &management.MachineView{
		Description:  machine.Description,
		Name:         machine.Name,
		LastKeyAdded: lastKeyAdded,
	}
}

func authnKeyViewsFromModel(keys ...*key_model.AuthNKeyView) []*management.MachineKeyView {
	keyViews := make([]*management.MachineKeyView, len(keys))
	for i, key := range keys {
		keyViews[i] = machineKeyViewFromModel(key)
	}
	return keyViews
}

func machineKeyViewFromModel(key *key_model.AuthNKeyView) *management.MachineKeyView {
	creationDate, err := ptypes.TimestampProto(key.CreationDate)
	logging.Log("MANAG-gluk7").OnError(err).Debug("unable to parse timestamp")

	expirationDate, err := ptypes.TimestampProto(key.ExpirationDate)
	logging.Log("MANAG-gluk7").OnError(err).Debug("unable to parse timestamp")

	return &management.MachineKeyView{
		Id:             key.ID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
		Sequence:       key.Sequence,
		Type:           machineKeyTypeFromModel(key.Type),
	}
}

func addMachineKeyToModel(key *management.AddMachineKeyRequest) *usr_model.MachineKey {
	expirationDate := time.Time{}
	if key.ExpirationDate != nil {
		var err error
		expirationDate, err = ptypes.Timestamp(key.ExpirationDate)
		logging.Log("MANAG-iNshR").OnError(err).Debug("unable to parse expiration date")
	}

	return &usr_model.MachineKey{
		ExpirationDate: expirationDate,
		Type:           machineKeyTypeToModel(key.Type),
		ObjectRoot:     models.ObjectRoot{AggregateID: key.UserId},
	}
}

func addMachineKeyFromModel(key *usr_model.MachineKey) *management.AddMachineKeyResponse {
	creationDate, err := ptypes.TimestampProto(key.CreationDate)
	logging.Log("MANAG-dlb8m").OnError(err).Debug("unable to parse cretaion date")

	expirationDate, err := ptypes.TimestampProto(key.ExpirationDate)
	logging.Log("MANAG-dlb8m").OnError(err).Debug("unable to parse cretaion date")

	detail, err := json.Marshal(struct {
		Type   string `json:"type"`
		KeyID  string `json:"keyId"`
		Key    string `json:"key"`
		UserID string `json:"userId"`
	}{
		Type:   "serviceaccount",
		KeyID:  key.KeyID,
		Key:    string(key.PrivateKey),
		UserID: key.AggregateID,
	})
	logging.Log("MANAG-lFQ2g").OnError(err).Warn("unable to marshall key")

	return &management.AddMachineKeyResponse{
		Id:             key.KeyID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
		Sequence:       key.Sequence,
		KeyDetails:     detail,
		Type:           machineKeyTypeFromModel(key.Type),
	}
}

func machineKeyTypeToModel(typ management.MachineKeyType) key_model.AuthNKeyType {
	switch typ {
	case management.MachineKeyType_MACHINEKEY_JSON:
		return key_model.AuthNKeyTypeJSON
	default:
		return key_model.AuthNKeyTypeNONE
	}
}

func machineKeyTypeFromModel(typ key_model.AuthNKeyType) management.MachineKeyType {
	switch typ {
	case key_model.AuthNKeyTypeJSON:
		return management.MachineKeyType_MACHINEKEY_JSON
	default:
		return management.MachineKeyType_MACHINEKEY_UNSPECIFIED
	}
}

func machineKeySearchRequestToModel(req *management.MachineKeySearchRequest) *key_model.AuthNKeySearchRequest {
	return &key_model.AuthNKeySearchRequest{
		Offset: req.Offset,
		Limit:  req.Limit,
		Asc:    req.Asc,
		Queries: []*key_model.AuthNKeySearchQuery{
			{
				Key:    key_model.AuthNKeyObjectType,
				Method: model.SearchMethodEquals,
				Value:  key_model.AuthNKeyObjectTypeUser,
			}, {
				Key:    key_model.AuthNKeyObjectID,
				Method: model.SearchMethodEquals,
				Value:  req.UserId,
			},
		},
	}
}

func machineKeySearchResponseFromModel(req *key_model.AuthNKeySearchResponse) *management.MachineKeySearchResponse {
	viewTimestamp, err := ptypes.TimestampProto(req.Timestamp)
	logging.Log("MANAG-Sk9ds").OnError(err).Debug("unable to parse cretaion date")

	return &management.MachineKeySearchResponse{
		Offset:            req.Offset,
		Limit:             req.Limit,
		TotalResult:       req.TotalResult,
		ProcessedSequence: req.Sequence,
		ViewTimestamp:     viewTimestamp,
		Result:            authnKeyViewsFromModel(req.Result...),
	}
}
