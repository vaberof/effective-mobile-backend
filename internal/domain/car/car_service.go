package car

import (
	"errors"
	"fmt"
	"github.com/vaberof/effective-mobile-backend/internal/infra/storage"
	"github.com/vaberof/effective-mobile-backend/pkg/logging/logs"
	"log/slog"
)

var (
	ErrCarNotFound = errors.New("car not found")
)

type CarService interface {
	Create(regNums []string) error
	Update(id int64, regNum, mark, model *string, year *int16, ownerName, ownerSurname, ownerPatronymic *string) (*Car, error)
	Delete(id int64) error
	ListByFilters(regNum, mark, model *string, year *int16, limit, offset int) ([]*Car, error)
}

type carServiceImpl struct {
	carStorage    CarStorage
	carApiService CarApiService
	logger        *slog.Logger
}

func NewCarService(carStorage CarStorage, carApiService CarApiService, logsBuilder *logs.Logs) CarService {
	return &carServiceImpl{
		carStorage:    carStorage,
		carApiService: carApiService,
		logger:        logsBuilder.WithName("domain.service.car"),
	}
}

func (c *carServiceImpl) Create(regNums []string) error {
	const operation = "Create"

	log := c.logger.With(
		slog.String("operation", operation))

	log.Info("creating a car")

	// call external api to get car info
	for _, regNum := range regNums {
		carInfo, err := c.carApiService.GetCar(regNum)
		if err != nil {
			log.Error("failed to create a car", "error", err)

			return fmt.Errorf("failed to create car: %w", err)
		}

		err = c.carStorage.Create(regNum, carInfo.Mark, carInfo.Model, &carInfo.Year, carInfo.Owner)
		if err != nil {
			log.Error("failed to create a car", "error", err)

			return err
		}
	}

	log.Info("cars have created")

	return nil
}

func (c *carServiceImpl) Update(id int64, regNum, mark, model *string, year *int16, ownerName, ownerSurname, ownerPatronymic *string) (*Car, error) {
	const operation = "Update"

	log := c.logger.With(
		slog.String("operation", operation),
		slog.Int64("id", id))

	log.Info("updating a car")

	domainCar, err := c.carStorage.Update(id, regNum, mark, model, year, ownerName, ownerSurname, ownerPatronymic)
	if err != nil {
		if errors.Is(err, storage.ErrCarNotFound) {
			log.Warn("failed to update a car", "error", err)

			return nil, ErrCarNotFound
		}

		log.Error("failed to update a car", "error", err)

		return nil, err
	}

	log.Info("car has updated")

	return domainCar, nil
}

func (c *carServiceImpl) Delete(id int64) error {
	const operation = "Delete"

	log := c.logger.With(
		slog.String("operation", operation),
		slog.Int64("id", id))

	log.Info("deleting a car")

	err := c.carStorage.Delete(id)
	if err != nil {
		if errors.Is(err, storage.ErrCarNotFound) {
			log.Warn("failed to delete a car", "error", err)

			return ErrCarNotFound
		}

		log.Error("failed to delete a car", "error", err)

		return err
	}

	log.Info("car has deleted")

	return nil
}

func (c *carServiceImpl) ListByFilters(regNum, mark, model *string, year *int16, limit, offset int) ([]*Car, error) {
	const operation = "ListByFilters"

	log := c.logger.With(
		slog.String("operation", operation),
		slog.Int("limit", limit),
		slog.Int("offset", offset),
	)

	log.Info("listing cars")

	domainCars, err := c.carStorage.ListByFilters(regNum, mark, model, year, limit, offset)
	if err != nil {
		log.Error("failed to list cars", "error", err)
		return nil, err
	}

	log.Info("cars have listed")

	return domainCars, nil
}
