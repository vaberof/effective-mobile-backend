package logging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vaberof/effective-mobile-backend/pkg/logging/logs"
	"log/slog"
)

type Middleware struct {
	Handler gin.HandlerFunc
	Logger  *slog.Logger
}

func New(logs *logs.Logs) *Middleware {
	return impl(logs, "")
}

func impl(logs *logs.Logs, serverName string) *Middleware {
	loggerName := "http-server"
	if serverName != "" {
		loggerName = fmt.Sprintf("%s.%s", loggerName, serverName)
	}
	logger := logs.WithName(loggerName)

	handler := func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "" {
			path = "/"
		}
		method := c.Request.Method
		logger.Info("Request started", slog.Group("http", "path", path, "method", method))

		defer func() {
			status := c.Writer.Status()

			if status >= 500 {
				logger.Info("Request finished", slog.Group("http", "path", path, "method", method, "result", "error", "status", status))
				return
			}

			logger.Info("Request finished", slog.Group("http", "path", path, "method", method, "result", "success", "status", status))
		}()

		c.Next()
	}

	return &Middleware{
		Handler: handler,
		Logger:  logger,
	}
}
