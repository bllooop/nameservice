package usecase

type Name interface {
}
type Usecase struct {
	Name
}

func NewUsecase(repo *repository.Repository) *Usecase {
	return &Usecase{
		Name: NewNameUsecase(repo),
	}
}
