package apicar

import "github.com/vaberof/effective-mobile-backend/pkg/integration/api/apicar"

type CarApi interface {
	GetCarInfo(regNum string) (*apicar.GetCarResponse, error)
}
