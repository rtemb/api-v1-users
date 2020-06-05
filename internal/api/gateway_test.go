package api_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/rtemb/api-v1-users/internal/api"
	"github.com/rtemb/api-v1-users/internal/auth"
	"github.com/rtemb/api-v1-users/internal/config"
	"github.com/rtemb/api-v1-users/internal/middleware"
	apiUsers "github.com/rtemb/api-v1-users/internal/proto/api-v1-users"
	"github.com/rtemb/api-v1-users/internal/testing/mocks"
	"github.com/rtemb/api-v1-users/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// nolint
type APITestSuite struct {
	suite.Suite
	handlerServiceMock *mocks.HandlerServiceMock
	Logger             *logrus.Entry
}

func (a *APITestSuite) SetupSuite() {
	a.handlerServiceMock = &mocks.HandlerServiceMock{}
	a.handlerServiceMock.AuthCalls(func(ctx context.Context, request *apiUsers.AuthRequest) (response *apiUsers.AuthResponse, err error) {
		a.Equal("test@example.com", request.Email)
		a.Equal("qwerty", request.Password)
		rsp := &apiUsers.AuthResponse{
			Token: "test-token",
			Valid: true,
		}

		return rsp, nil
	})
	a.handlerServiceMock.CreateUserCalls(func(ctx context.Context, request *apiUsers.CreateUserRequest) error {
		a.Equal("test@example.com", request.Email)
		a.Equal("qwerty", request.Password)
		a.Equal("test", request.Company)
		return nil
	})
	go a.initServerWithGateway()

}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, &APITestSuite{})
}

func (a *APITestSuite) Test_CreateUser() {
	req := `{"email": "test@example.com", "password": "qwerty", "company":"test"}`
	rsp, err := http.Post(
		"http://localhost:8081/v1/users/user",
		"application/json",
		strings.NewReader(req))
	a.Require().NoError(err)

	body, err := ioutil.ReadAll(rsp.Body)
	a.Require().NoError(err)
	s := string(body)

	a.Equal(http.StatusOK, rsp.StatusCode)
	a.Equal(`{"stateCode":"CREATED"}`, s)
}

func (a *APITestSuite) Test_Auth() {
	req := `{"email": "test@example.com", "password": "qwerty"}`
	rsp, err := http.Post(
		"http://localhost:8081/v1/users/auth",
		"application/json",
		strings.NewReader(req))
	a.Require().NoError(err)

	body, err := ioutil.ReadAll(rsp.Body)
	a.Require().NoError(err)
	s := string(body)

	a.Equal(http.StatusOK, rsp.StatusCode)
	a.Equal(`{"token":"test-token","valid":true,"errors":[]}`, s)
}

func (a APITestSuite) initServerWithGateway() {
	a.Logger = logrus.New().WithFields(logrus.Fields{
		"gitSha":  version.GitSha,
		"version": version.ServiceVersion,
		"logger":  "cmd/api-v1-users",
	})
	lvl, err := logrus.ParseLevel("trace")
	a.Require().NoError(err)
	a.Logger.Logger.SetLevel(lvl)

	err = os.Setenv("ADMIN_HASH", "password")
	a.Require().NoError(err)
	cfg, err := config.Load()
	a.Require().NoError(err)

	simpleAuth := auth.NewSimpleAuth("1234")
	mw := middleware.NewMiddleware(simpleAuth, a.Logger)

	apiHandler := api.NewHandler(a.handlerServiceMock, a.Logger)
	api.StartService(cfg.Server, mw, apiHandler, a.Logger)
}
