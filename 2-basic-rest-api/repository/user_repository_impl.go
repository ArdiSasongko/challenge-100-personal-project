package repository

import (
	"basic-rest-api/model/domain"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(DB *sql.DB) *UserRepository {
	return &UserRepository{
		DB: DB,
	}
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	rows, err := r.DB.Query("SELECT * FROM users WHERE email = $1", email) // Konsistensi nama tabel users
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Cek apakah ada hasil query
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	user, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(id int) (*domain.User, error) {
	rows, err := r.DB.Query("SELECT * FROM users WHERE id = $1", id) // Konsistensi nama tabel users
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	user, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func scanRows(rows *sql.Rows) (*domain.User, error) {
	user := new(domain.User)

	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CreateUser(user domain.User) error {
	stmt, err := r.DB.Prepare("INSERT INTO users(username, password, name, email) VALUES ($1, $2, $3, $4) RETURNING id") // Konsistensi nama tabel users
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(user.Username, user.Password, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}
