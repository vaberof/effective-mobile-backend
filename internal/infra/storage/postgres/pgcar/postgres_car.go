package pgcar

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type PgCar struct {
	Id                 int64
	RegistrationNumber string
	Mark               sql.NullString
	Model              sql.NullString
	Year               sql.NullInt16
	Owner              *Owner `json:"owner,omitempty"`
}

type Owner struct {
	Name       sql.NullString `json:"name"`
	Surname    sql.NullString `json:"surname"`
	Patronymic sql.NullString `json:"patronymic,omitempty"`
}

func (o *Owner) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func (o *Owner) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &o)
}
