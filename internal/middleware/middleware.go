package middleware

import (
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 --fake-name AuthCheckerMock -o ../testing/mocks/auth_checker.go . AuthChecker
type AuthChecker interface {
	CheckAccess(key string) (bool, error)
}

type RequestMiddleware struct {
	logger      *logrus.Entry
	authChecker AuthChecker
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 --fake-name AppMiddlewareMock -o ../testing/mocks/app_middleware.go . AppMiddleware
type AppMiddleware interface {
	SetupMiddleware(gw *runtime.ServeMux) *http.ServeMux
	LogRequests(next http.Handler) http.Handler
	Authorization(next http.Handler) http.Handler
	CORSProtection(next http.Handler) http.Handler
	ReadinesslivenessProbe(next http.Handler) http.Handler
}

const (
	URI_ROOT   = "/"
	API_PREFIX = "/v1/users"
)

const (
	EP_LIVENESS  = "/liveness"
	EP_READUNESS = "/readiness"
)

func NewMiddleware(ac AuthChecker, l *logrus.Entry) *RequestMiddleware {
	return &RequestMiddleware{authChecker: ac, logger: l}
}

func (m *RequestMiddleware) SetupMiddleware(gw *runtime.ServeMux) *http.ServeMux {
	muxRouter := http.NewServeMux()

	muxRouter.Handle(URI_ROOT, m.LogRequests(gw))
	//muxRouter.Handle(API_PREFIX+PREFIX_ADMIN, m.LogRequests(m.Authorization(gw)))
	muxRouter.Handle(API_PREFIX+EP_LIVENESS, m.LogRequests(m.ReadinesslivenessProbe(gw)))
	muxRouter.Handle(API_PREFIX+EP_READUNESS, m.LogRequests(m.ReadinesslivenessProbe(gw)))

	return muxRouter
}

// LogRequests print to log URL of each request
func (m *RequestMiddleware) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.WithFields(logrus.Fields{
			"method":    "RequestMiddleware.LogRequests",
			"timestamp": time.Now().Format(time.RFC3339),
		}).Trace(r.Method + " " + r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

// CORSProtection performs CORS protection
func (m *RequestMiddleware) CORSProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		next.ServeHTTP(w, r)
	})
}

// Authorization auth for admin section
func (m *RequestMiddleware) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := m.authChecker.CheckAccess(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		l := m.logger.WithFields(logrus.Fields{
			"method":    "RequestMiddleware.Authorization",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		if res {
			l.Debugln("authorized access")
			next.ServeHTTP(w, r)
		} else {
			l.Debugln("unauthorized access")
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		}
	})
}

func (m *RequestMiddleware) ReadinesslivenessProbe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
}
