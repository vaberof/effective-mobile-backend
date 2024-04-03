package apicar

import (
	"fmt"
	"github.com/vaberof/effective-mobile-backend/internal/domain/car"
	"github.com/vaberof/effective-mobile-backend/pkg/integration/api/apicar"
)

type CarApiService interface {
	GetCar(regNum string) (*car.Car, error)
}

type carApiServiceImpl struct {
	carApi CarApi
}

func NewCarApiService(carApi CarApi) CarApiService {
	return &carApiServiceImpl{carApi: carApi}
}

func (c *carApiServiceImpl) GetCar(regNum string) (*car.Car, error) {
	getCarResponse, err := c.carApi.GetCarInfo(regNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get car: %w", err)
	}
	return c.buildDomainCar(getCarResponse), err
}

func (c *carApiServiceImpl) buildDomainCar(getCarResponse *apicar.GetCarResponse) *car.Car {
	return &car.Car{
		RegistrationNumber: getCarResponse.RegNum,
		Mark:               getCarResponse.Mark,
		Model:              getCarResponse.Model,
		Year:               getCarResponse.Year,
		Owner: &car.Owner{
			Name:       getCarResponse.Owner.Name,
			Surname:    getCarResponse.Owner.Surname,
			Patronymic: getCarResponse.Owner.Patronymic,
		},
	}
}
