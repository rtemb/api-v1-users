package api

import (
	"github.com/sirupsen/logrus"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 --fake-name HandlerServiceMock -o ../testing/mocks/handler_service.go . HandlerService
type HandlerService interface {
	CreateUser() error
	Auth() error
}

type Handler struct {
	service HandlerService
	logger  *logrus.Entry
}

func NewHandler(s HandlerService, l *logrus.Entry) *Handler {
	return &Handler{service: s, logger: l}
}
