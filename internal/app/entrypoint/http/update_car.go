package http

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vaberof/effective-mobile-backend/internal/domain/car"
	"github.com/vaberof/effective-mobile-backend/pkg/http/protocols/apiv1"
	"log/slog"
	"net/http"
	"strconv"
)

// @Summary		Update a car
// @Tags			cars
// @Description	Update a car
// @ID				update-cars
// @Accept			json
// @Produce		json
// @Param			id		path		int						true	"Car`s id that needs to be updated"
// @Param			input	body		updateCarRequestBody	true	"Car object that needs to be updated"
// @Success		200		{object}	updateCarResponseBody
// @Failure		400		{object}	apiv1.Response
// @Failure		404		{object}	apiv1.Response
// @Failure		500		{object}	apiv1.Response
// @Router			/cars/{id} [patch]
func (h *Handler) UpdateCar(ctx *gin.Context) {
	const operation = "UpdateCar"

	log := h.logger.With(
		slog.String("operation", operation),
	)

	id := ctx.Param("id")

	carId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeInternalError, err.Error()))

		log.Warn("failed to update a car", "error", err)

		return
	}

	var updateCarReqBody updateCarRequestBody

	if err := ctx.Bind(&updateCarReqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, "invalid request body here"))

		log.Warn("failed to update a car", "error", err)

		return
	}

	domainCar, err := h.carService.Update(
		carId,
		updateCarReqBody.RegNum,
		updateCarReqBody.Mark,
		updateCarReqBody.Model,
		updateCarReqBody.Year,
		updateCarReqBody.Owner.Name,
		updateCarReqBody.Owner.Surname,
		updateCarReqBody.Owner.Patronymic,
	)

	if err != nil {
		if errors.Is(err, car.ErrCarNotFound) {
			ctx.JSON(http.StatusNotFound, apiv1.Error(apiv1.CodeNotFound, err.Error()))

			log.Warn("failed to update a car", "error", err)
		} else {
			ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))

			log.Error("failed to update a car", "error", err)
		}

		return
	}

	resp := &updateCarResponseBody{
		Id:     domainCar.Id,
		RegNum: domainCar.RegistrationNumber,
		Mark:   domainCar.Mark,
		Model:  domainCar.Model,
		Year:   domainCar.Year,
	}
	if domainCar.Owner != nil {
		resp.Owner = &carOwnerPayload{
			Name:       domainCar.Owner.Name,
			Surname:    domainCar.Owner.Surname,
			Patronymic: domainCar.Owner.Patronymic,
		}
	}
	payload, _ := json.Marshal(resp)

	ctx.JSON(http.StatusOK, apiv1.Success(payload))
}
