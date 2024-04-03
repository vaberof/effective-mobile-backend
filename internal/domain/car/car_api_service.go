package car

type CarApiService interface {
	GetCar(regNum string) (*Car, error)
}
