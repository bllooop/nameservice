package domain

type Person struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Age         *int    `json:"age,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	Nationality *string `json:"nationality,omitempty"`
}

type FilterParams struct {
	Name        *string
	Surname     *string
	Gender      *string
	Nationality *string
	Patronymic  *string
	AgeMin      *int
	AgeMax      *int
	Limit       int
	Page        int
}

type UpdatePerson struct {
	ID          *int64  `json:"id"`
	Name        *string `json:"name"`
	Surname     *string `json:"surname"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Age         *int    `json:"age,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	Nationality *string `json:"nationality,omitempty"`
}
