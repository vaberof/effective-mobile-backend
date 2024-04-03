package car

type Car struct {
	Id                 int64
	RegistrationNumber string
	Mark               string
	Model              string
	Year               int16
	Owner              *Owner
}

type Owner struct {
	Name       string
	Surname    string
	Patronymic string
}
