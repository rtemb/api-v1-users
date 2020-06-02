package api_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/rtemb/api-v1-users/internal/api"
	proto "gitlab.com/rtemb/api-v1-users/internal/proto/api-v1-users"
	"gitlab.com/rtemb/api-v1-users/internal/testing/mocks"
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
