package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/vaberof/effective-mobile-backend/internal/domain/car"
	"github.com/vaberof/effective-mobile-backend/pkg/logging/logs"
	"log/slog"
)

type Handler struct {
	carService car.CarService
	logger     *slog.Logger
}

func NewHandler(carService car.CarService, logsBuilder *logs.Logs) *Handler {
	logger := logsBuilder.WithName("handler")
	return &Handler{
		carService: carService,
		logger:     logger,
	}
}

func (h *Handler) InitRoutes(router *gin.Engine) *gin.Engine {
	router.Use(cors.Default())

	apiV1 := router.Group("/api/v1")

	// ====== Cars routes ======

	cars := apiV1.Group("/cars")
	cars.POST("/", h.CreateCars)
	cars.PATCH("/:id", h.UpdateCar)
	cars.DELETE("/:id", h.DeleteCar)
	cars.GET("/", h.ListCars)

	// ====== End of Cars routes ======

	// ====== Swagger route ======

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ====== End of Swagger route ======

	return router
}
