package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/vaberof/effective-mobile-backend/pkg/http/protocols/apiv1"
	"log/slog"
	"net/http"
	"strconv"
)

const (
	defaultListCarsLimit  = 100
	defaultListCarsOffset = 0
)

// @Summary		List cars
// @Tags			cars
// @Description	List cars by filters
// @ID				list-cars
// @Accept			json
// @Produce		json
// @Param			limit	query		integer	false	"An optional query parameter 'limit' that limits total number of returned cars. By default 'limit' = 100"
// @Param			offset	query		integer	false	"An optional query parameter 'offset' that indicates how many records should be skipped while listing cars. By default 'offset' = 0"
// @Param			regNum	query		string	false	"An optional query parameter 'regNum'"
// @Param			mark	query		string	false	"An optional query parameter 'mark'"
// @Param			model	query		string	false	"An optional query parameter 'model'"
// @Param			year	query		int		false	"An optional query parameter 'year'"
// @Success		200		{object}	listCarsResponseBody
// @Failure		400		{object}	apiv1.Response
// @Failure		404		{object}	apiv1.Response
// @Failure		500		{object}	apiv1.Response
// @Router			/cars [get]
func (h *Handler) ListCars(ctx *gin.Context) {
	const operation = "ListCars"

	log := h.logger.With(
		slog.String("operation", operation),
	)

	var limit, offset int
	var err error

	limitStr := ctx.Query("limit")

	if limitStr == "" {
		limit = defaultListCarsLimit
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Error("failed to convert 'limit' parameter", "limit", limitStr, "error", err.Error())

			ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))

			return
		}
		if limit < 0 {
			ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, "'limit' must not be negative"))

			return
		}
	}

	offsetStr := ctx.Query("offset")

	if offsetStr == "" {
		offset = defaultListCarsOffset
	} else {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			log.Error("failed to convert 'offset' parameter", "offset", offsetStr, "error", err.Error())

			ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))

			return
		}
		if offset < 0 {
			ctx.JSON(http.StatusBadRequest, apiv1.Error(apiv1.CodeBadRequest, "'offset' must not be negative"))

			return
		}
	}

	regNumStr := ctx.Query("regNum")
	markStr := ctx.Query("mark")
	modelStr := ctx.Query("model")
	yearStr := ctx.Query("year")

	var regNum, mark, model *string
	var year *int16

	if regNumStr != "" {
		regNum = &regNumStr
	}
	if markStr != "" {
		mark = &markStr
	}
	if modelStr != "" {
		model = &modelStr
	}

	if yearStr != "" {
		convYear, err := strconv.ParseInt(yearStr, 10, 16)
		if err != nil {
			log.Error("failed to convert 'yearStr' parameter", "yearStr", yearStr, "error", err.Error())

			ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))

			return
		}

		convYearInt16 := int16(convYear)
		year = &convYearInt16
	}

	domainCars, err := h.carService.ListByFilters(regNum, mark, model, year, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiv1.Error(apiv1.CodeInternalError, err.Error()))

		log.Error("failed to list cars", "error", err)

		return
	}

	if len(domainCars) == 0 {
		ctx.JSON(http.StatusNotFound, apiv1.Error(apiv1.CodeNotFound, "cars not found"))

		log.Warn("failed to list cars", "error", "cars not found")

		return
	}

	payload, _ := json.Marshal(&listCarsResponseBody{
		Cars: buildCarPayloads(domainCars),
	})

	ctx.JSON(http.StatusOK, apiv1.Success(payload))
}
