package api

import (
	"context"

	apiUsers "gitlab.com/rtemb/api-v1-users/internal/proto/api-v1-users"
)

func (s *Handler) CreateUser(ctx context.Context, req *apiUsers.CreateUserRequest) (*apiUsers.CreateUserResponse, error) {
	s.logger.Debugln("Create - 123")
	// TODO
	r := &apiUsers.CreateUserResponse{}

	r.StateCode = apiUsers.ResponseStateCode_CREATED
	return r, nil
}
