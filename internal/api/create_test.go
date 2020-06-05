package api_test

import (
	"context"
	"testing"

	"github.com/rtemb/api-v1-users/internal/api"
	apiUsers "github.com/rtemb/api-v1-users/internal/proto/api-v1-users"
	"github.com/rtemb/api-v1-users/internal/service"
	"github.com/rtemb/srv-users/pkg/client/srv-users/mocks"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateUser_Success(t *testing.T) {
	req := &apiUsers.CreateUserRequest{}
	l, _ := test.NewNullLogger()
	logger := logrus.NewEntry(l)

	srvUsersMock := &mocks.SrvUsersMock{}
	h := api.NewHandler(service.NewService(logger, srvUsersMock), logger)

	rsp, err := h.CreateUser(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, apiUsers.ResponseStateCode_CREATED, rsp.StateCode)
}
