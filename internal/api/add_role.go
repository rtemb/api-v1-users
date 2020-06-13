package api

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	apiUsers "github.com/rtemb/api-v1-users/internal/proto/api-v1-users"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Handler) AddRole(ctx context.Context, req *apiUsers.AddRoleRequest) (*empty.Empty, error) {
	s.logger.WithFields(logrus.Fields{"method": "api.AddRole"}).Trace(req)
	if req.Role == apiUsers.Role_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "role is required")
	}

	err := s.service.AddRole(ctx, req)
	if err != nil {
		return nil, mapInternalToApiErrors(err)
	}

	return &empty.Empty{}, nil
}
