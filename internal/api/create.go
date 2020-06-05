package api

import (
	"context"

	"github.com/pkg/errors"
	apiUsers "github.com/rtemb/api-v1-users/internal/proto/api-v1-users"
	"github.com/sirupsen/logrus"
)

func (s *Handler) CreateUser(ctx context.Context, req *apiUsers.CreateUserRequest) (*apiUsers.CreateUserResponse, error) {
	s.logger.WithFields(logrus.Fields{"method": "api.CreateUser"}).Trace(req)

	r := &apiUsers.CreateUserResponse{}
	err := s.service.CreateUser(ctx, req)
	if err != nil {
		return r, errors.Wrap(err, "unable to crate user")
	}

	r.StateCode = apiUsers.ResponseStateCode_CREATED
	return r, nil
}
