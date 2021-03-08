package management

import (
	"context"
	"time"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/authn"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	key_model "github.com/caos/zitadel/internal/key/model"
	user_model "github.com/caos/zitadel/internal/user/model"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
	user_pb "github.com/caos/zitadel/pkg/grpc/user"
)

func ListUsersRequestToModel(ctx context.Context, req *mgmt_pb.ListUsersRequest) *user_model.UserSearchRequest {
	req.Queries = append(req.Queries, &user_pb.SearchQuery{
		Query: &user_pb.SearchQuery_ResourceOwner{
			ResourceOwner: &user_pb.ResourceOwnerQuery{
				OrgID: authz.GetCtxData(ctx).OrgID,
			},
		},
	})

	return &user_model.UserSearchRequest{
		Offset:  req.Query.Offset,
		Limit:   uint64(req.Query.Limit),
		Asc:     req.Query.Asc,
		Queries: user_grpc.UserQueriesToModel(req.Queries),
	}
}

func AddHumanUserRequestToDomain(req *mgmt_pb.AddHumanUserRequest) *domain.Human {
	h := &domain.Human{
		Username: req.UserName,
	}
	preferredLanguage, err := language.Parse(req.Profile.PreferredLanguage)
	logging.Log("MANAG-3GUFJ").OnError(err).Debug("language malformed")
	h.Profile = &domain.Profile{
		FirstName:         req.Profile.FirstName,
		LastName:          req.Profile.LastName,
		NickName:          req.Profile.NickName,
		DisplayName:       req.Profile.DisplayName,
		PreferredLanguage: preferredLanguage,
		Gender:            user_grpc.GenderToDomain(req.Profile.Gender),
	}
	h.Email = &domain.Email{
		EmailAddress:    req.Email.Email,
		IsEmailVerified: req.Email.IsEmailVerified,
	}
	if req.Phone != nil {
		h.Phone = &domain.Phone{
			PhoneNumber:     req.Phone.Phone,
			IsPhoneVerified: req.Phone.IsPhoneVerified,
		}
	}
	if req.InitialPassword != "" {
		h.Password = &domain.Password{SecretString: req.InitialPassword}
	}

	return h
}

func AddMachineUserRequestToDomain(req *mgmt_pb.AddMachineUserRequest) *domain.Machine {
	return &domain.Machine{
		Username:    req.UserName,
		Name:        req.Name,
		Description: req.Description,
	}
}

func UpdateHumanProfileRequestToDomain(req *mgmt_pb.UpdateHumanProfileRequest) *domain.Profile {
	preferredLanguage, err := language.Parse(req.PreferredLanguage)
	logging.Log("MANAG-GPcYv").OnError(err).Debug("language malformed")
	return &domain.Profile{
		ObjectRoot:        models.ObjectRoot{AggregateID: req.UserId},
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		NickName:          req.NickName,
		DisplayName:       req.DisplayName,
		PreferredLanguage: preferredLanguage,
		Gender:            user_grpc.GenderToDomain(req.Gender),
	}
}

func UpdateHumanEmailRequestToDomain(req *mgmt_pb.UpdateHumanEmailRequest) *domain.Email {
	return &domain.Email{
		EmailAddress:    req.Email,
		IsEmailVerified: req.IsEmailVerified,
	}
}

func UpdateHumanPhoneRequestToDomain(req *mgmt_pb.UpdateHumanPhoneRequest) *domain.Phone {
	return &domain.Phone{
		PhoneNumber:     req.Phone,
		IsPhoneVerified: req.IsPhoneVerified,
	}
}

func notifyTypeToDomain(state mgmt_pb.SendHumanResetPasswordNotificationRequest_Type) domain.NotificationType {
	switch state {
	case mgmt_pb.SendHumanResetPasswordNotificationRequest_TYPE_EMAIL:
		return domain.NotificationTypeEmail
	case mgmt_pb.SendHumanResetPasswordNotificationRequest_TYPE_SMS:
		return domain.NotificationTypeSms
	default:
		return domain.NotificationTypeEmail
	}
}

func UpdateMachineRequestToDomain(ctx context.Context, req *mgmt_pb.UpdateMachineRequest) *domain.Machine {
	return &domain.Machine{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.UserId,
			ResourceOwner: authz.GetCtxData(ctx).OrgID,
		},
		Name:        req.Name,
		Description: req.Description,
	}
}

func ListMachineKeysRequestToModel(req *mgmt_pb.ListMachineKeysRequest) *key_model.AuthNKeySearchRequest {
	return &key_model.AuthNKeySearchRequest{
		Offset: req.Query.Offset,
		Limit:  uint64(req.Query.Limit),
		Asc:    req.Query.Asc,
		Queries: []*key_model.AuthNKeySearchQuery{
			{
				Key:    key_model.AuthNKeyObjectType,
				Method: domain.SearchMethodEquals,
				Value:  key_model.AuthNKeyObjectTypeUser,
			}, {
				Key:    key_model.AuthNKeyObjectID,
				Method: domain.SearchMethodEquals,
				Value:  req.UserId,
			},
		},
	}
}

func AddMachineKeyRequestToDomain(req *mgmt_pb.AddMachineKeyRequest) *domain.MachineKey {
	expDate := time.Time{}
	if req.ExpirationDate != nil {
		var err error
		expDate, err = ptypes.Timestamp(req.ExpirationDate)
		logging.Log("MANAG-iNshR").OnError(err).Debug("unable to parse expiration date")
	}

	return &domain.MachineKey{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.UserId,
		},
		ExpirationDate: expDate,
		Type:           authn.KeyTypeToDomain(req.Type),
	}
}

func RemoveHumanLinkedIDPRequestToDomain(ctx context.Context, req *mgmt_pb.RemoveHumanLinkedIDPRequest) *domain.ExternalIDP {
	return &domain.ExternalIDP{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.UserId,
			ResourceOwner: authz.GetCtxData(ctx).OrgID,
		},
		IDPConfigID:    req.IdpId,
		ExternalUserID: req.LinkedUserId,
	}
}

func ListHumanLinkedIDPsRequestToModel(req *mgmt_pb.ListHumanLinkedIDPsRequest) *user_model.ExternalIDPSearchRequest {
	return &user_model.ExternalIDPSearchRequest{
		Limit:   uint64(req.Query.Limit),
		Offset:  req.Query.Offset,
		Queries: []*user_model.ExternalIDPSearchQuery{{Key: user_model.ExternalIDPSearchKeyUserID, Method: domain.SearchMethodEquals, Value: req.UserId}},
	}
}

func ListUserMembershipsRequestToModel(req *mgmt_pb.ListUserMembershipsRequest) (*user_model.UserMembershipSearchRequest, error) {
	queries, err := user_grpc.MembershipQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	queries = append(queries, &user_model.UserMembershipSearchQuery{
		Key:    user_model.UserMembershipSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  req.UserId,
	})
	return &user_model.UserMembershipSearchRequest{
		Offset: req.Query.Offset,
		Limit:  uint64(req.Query.Limit),
		Asc:    req.Query.Asc,
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}
