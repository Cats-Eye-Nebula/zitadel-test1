package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
	"github.com/caos/zitadel/v2/internal/api/grpc/object"
	"github.com/caos/zitadel/v2/internal/api/grpc/user"
)

func (s *Server) GetMyEmail(ctx context.Context, _ *auth_pb.GetMyEmailRequest) (*auth_pb.GetMyEmailResponse, error) {
	email, err := s.query.GetHumanEmail(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyEmailResponse{
		Email: user.ModelEmailToPb(email),
		Details: object.ToViewDetailsPb(
			email.Sequence,
			email.CreationDate,
			email.ChangeDate,
			email.ResourceOwner,
		),
	}, nil
}

func (s *Server) SetMyEmail(ctx context.Context, req *auth_pb.SetMyEmailRequest) (*auth_pb.SetMyEmailResponse, error) {
	email, err := s.command.ChangeHumanEmail(ctx, UpdateMyEmailToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &auth_pb.SetMyEmailResponse{
		Details: object.ChangeToDetailsPb(
			email.Sequence,
			email.ChangeDate,
			email.ResourceOwner,
		),
	}, nil
}

func (s *Server) VerifyMyEmail(ctx context.Context, req *auth_pb.VerifyMyEmailRequest) (*auth_pb.VerifyMyEmailResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.VerifyHumanEmail(ctx, ctxData.UserID, req.Code, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.VerifyMyEmailResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ResendMyEmailVerification(ctx context.Context, _ *auth_pb.ResendMyEmailVerificationRequest) (*auth_pb.ResendMyEmailVerificationResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.CreateHumanEmailVerificationCode(ctx, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ResendMyEmailVerificationResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
