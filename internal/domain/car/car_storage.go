package car

type CarStorage interface {
	Create(regNum, mark, model string, year *int16, owner *Owner) error
	Update(id int64, regNum, mark, model *string, year *int16, ownerName, ownerSurname, ownerPatronymic *string) (*Car, error)
	Delete(id int64) error
	ListByFilters(regNum, mark, model *string, year *int16, limit, offset int) ([]*Car, error)
}
