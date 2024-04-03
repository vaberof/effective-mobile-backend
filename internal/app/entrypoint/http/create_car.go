package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vaberof/effective-mobile-backend/internal/domain/car"
	"github.com/vaberof/effective-mobile-backend/pkg/http/protocols/apiv1"
	"log/slog"
	"net/http"
)

// @Summary		Create a new cars
// @Tags			cars
// @Description	Create a new cars
// @ID				create-cars
// @Accept			json
// @Produce		json
// @Param			input	body		createCarsRequestBody	true	"Payload with array of car registration numbers that needs to be created"
// @Success		200		{object}	createCarsResponseBody
// @Failure		400		{object}	apiv1.Response
// @Failure		404		{object}	apiv1.Response
// @Failure		500		{object}	apiv1.Response
// @Router			/cars [post]
func (h *Handler) CreateCars(ctx *gin.Context) {
	const operation = "CreateCars"

	log := h.logger.With(
		slog.String("operation", operation),
	)

	var createCarReqBody createCarsRequestBody

	if err := ctx.Bind(&createCarReqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, "invalid request body here"))

		log.Warn("failed to create a car", "error", err)

		return
	}

	err := h.carService.Create(createCarReqBody.RegNums)
	if err != nil {
		if errors.Is(err, car.ErrCarNotFound) {
			ctx.JSON(http.StatusNotFound, apiv1.Error(apiv1.CodeNotFound, err.Error()))

			log.Warn("failed to create a car", "error", err)
		} else {
			ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))

			log.Error("failed to create a car", "error", err)
		}

		return
	}

	payload, _ := json.Marshal(&createCarsResponseBody{
		Message: fmt.Sprintf("cars with reg nums '%v' have created successfully", createCarReqBody.RegNums),
	})

	ctx.JSON(http.StatusOK, apiv1.Success(payload))
}
