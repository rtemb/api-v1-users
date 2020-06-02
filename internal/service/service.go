package service

import (
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger *logrus.Entry
}

func NewService(l *logrus.Entry) *Service {
	return &Service{logger: l}
}

func (s *Service) CreateUser() error {
	return nil
}

func (s *Service) Auth() error {
	return nil
}
