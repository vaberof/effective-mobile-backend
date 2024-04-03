package http

import "github.com/vaberof/effective-mobile-backend/internal/domain/car"

type createCarsRequestBody struct {
	RegNums []string `json:"regNums"`
}

type createCarsResponseBody struct {
	Message string `json:"message"`
}

type carOwnerPayload struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type updateCarRequestBody struct {
	RegNum *string                `json:"regNum,omitempty"`
	Mark   *string                `json:"mark,omitempty"`
	Model  *string                `json:"model,omitempty"`
	Year   *int16                 `json:"year,omitempty"`
	Owner  *updateCarOwnerPayload `json:"owner,omitempty"`
}

type updateCarOwnerPayload struct {
	Name       *string `json:"name,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	Patronymic *string `json:"patronymic,omitempty"`
}

type updateCarResponseBody struct {
	Id     int64            `json:"id"`
	RegNum string           `json:"regNum"`
	Mark   string           `json:"mark"`
	Model  string           `json:"model"`
	Year   int16            `json:"year,omitempty"`
	Owner  *carOwnerPayload `json:"owner"`
}

type deleteCarResponseBody struct {
	Message string `json:"message"`
}

type listCarsResponseBody struct {
	Cars []*carPayload `json:"cars"`
}

type carPayload struct {
	Id     int64            `json:"id"`
	RegNum string           `json:"regNum"`
	Mark   string           `json:"mark"`
	Model  string           `json:"model"`
	Year   int16            `json:"year,omitempty"`
	Owner  *carOwnerPayload `json:"owner"`
}

func buildCarPayloads(domainCars []*car.Car) []*carPayload {
	carPayloads := make([]*carPayload, len(domainCars))
	for i := range domainCars {
		carPayloads[i] = buildCarPayload(domainCars[i])
	}
	return carPayloads
}

func buildCarPayload(domainCar *car.Car) *carPayload {
	c := &carPayload{
		Id:     domainCar.Id,
		RegNum: domainCar.RegistrationNumber,
		Mark:   domainCar.Mark,
		Model:  domainCar.Model,
		Year:   domainCar.Year,
	}
	if domainCar.Owner != nil {
		c.Owner = &carOwnerPayload{
			Name:       domainCar.Owner.Name,
			Surname:    domainCar.Owner.Surname,
			Patronymic: domainCar.Owner.Patronymic,
		}
	}
	return c
}
