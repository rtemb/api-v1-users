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
	req := &proto.AuthRequest{}
	l, _ := test.NewNullLogger()
	logger := logrus.NewEntry(l)
	service := &mocks.HandlerServiceMock{}
	h := api.NewHandler(service, logger)

	rsp, err := h.Auth(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "test-token", rsp.Token)
}
