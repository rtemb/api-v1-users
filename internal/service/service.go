package service

import (
	"context"

	"github.com/pkg/errors"
	apiUsers "github.com/rtemb/api-v1-users/internal/proto/api-v1-users"
	srvUsers "github.com/rtemb/srv-users/pkg/client/srv-users"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger         *logrus.Entry
	srvUsersClient srvUsers.UsersServiceClient
}

func NewService(l *logrus.Entry, client srvUsers.UsersServiceClient) *Service {
	return &Service{logger: l, srvUsersClient: client}
}

func (s *Service) CreateUser(ctx context.Context, req *apiUsers.CreateUserRequest) error {
	r := srvUsers.CreateUserRequest{
		Email:    req.Email,
		Password: req.Password,
		Company:  req.Company,
	}

	_, err := s.srvUsersClient.CreateUser(ctx, &r)
	if err != nil {
		return errors.Wrap(err, "unable to create user")
	}

	return nil
}

func (s *Service) Auth(ctx context.Context, req *apiUsers.AuthRequest) (*apiUsers.AuthResponse, error) {
	r := &srvUsers.AuthRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	rsp, err := s.srvUsersClient.Auth(ctx, r)
	if err != nil {
		return nil, errors.Wrap(err, "unable to auth user")
	}
	res := &apiUsers.AuthResponse{
		Token: rsp.Token,
		Valid: rsp.Valid,
	}

	return res, nil
}

func (s *Service) AddRole(ctx context.Context, req *apiUsers.AddRoleRequest) error {
	return nil
}
