package repository

import (
	"github.com/bllooop/nameservice/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Person interface {
	CreatePerson(input domain.Person) (*domain.Person, error)
	DeleteName(nameId int) error
	UpdateName(nameId int, input domain.UpdatePerson) (*domain.Person, error)
	GetPeople(filters domain.FilterParams) ([]domain.Person, error)
}

type Repository struct {
	Person
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Person: NewPersonPostgres(db),
	}
}
