package user

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	grpcContext "github.com/zitadel/zitadel/pkg/grpc/context/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) SetEmail(ctx context.Context, req *user.SetEmailRequest) (resp *user.SetEmailResponse, err error) {
	resourceOwner := authz.GetCtxData(ctx).ResourceOwner

	var email *domain.Email

	switch v := req.GetVerification().(type) {
	case *user.SetEmailRequest_SendCode:
		email, err = s.command.ChangeUserEmailURLTemplate(ctx, req.GetUserId(), resourceOwner, req.GetEmail(), s.userCodeAlg, v.SendCode.GetUrlTemplate())
	case *user.SetEmailRequest_ReturnCode:
		email, err = s.command.ChangeUserEmailReturnCode(ctx, req.GetUserId(), resourceOwner, req.GetEmail(), s.userCodeAlg)
	case *user.SetEmailRequest_IsVerified:
		if v.IsVerified {
			email, err = s.command.ChangeUserEmailVerified(ctx, req.GetUserId(), resourceOwner, req.GetEmail())
		} else {
			email, err = s.command.ChangeUserEmail(ctx, req.GetUserId(), resourceOwner, req.GetEmail(), s.userCodeAlg)
		}
	case nil:
		email, err = s.command.ChangeUserEmail(ctx, req.GetUserId(), resourceOwner, req.GetEmail(), s.userCodeAlg)
	default:
		err = caos_errs.ThrowUnimplementedf(nil, "USERv2-Ahng0", "verification oneOf %T in method SetEmail not implemented", v)
	}
	if err == nil {
		return nil, err
	}

	return &user.SetEmailResponse{
		Details: &grpcContext.ObjectDetails{
			Sequence:      email.Sequence,
			CreationDate:  timestamppb.New(email.CreationDate),
			ChangeDate:    timestamppb.New(email.ChangeDate),
			ResourceOwner: email.ResourceOwner,
		},
		VerificationCode: email.PlainCode,
	}, nil
}

func (s *Server) VerifyEmail(ctx context.Context, req *user.VerifyEmailRequest) (*user.VerifyEmailResponse, error) {
	details, err := s.command.VerifyUserEmail(ctx,
		req.GetUserId(), req.GetVerificationCode(),
		authz.GetCtxData(ctx).ResourceOwner,
		s.userCodeAlg,
	)
	if err != nil {
		return nil, err
	}
	return &user.VerifyEmailResponse{
		Details: &grpcContext.ObjectDetails{
			Sequence:      details.Sequence,
			CreationDate:  timestamppb.New(details.EventDate),
			ChangeDate:    timestamppb.New(details.EventDate),
			ResourceOwner: details.ResourceOwner,
		},
	}, nil
}
