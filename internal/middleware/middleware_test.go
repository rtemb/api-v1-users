package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	mw "github.com/rtemb/api-v1-users/internal/middleware"
	"github.com/rtemb/api-v1-users/internal/testing/mocks"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestMiddleware_LogRequests(t *testing.T) {
	handler, hook, _, err := initHandler()
	require.NoError(t, err)

	req, err := http.NewRequest("GET", "/v1/users/", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected := `level=trace msg="GET /v1/users/" method=RequestMiddleware.LogRequests timestamp=`
	actual, err := hook.LastEntry().String()
	require.NoError(t, err)

	assert.Contains(t, actual, expected)
}

func TestRequestMiddleware_ReadinessProbe(t *testing.T) {
	handler, hook, _, err := initHandler()
	require.NoError(t, err)

	req, err := http.NewRequest("GET", "/v1/users/readiness", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected := `level=trace msg="GET /v1/users/readiness" method=RequestMiddleware.LogRequests timestamp=`
	actual, err := hook.LastEntry().String()
	require.NoError(t, err)

	assert.Contains(t, actual, expected)
}

func TestRequestMiddleware_LivenessProbe(t *testing.T) {
	handler, hook, _, err := initHandler()
	require.NoError(t, err)

	req, err := http.NewRequest("GET", "/v1/users/liveness", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected := `level=trace msg="GET /v1/users/liveness" method=RequestMiddleware.LogRequests timestamp=`
	actual, err := hook.LastEntry().String()
	require.NoError(t, err)

	assert.Contains(t, actual, expected)
}

func initHandler() (*http.ServeMux, *test.Hook, *mocks.AuthCheckerMock, error) {
	authMock := &mocks.AuthCheckerMock{}

	logger, hook := test.NewNullLogger()
	l := logrus.NewEntry(logger)
	lvl, err := logrus.ParseLevel("trace")
	if err != nil {
		return nil, nil, nil, err
	}

	l.Logger.SetLevel(lvl)
	l.Logger.Hooks.Add(hook)

	reqMW := mw.NewMiddleware(authMock, l)
	mux := runtime.NewServeMux()
	handlerToTest := reqMW.SetupMiddleware(mux)

	return handlerToTest, hook, authMock, nil
}
