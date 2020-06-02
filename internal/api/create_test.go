package api_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/rtemb/api-v1-users/internal/api"
	apiUsers "gitlab.com/rtemb/api-v1-users/internal/proto/api-v1-users"
	"gitlab.com/rtemb/api-v1-users/internal/service"
)

func TestHandler_CreateUser_Success(t *testing.T) {
	req := &apiUsers.CreateUserRequest{}
	l, _ := test.NewNullLogger()
	logger := logrus.NewEntry(l)
	h := api.NewHandler(service.NewService(logger), logger)

	rsp, err := h.CreateUser(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, apiUsers.ResponseStateCode_CREATED, rsp.StateCode)
}
