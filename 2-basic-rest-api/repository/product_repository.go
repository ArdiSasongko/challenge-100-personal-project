package repository

import "database/sql"

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(DB *sql.DB) *ProductRepository {
	return &ProductRepository{
		DB: DB,
	}
}
