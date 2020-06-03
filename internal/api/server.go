package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/rtemb/api-v1-users/internal/config"
	"github.com/rtemb/api-v1-users/internal/middleware"
	apiUsers "github.com/rtemb/api-v1-users/internal/proto/api-v1-users"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func StartService(cfg *config.Server, mw middleware.AppMiddleware, grpcHandler apiUsers.UsersAPIServiceServer, logger *logrus.Entry) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	terminate := make(chan struct{}, 2)
	go startGRPC(terminate, cfg, logger, grpcHandler)
	go startGateway(terminate, ctx, cfg, logger, mw)

	<-interrupt
	terminate <- struct{}{}
	terminate <- struct{}{}
	logger.Debug("Shutting down")

	time.Sleep(2 * time.Second)
	os.Exit(0)
}

func startGRPC(terminate chan struct{}, cfg *config.Server, logger *logrus.Entry, grpcHandler apiUsers.UsersAPIServiceServer) {
	grpcAddr := cfg.GRPCHost + ":" + cfg.GRPCPort
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Fatalln("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	apiUsers.RegisterUsersAPIServiceServer(grpcServer, grpcHandler)
	var closed bool
	go func() {
		<-terminate
		closed = true
		logger.Debug("Shutting down gRPC")
		grpcServer.GracefulStop()
	}()

	logger.Info("Serving gRPC on http://", grpcAddr)
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Error(errors.Wrap(err, "unable to start gRPC server"))
	}

	if !closed {
		close(terminate)
	}
}

func startGateway(terminate chan struct{}, ctx context.Context, cfg *config.Server, logger *logrus.Entry, mw middleware.AppMiddleware) {
	customMarshaller := &runtime.JSONPb{
		OrigName:     true,
		EmitDefaults: true,
	}
	muxOpt := runtime.WithMarshalerOption(runtime.MIMEWildcard, customMarshaller)
	gw := runtime.NewServeMux(muxOpt)

	opts := []grpc.DialOption{grpc.WithInsecure()}
	dialAddr := fmt.Sprintf("dns:///%s", net.JoinHostPort(cfg.GRPCHost, cfg.GRPCPort))
	err := apiUsers.RegisterUsersAPIServiceHandlerFromEndpoint(ctx, gw, dialAddr, opts)
	if err != nil {
		logger.Fatalln(err)
	}

	apiServer := &http.Server{
		Addr:         ":" + cfg.GatewayPort,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      mw.SetupMiddleware(gw),
	}

	var closed bool
	go func() {
		<-terminate
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)
		defer shutdownCancel()

		logger.Debug("Shutting gateway...")
		err = apiServer.Shutdown(shutdownCtx)
		if err != nil {
			logger.Warnln(errors.Wrap(err, "unable to shutdown gateway"))
		}
	}()

	logger.Info("Serving gateway on " + cfg.GatewayPort)
	err = apiServer.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			logger.Error(errors.Wrap(err, "unable to start gateway"))
		}
		logger.Debug(err)
	}

	if !closed {
		close(terminate)
	}
}
