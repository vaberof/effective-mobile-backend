package car_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/vaberof/effective-mobile-backend/internal/domain/car"
	mocks "github.com/vaberof/effective-mobile-backend/internal/domain/car/mocks"
	"github.com/vaberof/effective-mobile-backend/pkg/logging/logs"
	"go.uber.org/mock/gomock"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	carApiService := mocks.NewMockCarApiService(ctrl)
	carStorage := mocks.NewMockCarStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	carService := car.NewCarService(carStorage, carApiService, logsBuilder)

	regNums := []string{"X123XX123"}

	id := int64(0)
	regNum := regNums[0]
	mark := "Mark1"
	model := "Model1"
	year := int16(2000)
	ownerName := "Name1"
	ownerSurname := "Surname1"
	ownerPatronymic := "Patronymic"

	owner := car.Owner{
		Name:       ownerName,
		Surname:    ownerSurname,
		Patronymic: ownerPatronymic,
	}

	apiServiceExpected := &car.Car{
		Id:                 id,
		RegistrationNumber: regNum,
		Mark:               mark,
		Model:              model,
		Year:               year,
		Owner:              &owner,
	}

	carApiService.EXPECT().GetCar(regNum).Return(apiServiceExpected, nil).AnyTimes()
	carStorage.EXPECT().Create(regNum, mark, model, &year, apiServiceExpected.Owner).Return(nil).AnyTimes()
	err := carService.Create(regNums)
	require.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	carApiService := mocks.NewMockCarApiService(ctrl)
	carStorage := mocks.NewMockCarStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	carService := car.NewCarService(carStorage, carApiService, logsBuilder)

	id := int64(1)
	regNum := "X123XX123"
	mark := "Mark1"
	model := "Model1"
	year := int16(2000)
	ownerName := "Name1"
	ownerSurname := "Surname1"
	ownerPatronymic := "Patronymic1"

	expected := &car.Car{
		Id:                 id,
		RegistrationNumber: regNum,
		Mark:               mark,
		Model:              model,
		Year:               year,
		Owner: &car.Owner{
			Name:       ownerName,
			Surname:    ownerSurname,
			Patronymic: ownerPatronymic,
		},
	}

	carStorage.EXPECT().Update(id, &regNum, &mark, &model, &year, &ownerName, &ownerSurname, &ownerPatronymic).Return(expected, nil).Times(1)
	domainCar, err := carService.Update(id, &regNum, &mark, &model, &year, &ownerName, &ownerSurname, &ownerPatronymic)
	require.NoError(t, err)
	require.Equal(t, expected, domainCar)
}

func TestUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	carApiService := mocks.NewMockCarApiService(ctrl)
	carStorage := mocks.NewMockCarStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	carService := car.NewCarService(carStorage, carApiService, logsBuilder)

	type in struct {
		Id                 int64
		RegistrationNumber string
		Mark               string
		Model              string
		Year               int16
		Owner              struct {
			Name       string
			Surname    string
			Patronymic string
		}
	}

	testCases := []struct {
		name   string
		in     in
		out    *car.Car
		expErr error
	}{
		{
			name: "err_car_not_found",
			in: in{
				Id:                 1,
				RegistrationNumber: "X123XX123",
				Mark:               "Mark1",
				Model:              "Model1",
				Year:               2000,
				Owner: struct {
					Name       string
					Surname    string
					Patronymic string
				}{Name: "Name1", Surname: "Surname1", Patronymic: "Patronymic1"},
			},
			out:    nil,
			expErr: car.ErrCarNotFound,
		},
		{
			name: "err_other",
			in: in{
				Id:                 2,
				RegistrationNumber: "X124XX124",
				Mark:               "Mark2",
				Model:              "Model2",
				Year:               2001,
				Owner: struct {
					Name       string
					Surname    string
					Patronymic string
				}{Name: "Name2", Surname: "Surname2", Patronymic: "Patronymic2"},
			},
			out:    nil,
			expErr: fmt.Errorf("failed to update a car: %w", errors.New("database is down")),
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			carStorage.EXPECT().Update(tCase.in.Id, &tCase.in.RegistrationNumber, &tCase.in.Mark, &tCase.in.Model, &tCase.in.Year, &tCase.in.Owner.Name, &tCase.in.Owner.Surname, &tCase.in.Owner.Patronymic).Return(tCase.out, tCase.expErr).AnyTimes()
			domainCar, err := carService.Update(tCase.in.Id, &tCase.in.RegistrationNumber, &tCase.in.Mark, &tCase.in.Model, &tCase.in.Year, &tCase.in.Owner.Name, &tCase.in.Owner.Surname, &tCase.in.Owner.Patronymic)
			require.Error(t, err)
			require.EqualError(t, tCase.expErr, err.Error())
			require.Nil(t, domainCar)
		})
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	carApiService := mocks.NewMockCarApiService(ctrl)
	carStorage := mocks.NewMockCarStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	carService := car.NewCarService(carStorage, carApiService, logsBuilder)

	carId := int64(1)

	carStorage.EXPECT().Delete(carId).Return(nil).Times(1)

	err := carService.Delete(carId)
	require.NoError(t, err)
}

func TestDeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	carApiService := mocks.NewMockCarApiService(ctrl)
	carStorage := mocks.NewMockCarStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	carService := car.NewCarService(carStorage, carApiService, logsBuilder)

	testCases := []struct {
		name   string
		in     int64
		expErr error
	}{
		{
			name:   "err_car_not_found",
			in:     1,
			expErr: car.ErrCarNotFound,
		},
		{
			name:   "err_other",
			in:     2,
			expErr: fmt.Errorf("failed to delete a car: %w", errors.New("database is down")),
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			carStorage.EXPECT().Delete(tCase.in).Return(tCase.expErr).AnyTimes()
			err := carService.Delete(tCase.in)
			require.EqualError(t, err, tCase.expErr.Error())
		})
	}
}

func TestListByFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	carApiService := mocks.NewMockCarApiService(ctrl)
	carStorage := mocks.NewMockCarStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	carService := car.NewCarService(carStorage, carApiService, logsBuilder)

	expected := []*car.Car{
		{
			Id:                 1,
			RegistrationNumber: "X123XX123",
			Mark:               "Mark1",
			Model:              "Model1",
			Year:               2000,
			Owner: &car.Owner{
				Name:       "Name1",
				Surname:    "Surname1",
				Patronymic: "Patronymic1",
			},
		},
		{
			Id:                 2,
			RegistrationNumber: "X124XX124",
			Mark:               "Mark2",
			Model:              "Model2",
			Year:               2001,
			Owner: &car.Owner{
				Name:    "Name2",
				Surname: "Surname2",
			},
		},
	}

	limit := 100
	offset := 0

	mark := "Mark"
	model := "Model"

	carStorage.EXPECT().ListByFilters(nil, &mark, &model, nil, limit, offset).Return(expected, nil).Times(1)

	cars, err := carService.ListByFilters(nil, &mark, &model, nil, limit, offset)
	require.NoError(t, err)
	require.Equal(t, expected, cars)
}

func TestListByFiltersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	carApiService := mocks.NewMockCarApiService(ctrl)
	carStorage := mocks.NewMockCarStorage(ctrl)
	logsBuilder := logs.New(os.Stdout, nil)

	carService := car.NewCarService(carStorage, carApiService, logsBuilder)

	fakeError := errors.New("database is down")
	expectedErr := fmt.Errorf("failed to list cars: %w", fakeError)

	regNum := "X123XX123"
	mark := "Mark1"
	model := "Model1"
	year := int16(2000)
	limit := 100
	offset := 0

	carStorage.EXPECT().ListByFilters(&regNum, &mark, &model, &year, limit, offset).Return(nil, expectedErr).Times(1)

	cars, err := carService.ListByFilters(&regNum, &mark, &model, &year, limit, offset)
	require.Error(t, err)
	require.EqualError(t, expectedErr, err.Error())
	require.Nil(t, cars)
}
