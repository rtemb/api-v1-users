package api

import (
	"context"

	apiUsers "github.com/rtemb/api-v1-users/internal/proto/api-v1-users"
	"github.com/sirupsen/logrus"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 --fake-name HandlerServiceMock -o ../testing/mocks/handler_service.go . HandlerService
type HandlerService interface {
	CreateUser(ctx context.Context, req *apiUsers.CreateUserRequest) error
	Auth(ctx context.Context, req *apiUsers.AuthRequest) (*apiUsers.AuthResponse, error)
}

type Handler struct {
	service HandlerService
	logger  *logrus.Entry
}

func NewHandler(s HandlerService, l *logrus.Entry) *Handler {
	return &Handler{service: s, logger: l}
}
