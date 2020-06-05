package api_test

import (
	"context"
	"testing"

	"github.com/rtemb/api-v1-users/internal/api"
	proto "github.com/rtemb/api-v1-users/internal/proto/api-v1-users"
	"github.com/rtemb/api-v1-users/internal/testing/mocks"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_AuthSuccess(t *testing.T) {
	req := &proto.AuthRequest{
		Email:    "test@example.com",
		Password: "1234",
	}

	l, _ := test.NewNullLogger()
	logger := logrus.NewEntry(l)

	service := &mocks.HandlerServiceMock{}
	service.AuthCalls(func(ctx context.Context, request *proto.AuthRequest) (response *proto.AuthResponse, err error) {
		assert.Equal(t, req.Email, request.Email)
		assert.Equal(t, req.Password, request.Password)
		rsp := &proto.AuthResponse{
			Token: "test-token",
			Valid: true,
		}

		return rsp, nil
	})

	h := api.NewHandler(service, logger)

	rsp, err := h.Auth(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, rsp)
	assert.Equal(t, "test-token", rsp.Token)
	assert.True(t, rsp.Valid)
}
