package repository

import (
	"fmt"
	"strings"

	"github.com/bllooop/nameservice/internal/domain"
	applog "github.com/bllooop/nameservice/internal/log"
	"github.com/jmoiron/sqlx"
)

type PersonPostgres struct {
	db *sqlx.DB
}

func NewPersonPostgres(db *sqlx.DB) *PersonPostgres {
	return &PersonPostgres{
		db: db,
	}
}

func (r *PersonPostgres) GetPeople(filters domain.FilterParams) ([]domain.Person, error) {
	conditions, args := buildConditions(filters)
	offset := (filters.Page - 1) * filters.Limit
	args = append(args, filters.Limit, offset)
	query := fmt.Sprintf("SELECT * FROM %s", peopleListTable)
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += fmt.Sprintf(" ORDER BY  %s  %s LIMIT $%d OFFSET $%d", *filters.SortBy, *filters.OrderBy, len(args)-1, len(args))
	var people []domain.Person
	err := r.db.Select(&people, query, args...)
	if err != nil {
		return nil, err
	}
	applog.Logger.Debug().Any("query", query).Msg("Запрос данных о людях")
	return people, nil
}

func buildConditions(input domain.FilterParams) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}
	argId := 1
	if input.Name != nil && *input.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name=$%d", argId))
		args = append(args, input.Name)
		argId++
	}
	if input.Surname != nil && *input.Surname != "" {
		conditions = append(conditions, fmt.Sprintf("surname=$%d", argId))
		args = append(args, input.Surname)
		argId++
	}
	if input.Patronymic != nil && *input.Patronymic != "" {
		conditions = append(conditions, fmt.Sprintf("patronymic=$%d", argId))
		args = append(args, input.Patronymic)
		argId++
	}
	if input.Gender != nil && *input.Gender != "" {
		conditions = append(conditions, fmt.Sprintf("gender=$%d", argId))
		args = append(args, input.Gender)
		argId++
	}
	if input.Nationality != nil && *input.Nationality != "" {
		conditions = append(conditions, fmt.Sprintf("nationality=$%d", argId))
		args = append(args, input.Nationality)
		argId++
	}
	if input.AgeMin != nil && *input.AgeMin >= 0 {
		conditions = append(conditions, fmt.Sprintf("age >= $%d", argId))
		args = append(args, input.AgeMin)
		argId++
	}
	if input.AgeMax != nil && *input.AgeMax >= 0 {
		conditions = append(conditions, fmt.Sprintf("age <= $%d", argId))
		args = append(args, input.AgeMax)
		argId++
	}
	return conditions, args
}
func (r *PersonPostgres) CreatePerson(input domain.Person) (*domain.Person, error) {
	var result domain.Person
	query := fmt.Sprintf(`INSERT INTO %s (name,surname, patronymic, age, gender, nationality) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id,name,surname, patronymic, age, gender, nationality`, peopleListTable)
	row := r.db.QueryRowx(query, input.Name, input.Surname, input.Patronymic, input.Age, input.Gender, input.Nationality)
	applog.Logger.Debug().Str("query", query).Msg("Выполнение запроса сохранения сущности человека в БД")
	if err := row.Scan(&result.ID, &result.Name, &result.Surname, &result.Patronymic, &result.Age, &result.Gender, &result.Nationality); err != nil {
		return &domain.Person{}, err
	}
	applog.Logger.Debug().Any("person response", result).Msg("Успешно добавлена обогащенная сущность человека в БД")
	return &result, nil
}
func (r *PersonPostgres) DeleteName(nameId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", peopleListTable)
	applog.Logger.Debug().
		Str("query", query).
		Int("id", nameId).
		Msg("Удаление сущности фио")
	_, err := r.db.Exec(query, nameId)
	return err
}
func (r *PersonPostgres) UpdateName(nameId int, input domain.UpdatePerson) (*domain.Person, error) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}
	if input.Surname != nil {
		setValues = append(setValues, fmt.Sprintf("surname=$%d", argId))
		args = append(args, *input.Surname)
		argId++
	}
	if input.Patronymic != nil {
		setValues = append(setValues, fmt.Sprintf("patronymic=$%d", argId))
		args = append(args, *input.Patronymic)
		argId++
	}
	if input.Age != nil {
		setValues = append(setValues, fmt.Sprintf("age=$%d", argId))
		args = append(args, *input.Age)
		argId++
	}
	if input.Gender != nil {
		setValues = append(setValues, fmt.Sprintf("gender=$%d", argId))
		args = append(args, *input.Gender)
		argId++
	}
	if input.Nationality != nil {
		setValues = append(setValues, fmt.Sprintf("nationality=$%d", argId))
		args = append(args, *input.Nationality)
		argId++
	}
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d RETURNING id, name,surname, patronymic, age, gender, nationality", peopleListTable, setQuery, argId)
	applog.Logger.Debug().
		Str("query", query).
		Msgf("Выполняем запрос обновления следующих параметров: %v", args)
	args = append(args, nameId)
	result := domain.Person{}
	err := r.db.Get(&result, query, args...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func IntPointer(s int) *int {
	return &s
}
