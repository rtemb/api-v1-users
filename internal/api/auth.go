package api

import (
	"context"

	"github.com/pkg/errors"
	apiUsers "github.com/rtemb/api-v1-users/internal/proto/api-v1-users"
	"github.com/sirupsen/logrus"
)

func (s *Handler) Auth(ctx context.Context, req *apiUsers.AuthRequest) (*apiUsers.AuthResponse, error) {
	s.logger.WithFields(logrus.Fields{"method": "api.Auth"}).Trace(req)

	rsp, err := s.service.Auth(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to auth user")
	}

	return rsp, nil
}
