package api

import (
	"context"

	apiUsers "gitlab.com/rtemb/api-v1-users/internal/proto/api-v1-users"
)

func (s *Handler) Auth(ctx context.Context, req *apiUsers.AuthRequest) (*apiUsers.AuthResponse, error) {
	s.logger.Debugln("Auth - 123")
	r := &apiUsers.AuthResponse{}
	r.Token = "test-token"
	r.Valid = true

	return r, nil
}
