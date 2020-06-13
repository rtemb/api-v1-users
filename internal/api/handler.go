package api

import (
	"context"
	"errors"

	apiUsers "github.com/rtemb/api-v1-users/internal/proto/api-v1-users"
	srvErr "github.com/rtemb/srv-users/pkg/client/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 --fake-name HandlerServiceMock -o ../testing/mocks/handler_service.go . HandlerService
type HandlerService interface {
	CreateUser(ctx context.Context, req *apiUsers.CreateUserRequest) error
	Auth(ctx context.Context, req *apiUsers.AuthRequest) (*apiUsers.AuthResponse, error)
	AddRole(ctx context.Context, req *apiUsers.AddRoleRequest) error
}

type Handler struct {
	service HandlerService
	logger  *logrus.Entry
}

func NewHandler(s HandlerService, l *logrus.Entry) *Handler {
	return &Handler{service: s, logger: l}
}

func mapInternalToApiErrors(e error) error {
	if errors.Is(e, srvErr.UserNotFound) {
		return status.Error(codes.NotFound, e.Error())
	}

	return status.Error(codes.Internal, e.Error())

}
