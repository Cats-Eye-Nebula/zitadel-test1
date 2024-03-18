package execution

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
)

var _ execution.ExecutionServiceServer = (*Server)(nil)

type Server struct {
	execution.UnimplementedExecutionServiceServer
	command             *command.Commands
	query               *query.Queries
	ListActionFunctions func() []string
	ListGRPCMethods     func() []string
	ListGRPCServices    func() []string
}

type Config struct{}

func CreateServer(
	command *command.Commands,
	query *query.Queries,
	listActionFunctions func() []string,
	listGRPCMethods func() []string,
	listGRPCServices func() []string,
) *Server {
	return &Server{
		command:             command,
		query:               query,
		ListActionFunctions: listActionFunctions,
		ListGRPCMethods:     listGRPCMethods,
		ListGRPCServices:    listGRPCServices,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	execution.RegisterExecutionServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return execution.ExecutionService_ServiceDesc.ServiceName
}

func (s *Server) MethodPrefix() string {
	return execution.ExecutionService_ServiceDesc.ServiceName
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return execution.ExecutionService_AuthMethods
}

func (s *Server) RegisterGateway() server.RegisterGatewayFunc {
	return execution.RegisterExecutionServiceHandler
}

func checkExecutionEnabled(ctx context.Context) error {
	if authz.GetInstance(ctx).Features().Execution {
		return nil
	}
	return zerrors.ThrowPreconditionFailed(nil, "SCHEMA-141bwx3lef", "Errors.Execution.NotEnabled")
}
