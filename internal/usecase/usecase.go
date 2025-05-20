package usecase

import (
	"github.com/bllooop/nameservice/internal/domain"
	"github.com/bllooop/nameservice/internal/repository"
)

type Person interface {
	CreatePerson(input domain.Person) (*domain.Person, error)
	DeleteName(nameId int) error
	UpdateName(nameId int, input domain.UpdatePerson) (*domain.Person, error)
	GetPeople(filters domain.FilterParams) ([]domain.Person, error)
}
type Usecase struct {
	Person
}

func NewUsecase(repo *repository.Repository) *Usecase {
	return &Usecase{
		Person: NewPersonUsecase(repo),
	}
}
