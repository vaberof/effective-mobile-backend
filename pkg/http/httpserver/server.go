package httpserver

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vaberof/effective-mobile-backend/pkg/http/httpserver/middleware/logging"
	"github.com/vaberof/effective-mobile-backend/pkg/logging/logs"
	"log/slog"
	"net/http"
)

type AppServer struct {
	Router *gin.Engine
	config *ServerConfig
	logger *slog.Logger
	addr   string
}

func New(config *ServerConfig, logs *logs.Logs) *AppServer {
	loggingMw := logging.New(logs)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(loggingMw.Handler)

	return &AppServer{
		Router: router,
		config: config,
		logger: loggingMw.Logger,
		addr:   fmt.Sprintf("%s:%d", config.Host, config.Port),
	}
}

func (server *AppServer) StartAsync() <-chan error {
	server.logger.Info("Starting http server")

	exitChannel := make(chan error)

	go func() {
		err := http.ListenAndServe(server.addr, server.Router)
		if !errors.Is(err, http.ErrServerClosed) {
			server.logger.Error("Failed to start HTTP server", slog.Group("error", err))
			exitChannel <- err
			return
		} else {
			exitChannel <- nil
		}
	}()

	server.logger.Info("Started HTTP server", slog.Group("http-server", "address", server.addr))

	return exitChannel
}

func (server *AppServer) GetLogger() *slog.Logger {
	return server.logger
}
