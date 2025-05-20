package usecase

import (
	"github.com/bllooop/nameservice/internal/domain"
	"github.com/bllooop/nameservice/internal/repository"
)

type PersonUsecase struct {
	repo repository.Person
}

func NewPersonUsecase(repo repository.Person) *PersonUsecase {
	return &PersonUsecase{
		repo: repo,
	}
}

func (s *PersonUsecase) CreatePerson(input domain.Person) (*domain.Person, error) {
	return s.repo.CreatePerson(input)
}
func (s *PersonUsecase) DeleteName(nameId int) error {
	return s.repo.DeleteName(nameId)

}
func (s *PersonUsecase) UpdateName(nameId int, input domain.UpdatePerson) (*domain.Person, error) {
	return s.repo.UpdateName(nameId, input)

}

func (s *PersonUsecase) GetPeople(filters domain.FilterParams) ([]domain.Person, error) {
	return s.repo.GetPeople(filters)
}
