package management

import (
	"encoding/json"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func machineCreateToDomain(machine *management.CreateMachineRequest) *domain.Machine {
	return &domain.Machine{
		Name:        machine.Name,
		Description: machine.Description,
	}
}

func updateMachineToDomain(ctxData authz.CtxData, machine *management.UpdateMachineRequest) *domain.Machine {
	return &domain.Machine{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   machine.Id,
			ResourceOwner: ctxData.ResourceOwner,
		},
		Name:        machine.Name,
		Description: machine.Description,
	}
}

func machineFromDomain(account *domain.Machine) *management.MachineResponse {
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

func addMachineKeyToDomain(key *management.AddMachineKeyRequest) *domain.MachineKey {
	expirationDate := time.Time{}
	if key.ExpirationDate != nil {
		var err error
		expirationDate, err = ptypes.Timestamp(key.ExpirationDate)
		logging.Log("MANAG-iNshR").OnError(err).Debug("unable to parse expiration date")
	}

	return &domain.MachineKey{
		ExpirationDate: expirationDate,
		Type:           machineKeyTypeToDomain(key.Type),
		ObjectRoot:     models.ObjectRoot{AggregateID: key.UserId},
	}
}

func addMachineKeyFromDomain(key *domain.MachineKey) *management.AddMachineKeyResponse {
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
		CreationDate:   timestamppb.New(key.CreationDate),
		ExpirationDate: timestamppb.New(key.ExpirationDate),
		Sequence:       key.Sequence,
		KeyDetails:     detail,
		Type:           machineKeyTypeFromDomain(key.Type),
	}
}

func machineKeyTypeToDomain(typ management.MachineKeyType) domain.AuthNKeyType {
	switch typ {
	case management.MachineKeyType_MACHINEKEY_JSON:
		return domain.AuthNKeyTypeJSON
	default:
		return domain.AuthNKeyTypeNONE
	}
}

func machineKeyTypeFromDomain(typ domain.AuthNKeyType) management.MachineKeyType {
	switch typ {
	case domain.AuthNKeyTypeJSON:
		return management.MachineKeyType_MACHINEKEY_JSON
	default:
		return management.MachineKeyType_MACHINEKEY_UNSPECIFIED
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
