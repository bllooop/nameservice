package domain

type Person struct {
	ID          int64   `json:"-"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Age         *int    `json:"age,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	Nationality *string `json:"nationality,omitempty"`
}

type InputPerson struct {
	ID         int64   `json:"-"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic,omitempty"`
}

type PersonResponse struct {
	Status string `json:"status" example:"ok"`
	Data   Person `json:"data"`
}
type FilterParams struct {
	Name        *string
	Surname     *string
	Gender      *string
	Nationality *string
	Patronymic  *string
	AgeMin      *int
	AgeMax      *int
	SortBy      *string
	OrderBy     *string
	Limit       int
	Page        int
}

type UpdatePerson struct {
	ID          *int64  `json:"-"`
	Name        *string `json:"name"`
	Surname     *string `json:"surname"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Age         *int    `json:"age,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	Nationality *string `json:"nationality,omitempty"`
}
