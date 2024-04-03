package pgcar

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vaberof/effective-mobile-backend/internal/domain/car"
	"github.com/vaberof/effective-mobile-backend/internal/infra/storage"
	"strconv"
)

type PgCarStorage struct {
	db *sqlx.DB
}

func NewPgCarStorage(db *sqlx.DB) *PgCarStorage {
	return &PgCarStorage{db: db}
}

func (s *PgCarStorage) Create(regNum, mark, model string, year *int16, owner *car.Owner) error {
	var car PgCar
	query := `
			INSERT INTO cars (
			                    reg_num,
			                    mark,
			                    model,
			                    year,
			                  	owner
				) VALUES ($1, $2, $3, $4, $5)
				RETURNING
				    id,
					reg_num,
					mark,
					model,
				  	year,
				  	owner
`
	pgOwner := &Owner{
		Name:       sql.NullString{String: owner.Name, Valid: true},
		Surname:    sql.NullString{String: owner.Surname, Valid: true},
		Patronymic: sql.NullString{String: owner.Patronymic, Valid: true},
	}

	pgOwnerBytes, _ := pgOwner.Value()

	row := s.db.QueryRow(query, regNum, mark, model, year, pgOwnerBytes)
	if err := row.Scan(
		&car.Id,
		&car.RegistrationNumber,
		&car.Mark,
		&car.Model,
		&car.Year,
		&car.Owner,
	); err != nil {
		return fmt.Errorf("failed to create a car: %w", err)
	}

	return nil
}

func (s *PgCarStorage) Update(id int64, regNum, mark, model *string, year *int16, ownerName, ownerSurname, ownerPatronymic *string) (*car.Car, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction while updating car: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("LOCK TABLE cars IN SHARE ROW EXCLUSIVE MODE")
	if err != nil {
		return nil, fmt.Errorf("failed to lock 'cars' table while updating car: %w", err)
	}

	var carOwner *Owner

	if ownerName != nil || ownerSurname != nil || ownerPatronymic != nil {
		carOwner, err = s.getUpdatedCarOwner(tx, id, ownerName, ownerSurname, ownerPatronymic)
		if err != nil {
			return nil, err
		}
	}

	var car PgCar

	carOwnerBytes, _ := carOwner.Value()

	query := `
			UPDATE cars
						SET 
							reg_num = COALESCE($1, reg_num),
							mark = COALESCE($2, mark),
							model = COALESCE($3, model),
							year = COALESCE($4, year),
							owner = COALESCE($5, owner)
						WHERE id=$6
			RETURNING
			    id,
				reg_num,
				mark,
				model,
				year,
				owner
	`

	row := tx.QueryRow(query, regNum, mark, model, year, carOwnerBytes, id)
	if err = row.Scan(
		&car.Id,
		&car.RegistrationNumber,
		&car.Mark,
		&car.Model,
		&car.Year,
		&car.Owner,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to update car in database: %w", storage.ErrCarNotFound)
		}
		return nil, fmt.Errorf("failed to update car in database: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction while updating car: %w", err)
	}

	return buildDomainCar(&car), nil
}

func (s *PgCarStorage) getUpdatedCarOwner(tx *sql.Tx, carId int64, ownerName, ownerSurname, ownerPatronymic *string) (*Owner, error) {
	queryGetOwner := `
					SELECT owner FROM cars WHERE id=$1
`
	var owner *Owner
	var ownerBytes []byte

	row := tx.QueryRow(queryGetOwner, carId)
	if err := row.Scan(&ownerBytes); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrCarNotFound
		}

		return nil, err
	}

	if len(ownerBytes) > 0 {
		if err := json.Unmarshal(ownerBytes, &owner); err != nil {
			return nil, err
		}
	}

	if ownerName != nil {
		owner.Name = sql.NullString{String: *ownerName, Valid: true}
	}
	if ownerSurname != nil {
		owner.Surname = sql.NullString{String: *ownerSurname, Valid: true}
	}
	if ownerPatronymic != nil {
		owner.Patronymic = sql.NullString{String: *ownerPatronymic, Valid: true}
	}

	return owner, nil
}

func (s *PgCarStorage) Delete(id int64) error {
	query := `DELETE FROM cars WHERE id=$1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("failed to delete car: %w", storage.ErrCarNotFound)
	}
	return nil
}

func (s *PgCarStorage) ListByFilters(regNum, mark, model *string, year *int16, limit, offset int) ([]*car.Car, error) {
	limitOffsetParams := fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, offset)

	query := `
			SELECT * FROM cars WHERE 1=1
`
	args := []interface{}{}

	if regNum != nil {
		query += "AND reg_num LIKE $" + strconv.Itoa(len(args)+1) + " "
		args = append(args, "%"+*regNum+"%")
	}
	if mark != nil {
		query += "AND mark LIKE $" + strconv.Itoa(len(args)+1) + " "
		args = append(args, "%"+*mark+"%")
	}
	if model != nil {
		query += "AND model LIKE $" + strconv.Itoa(len(args)+1) + " "
		args = append(args, "%"+*model+"%")
	}
	if year != nil {
		query += "AND year = $" + strconv.Itoa(len(args)+1) + " "
		args = append(args, *year)
	}

	query += limitOffsetParams

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []*PgCar

	for rows.Next() {
		var car PgCar
		err = rows.Scan(&car.Id, &car.RegistrationNumber, &car.Mark, &car.Model, &car.Year, &car.Owner)
		if err != nil {
			return nil, err
		}
		cars = append(cars, &car)
	}

	return buildDomainCars(cars), nil
}

func buildDomainCars(postgresCars []*PgCar) []*car.Car {
	domainCars := make([]*car.Car, len(postgresCars))
	for i := range postgresCars {
		domainCars[i] = buildDomainCar(postgresCars[i])
	}
	return domainCars
}

func buildDomainCar(postgresCar *PgCar) *car.Car {
	c := &car.Car{
		Id:                 postgresCar.Id,
		RegistrationNumber: postgresCar.RegistrationNumber,
		Mark:               postgresCar.Mark.String,
		Model:              postgresCar.Model.String,
		Year:               postgresCar.Year.Int16,
	}
	if postgresCar.Owner != nil {
		c.Owner = &car.Owner{
			Name:       postgresCar.Owner.Name.String,
			Surname:    postgresCar.Owner.Surname.String,
			Patronymic: postgresCar.Owner.Patronymic.String,
		}
	}
	return c
}
