package api

import (
	"context"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	"github.com/caos/zitadel/internal/config/systemdefaults"

	"github.com/caos/zitadel/internal/admin/repository"
	"github.com/caos/zitadel/internal/api/auth"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/pkg/admin/api/grpc"
)

type Config struct {
	GRPC grpc_util.Config
}

func Start(ctx context.Context, conf Config, authZRepo *authz_repo.EsRepository, authZ auth.Config, defaults systemdefaults.SystemDefaults, repo repository.Repository) {
	grpcServer := grpc.StartServer(conf.GRPC.ToServerConfig(), authZRepo, authZ, repo)
	grpcGateway := grpc.StartGateway(conf.GRPC.ToGatewayConfig())

	server.StartServer(ctx, grpcServer, defaults)
	server.StartGateway(ctx, grpcGateway)
}
