package user

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/v2/internal/domain"
	user "github.com/zitadel/zitadel/v2/pkg/grpc/user/v2beta"
)

func (s *Server) RegisterTOTP(ctx context.Context, req *user.RegisterTOTPRequest) (*user.RegisterTOTPResponse, error) {
	return totpDetailsToPb(
		s.command.AddUserTOTP(ctx, req.GetUserId(), authz.GetCtxData(ctx).OrgID),
	)

}

func totpDetailsToPb(totp *domain.TOTP, err error) (*user.RegisterTOTPResponse, error) {
	if err != nil {
		return nil, err
	}
	return &user.RegisterTOTPResponse{
		Details: object.DomainToDetailsPb(totp.ObjectDetails),
		Uri:     totp.URI,
		Secret:  totp.Secret,
	}, nil
}

func (s *Server) VerifyTOTPRegistration(ctx context.Context, req *user.VerifyTOTPRegistrationRequest) (*user.VerifyTOTPRegistrationResponse, error) {
	objectDetails, err := s.command.CheckUserTOTP(ctx, req.GetUserId(), req.GetCode(), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &user.VerifyTOTPRegistrationResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}
