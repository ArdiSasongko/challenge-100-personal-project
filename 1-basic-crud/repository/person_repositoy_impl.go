package repository

import (
	"basic-crud/model/domain"
	"context"
	"database/sql"
	"errors"
	"strconv"
)

type PersonRepository struct {
}

func NewPersonRepository() *PersonRepository {
	return &PersonRepository{}
}

func (r *PersonRepository) Create(ctx context.Context, tx *sql.Tx, person domain.Person) (*domain.Person, error) {
	sql := `INSERT INTO person (name, age) VALUES ($1, $2) RETURNING id`
	err := tx.QueryRowContext(ctx, sql, person.Name, person.Age).Scan(&person.ID)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *PersonRepository) Update(ctx context.Context, tx *sql.Tx, person domain.Person) (*domain.Person, error) {
	sql := `UPDATE person set name = $1, age = $2 WHERE id = $3`
	_, err := tx.ExecContext(ctx, sql, person.Name, person.Age, person.ID)
	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (r *PersonRepository) Delete(ctx context.Context, tx *sql.Tx, personID int) error {
	sql := `DELETE FROM person WHERE id = $1`
	_, err := tx.ExecContext(ctx, sql, personID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PersonRepository) FindByID(ctx context.Context, tx *sql.Tx, personID int) (*domain.Person, error) {
	sql := `SELECT id, name, age FROM person WHERE id = $1`
	rows, err := tx.QueryContext(ctx, sql, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	person := domain.Person{}
	if rows.Next() {
		err = rows.Scan(&person.ID, &person.Name, &person.Age)
		if err != nil {
			return nil, err
		}
		return &person, nil
	} else {
		return nil, errors.New("person with id " + strconv.Itoa(personID) + " not found")
	}
}

func (r *PersonRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Person, error) {
	sql := `SELECT id, name, age FROM person`
	rows, err := tx.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []domain.Person
	for rows.Next() {
		person := domain.Person{}
		err := rows.Scan(&person.ID, &person.Name, &person.Age)
		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return persons, nil
}
