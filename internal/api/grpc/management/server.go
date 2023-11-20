package management

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/v2/internal/api/assets"
	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/api/grpc/server"
	"github.com/zitadel/zitadel/v2/internal/command"
	"github.com/zitadel/zitadel/v2/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/v2/internal/crypto"
	"github.com/zitadel/zitadel/v2/internal/query"
	"github.com/zitadel/zitadel/v2/pkg/grpc/management"
)

const (
	mgmtName = "Management-API"
)

var _ management.ManagementServiceServer = (*Server)(nil)

type Server struct {
	management.UnimplementedManagementServiceServer
	command         *command.Commands
	query           *query.Queries
	systemDefaults  systemdefaults.SystemDefaults
	assetAPIPrefix  func(context.Context) string
	passwordHashAlg crypto.HashAlgorithm
	userCodeAlg     crypto.EncryptionAlgorithm
	externalSecure  bool
}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	sd systemdefaults.SystemDefaults,
	userCodeAlg crypto.EncryptionAlgorithm,
	externalSecure bool,
) *Server {
	return &Server{
		command:         command,
		query:           query,
		systemDefaults:  sd,
		assetAPIPrefix:  assets.AssetAPI(externalSecure),
		passwordHashAlg: crypto.NewBCrypt(sd.SecretGenerators.PasswordSaltCost),
		userCodeAlg:     userCodeAlg,
		externalSecure:  externalSecure,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	management.RegisterManagementServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return mgmtName
}

func (s *Server) MethodPrefix() string {
	return management.ManagementService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return management.ManagementService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return management.RegisterManagementServiceHandler
}

func (s *Server) GatewayPathPrefix() string {
	return GatewayPathPrefix()
}

func GatewayPathPrefix() string {
	return "/management/v1"
}
