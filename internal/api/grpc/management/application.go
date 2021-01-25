package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) SearchApplications(ctx context.Context, in *management.ApplicationSearchRequest) (*management.ApplicationSearchResponse, error) {
	response, err := s.project.SearchApplications(ctx, applicationSearchRequestsToModel(in))
	if err != nil {
		return nil, err
	}
	return applicationSearchResponseFromModel(response), nil
}

func (s *Server) ApplicationByID(ctx context.Context, in *management.ApplicationID) (*management.ApplicationView, error) {
	app, err := s.project.ApplicationByID(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return applicationViewFromModel(app), nil
}

func (s *Server) CreateOIDCApplication(ctx context.Context, in *management.OIDCApplicationCreate) (*management.Application, error) {
	app, err := s.command.AddOIDCApplication(ctx, oidcAppCreateToDomain(in), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return oidcAppFromDomain(app), nil
}
func (s *Server) UpdateApplication(ctx context.Context, in *management.ApplicationUpdate) (*management.Application, error) {
	app, err := s.command.ChangeApplication(ctx, appUpdateToDomain(in), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return appFromDomain(app), nil
}
func (s *Server) DeactivateApplication(ctx context.Context, in *management.ApplicationID) (*empty.Empty, error) {
	err := s.command.DeactivateApplication(ctx, in.ProjectId, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}
func (s *Server) ReactivateApplication(ctx context.Context, in *management.ApplicationID) (*empty.Empty, error) {
	err := s.command.ReactivateApplication(ctx, in.ProjectId, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) RemoveApplication(ctx context.Context, in *management.ApplicationID) (*empty.Empty, error) {
	err := s.command.RemoveApplication(ctx, in.ProjectId, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) UpdateApplicationOIDCConfig(ctx context.Context, in *management.OIDCConfigUpdate) (*management.OIDCConfig, error) {
	config, err := s.command.ChangeOIDCApplication(ctx, oidcConfigUpdateToDomain(in), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return oidcConfigFromDomain(config), nil
}

func (s *Server) RegenerateOIDCClientSecret(ctx context.Context, in *management.ApplicationID) (*management.ClientSecret, error) {
	config, err := s.command.ChangeOIDCApplicationSecret(ctx, in.ProjectId, in.Id, authz.GetCtxData(ctx).ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &management.ClientSecret{ClientSecret: config.ClientSecretString}, nil
}

func (s *Server) ApplicationChanges(ctx context.Context, changesRequest *management.ChangeRequest) (*management.Changes, error) {
	response, err := s.project.ApplicationChanges(ctx, changesRequest.Id, changesRequest.SecId, changesRequest.SequenceOffset, changesRequest.Limit, changesRequest.Asc)
	if err != nil {
		return nil, err
	}
	return appChangesToResponse(response, changesRequest.GetSequenceOffset(), changesRequest.GetLimit()), nil
}
