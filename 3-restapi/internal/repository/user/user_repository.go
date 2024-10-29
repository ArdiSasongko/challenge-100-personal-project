package userrepository

import (
	"context"
	"database/sql"
	"fmt"
	"restapi/internal/model/domain"
	"time"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(ctx context.Context, tx *sql.Tx, user domain.User) (*domain.User, error) {
	// set timeout after 5 second
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := tx.QueryRowContext(ctx,
		`INSERT INTO users(name, age, username, email, password) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Name, user.Age, user.Username, user.Email, user.Password,
	).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed create new user")
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, tx *sql.Tx, user domain.User, userID int) (*domain.User, error) {
	// set timeout after 5 second
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := tx.ExecContext(ctx,
		`UPDATE users SET name = $1, age = $2, username = $3, email = $4, password = $5 WHERE id = $6`,
		user.Name, user.Age, user.Username, user.Email, user.Password, userID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed update user")
	}

	err = CheckRows(result, userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, tx *sql.Tx, userID int) error {
	// set timeout after 5 second
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := tx.ExecContext(ctx,
		`DELETE users WHERE id = $1`,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed delete user")
	}

	err = CheckRows(result, userID)
	if err != nil {
		return err
	}

	return nil
}

func CheckRows(result sql.Result, userID int) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user with id %d not found", userID)
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, tx *sql.Tx, userID int) (*domain.User, error) {
	// set timeout after 5 second
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := tx.QueryContext(ctx,
		`SELECT id, name, age, username, email, password, created_at FROM users WHERE id = $1`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	user, err := ScanRows(rows)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*domain.User, error) {
	// set timeout after 5 second
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := tx.QueryContext(ctx,
		`SELECT id, name, age, username, email, password, created_at FROM users WHERE email = $1`,
		email,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	user, err := ScanRows(rows)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := tx.QueryContext(ctx,
		`SELECT id, name, age, username, email, password FROM users`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user, err := ScanRows(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}
func ScanRows(rows *sql.Rows) (*domain.User, error) {
	user := new(domain.User)

	err := rows.Scan(
		&user.ID,
		&user.Name,
		&user.Age,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
