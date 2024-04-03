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
	"strconv"
)

// @Summary		Delete a car
// @Tags			cars
// @Description	Delete a car
// @ID				delete-car
// @Accept			json
// @Produce		json
// @Param			id	path		integer	true	"Cars`s id that needs to be deleted"
// @Success		200	{object}	deleteCarResponseBody
// @Failure		400	{object}	apiv1.Response
// @Failure		404	{object}	apiv1.Response
// @Failure		500	{object}	apiv1.Response
// @Router			/cars/{id} [delete]
func (h *Handler) DeleteCar(ctx *gin.Context) {
	const operation = "DeleteCar"

	log := h.logger.With(
		slog.String("operation", operation),
	)

	id := ctx.Param("id")

	carId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeInternalError, err.Error()))

		log.Warn("failed to delete a car", "error", err)

		return
	}

	err = h.carService.Delete(carId)
	if err != nil {
		if errors.Is(err, car.ErrCarNotFound) {
			ctx.JSON(http.StatusNotFound, apiv1.Error(apiv1.CodeNotFound, err.Error()))

			log.Warn("failed to delete a car", "error", err)
		} else {
			ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))

			log.Error("failed to delete a car", "error", err)
		}

		return
	}

	payload, _ := json.Marshal(&deleteCarResponseBody{
		Message: fmt.Sprintf("car with id '%d' has deleted successfully", carId),
	})

	ctx.JSON(http.StatusOK, apiv1.Success(payload))
}
