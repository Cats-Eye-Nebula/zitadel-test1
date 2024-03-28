package middleware

import (
	"context"
	"encoding/json"
	"strings"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/execution"
	"github.com/zitadel/zitadel/internal/query"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func ExecutionHandler(queries *query.Queries) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		request, err := executeTargetsForRequest(ctx, queries, info.FullMethod, req)
		if err != nil {
			return nil, err
		}

		resp, err := handler(ctx, request)
		if err != nil {
			return nil, err
		}

		return executeTargetsForResponse(ctx, queries, info.FullMethod, req, resp)
	}
}

func executeTargetsForRequest(ctx context.Context, queries ExecutionQueries, fullMethod string, req interface{}) (interface{}, error) {
	targets, err := queryTargets(ctx, queries, fullMethod, domain.ExecutionTypeRequest)
	if err != nil {
		// if no targets are found, return without any calls
		if zerrors.IsNotFound(err) {
			return req, nil
		}
		return nil, err
	}

	ctxData := authz.GetCtxData(ctx)
	info := &ContextInfoRequest{
		FullMethod: fullMethod,
		InstanceID: authz.GetInstance(ctx).InstanceID(),
		ProjectID:  ctxData.ProjectID,
		OrgID:      ctxData.OrgID,
		UserID:     ctxData.UserID,
		Request:    req,
	}

	request, err := execution.CallTargets(ctx, targets, info)
	if err != nil {
		// if an error is returned still return also the original request
		return req, err
	}
	return request, err
}

func executeTargetsForResponse(ctx context.Context, queries ExecutionQueries, fullMethod string, req, resp interface{}) (interface{}, error) {
	targets, err := queryTargets(ctx, queries, fullMethod, domain.ExecutionTypeResponse)
	if err != nil {
		// if no targets are found, return without any calls
		if zerrors.IsNotFound(err) {
			return resp, nil
		}
		return nil, err
	}

	ctxData := authz.GetCtxData(ctx)
	info := &ContextInfoResponse{
		FullMethod: fullMethod,
		InstanceID: authz.GetInstance(ctx).InstanceID(),
		ProjectID:  ctxData.ProjectID,
		OrgID:      ctxData.OrgID,
		UserID:     ctxData.UserID,
		Request:    req,
		Response:   resp,
	}

	response, err := execution.CallTargets(ctx, targets, info)
	if err != nil {
		// if an error is returned still return also the original response
		return resp, err
	}
	return response, err
}

type ExecutionQueries interface {
	ExecutionTargetsRequestResponse(ctx context.Context, fullMethod, service, all string) (execution *query.ExecutionTargets, err error)
	SearchTargets(ctx context.Context, queries *query.TargetSearchQueries) (targets *query.Targets, err error)
}

func queryTargets(
	ctx context.Context,
	queries ExecutionQueries,
	fullMethod string,
	executionType domain.ExecutionType,
) ([]*query.Target, error) {
	exectargets, err := queries.ExecutionTargetsRequestResponse(ctx, exec_repo.ID(executionType, fullMethod), exec_repo.ID(executionType, serviceFromFullMethod(fullMethod)), exec_repo.IDAll(executionType))
	if err != nil {
		return nil, err
	}
	if exectargets == nil || len(exectargets.Targets) == 0 {
		return nil, zerrors.ThrowNotFound(err, "EXEC-m70fpc7a9q", "Errors.Execution.NotFound")
	}

	targetIDsQuery, err := query.NewTargetInIDsSearchQuery(exectargets.Targets)
	if err != nil {
		return nil, err
	}

	targets, err := queries.SearchTargets(ctx, &query.TargetSearchQueries{Queries: []query.SearchQuery{targetIDsQuery}})
	if err != nil {
		return nil, err
	}
	if targets == nil || len(targets.Targets) == 0 {
		return nil, zerrors.ThrowNotFound(err, "EXEC-x2r3cnfadi", "Errors.Execution.NotFound")
	}
	return targets.Targets, nil
}

func serviceFromFullMethod(s string) string {
	parts := strings.Split(s, "/")
	return parts[1]
}

var _ execution.ContextInfo = &ContextInfoRequest{}

type ContextInfoRequest struct {
	FullMethod string      `json:"fullMethod,omitempty"`
	InstanceID string      `json:"instanceID,omitempty"`
	OrgID      string      `json:"orgID,omitempty"`
	ProjectID  string      `json:"projectID,omitempty"`
	UserID     string      `json:"userID,omitempty"`
	Request    interface{} `json:"request,omitempty"`
}

func (c *ContextInfoRequest) GetHTTPRequestBody() []byte {
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return data
}

func (c *ContextInfoRequest) SetHTTPResponseBody(resp []byte) error {
	return json.Unmarshal(resp, c.Request)
}

func (c *ContextInfoRequest) GetContent() interface{} {
	return c.Request
}

var _ execution.ContextInfo = &ContextInfoResponse{}

type ContextInfoResponse struct {
	FullMethod string      `json:"fullMethod,omitempty"`
	InstanceID string      `json:"instanceID,omitempty"`
	OrgID      string      `json:"orgID,omitempty"`
	ProjectID  string      `json:"projectID,omitempty"`
	UserID     string      `json:"userID,omitempty"`
	Request    interface{} `json:"request,omitempty"`
	Response   interface{} `json:"response,omitempty"`
}

func (c *ContextInfoResponse) GetHTTPRequestBody() []byte {
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return data
}

func (c *ContextInfoResponse) SetHTTPResponseBody(resp []byte) error {
	return json.Unmarshal(resp, c.Response)
}

func (c *ContextInfoResponse) GetContent() interface{} {
	return c.Response
}
