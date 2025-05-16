package repository

import "github.com/jmoiron/sqlx"

type Name interface {
}

type Repository struct {
	Name
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Name: NewNamePostgres(db),
	}
}
