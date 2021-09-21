package admin

import (
	"context"

	idp_grpc "github.com/caos/zitadel/internal/api/grpc/idp"
	object_pb "github.com/caos/zitadel/internal/api/grpc/object"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetIDPByID(ctx context.Context, req *admin_pb.GetIDPByIDRequest) (*admin_pb.GetIDPByIDResponse, error) {
	idp, err := s.query.IDPByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetIDPByIDResponse{Idp: idp_grpc.IDPViewToPb(idp)}, nil
}

func (s *Server) ListIDPs(ctx context.Context, req *admin_pb.ListIDPsRequest) (*admin_pb.ListIDPsResponse, error) {
	queries, err := listIDPsToModel(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchIDPs(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListIDPsResponse{
		Result:  idp_grpc.IDPViewsToPb(resp.IDPs),
		Details: object_pb.ToListDetails(resp.Count, resp.Sequence, resp.Timestamp),
	}, nil
}

func (s *Server) AddOIDCIDP(ctx context.Context, req *admin_pb.AddOIDCIDPRequest) (*admin_pb.AddOIDCIDPResponse, error) {
	config, err := s.command.AddDefaultIDPConfig(ctx, addOIDCIDPRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddOIDCIDPResponse{
		IdpId: config.IDPConfigID,
		Details: object_pb.AddToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) AddJWTIDP(ctx context.Context, req *admin_pb.AddJWTIDPRequest) (*admin_pb.AddJWTIDPResponse, error) {
	config, err := s.command.AddDefaultIDPConfig(ctx, addJWTIDPRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddJWTIDPResponse{
		IdpId: config.IDPConfigID,
		Details: object_pb.AddToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateIDP(ctx context.Context, req *admin_pb.UpdateIDPRequest) (*admin_pb.UpdateIDPResponse, error) {
	config, err := s.command.ChangeDefaultIDPConfig(ctx, updateIDPToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateIDPResponse{
		Details: object_pb.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeactivateIDP(ctx context.Context, req *admin_pb.DeactivateIDPRequest) (*admin_pb.DeactivateIDPResponse, error) {
	objectDetails, err := s.command.DeactivateDefaultIDPConfig(ctx, req.IdpId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.DeactivateIDPResponse{Details: object_pb.DomainToChangeDetailsPb(objectDetails)}, nil
}

func (s *Server) ReactivateIDP(ctx context.Context, req *admin_pb.ReactivateIDPRequest) (*admin_pb.ReactivateIDPResponse, error) {
	objectDetails, err := s.command.ReactivateDefaultIDPConfig(ctx, req.IdpId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ReactivateIDPResponse{Details: object_pb.DomainToChangeDetailsPb(objectDetails)}, nil
}

func (s *Server) RemoveIDP(ctx context.Context, req *admin_pb.RemoveIDPRequest) (*admin_pb.RemoveIDPResponse, error) {
	idpProviders, err := s.iam.IDPProvidersByIDPConfigID(ctx, req.IdpId)
	if err != nil {
		return nil, err
	}
	externalIDPs, err := s.iam.ExternalIDPsByIDPConfigID(ctx, req.IdpId)
	if err != nil {
		return nil, err
	}
	objectDetails, err := s.command.RemoveDefaultIDPConfig(ctx, req.IdpId, idpProviderViewsToDomain(idpProviders), externalIDPViewsToDomain(externalIDPs)...)
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveIDPResponse{Details: object_pb.DomainToChangeDetailsPb(objectDetails)}, nil
}

func (s *Server) UpdateIDPOIDCConfig(ctx context.Context, req *admin_pb.UpdateIDPOIDCConfigRequest) (*admin_pb.UpdateIDPOIDCConfigResponse, error) {
	config, err := s.command.ChangeDefaultIDPOIDCConfig(ctx, updateOIDCConfigToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateIDPOIDCConfigResponse{
		Details: object_pb.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateIDPJWTConfig(ctx context.Context, req *admin_pb.UpdateIDPJWTConfigRequest) (*admin_pb.UpdateIDPJWTConfigResponse, error) {
	config, err := s.command.ChangeDefaultIDPJWTConfig(ctx, updateJWTConfigToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateIDPJWTConfigResponse{
		Details: object_pb.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}
