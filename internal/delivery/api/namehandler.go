package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bllooop/nameservice/internal/domain"
	applog "github.com/bllooop/nameservice/internal/log"

	"github.com/gin-gonic/gin"
)

var agify struct {
	Age int `json:"age"`
}
var genderize struct {
	Gender string `json:"gender"`
}

var nationalizeData struct {
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

// @Summary Get list of people
// @Tags people
// @Description Получение списка людей из бд с применением фильтров и пагинаций. Параметры, передающиеся в адресе запроса, представлены ниже, можно заполнять любое желаемое количество. Для возраста представлено два параметра, ограничиващих выборку.
// @ID get-people
// @Param name query string false "Имя"
// @Param surname query string false "Фамилия"
// @Param patronymic query string false "Отчество"
// @Param gender query string false "Пол"
// @Param nationality query string false "Национальность"
// @Param age_min query int false "Минимальный возраст"
// @Param age_max query int false "Максимальный возраст"
// @Param page query int false "Страница"
// @Param limit query int false "Лимит"
// @Produce  json
// @Success 200 {object} domain.PersonResponse
// @Failure 400,404,500 {object} api.ErrorResponse
// @Router /get_people [get]
func (h *Handler) GetPeople(c *gin.Context) {
	applog.Logger.Info().Msg("Получен запрос на получение данны о ПВЗ")
	if c.Request.Method != http.MethodGet {
		applog.Logger.Error().Msg("Требуется запрос GET")
		newErrorResponse(c, http.StatusBadRequest, "Требуется запрос GET")
		return
	}
	name := c.Query("name")
	surname := c.Query("surname")
	gender := c.Query("gender")
	nationality := c.Query("nationality")
	age_min := c.Query("age_min")
	age_max := c.Query("age_max")
	patronymic := c.Query("patronymic")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 10
	}
	ageInt, err := strconv.Atoi(age_min)
	if err != nil || ageInt < 0 {
		ageInt = -1
	}
	ageMaxInt, err := strconv.Atoi(age_max)
	if err != nil || ageMaxInt < 0 {
		ageMaxInt = -1
	}
	filters := domain.FilterParams{
		Name:        &name,
		Surname:     &surname,
		Gender:      &gender,
		Nationality: &nationality,
		Patronymic:  &patronymic,
		AgeMin:      &ageInt,
		AgeMax:      &ageMaxInt,
		Limit:       limitInt,
		Page:        pageInt,
	}
	applog.Logger.Debug().Msgf(
		"Успешно прочитаны параметры из запроса %s, %s, %s, %s, %s, %d, %d, %d, %d",
		strVal(filters.Name), strVal(filters.Surname), strVal(filters.Gender),
		strVal(filters.Nationality), strVal(filters.Patronymic),
		intVal(filters.AgeMin), intVal(filters.AgeMax), filters.Limit, filters.Page,
	)
	result, err := h.Usecases.GetPeople(filters)
	if err != nil {
		applog.Logger.Error().Err(err).Msg("")
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка выполнения запроса "+err.Error())
		return
	}

	applog.Logger.Info().Msg("Получен ответ на запрос сущностей")
	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
		"data":   result,
	})
}

// @Summary Create person
// @Tags people
// @Description Создание сущности человека в БД при использовании переданных в запросе данных. Структура тела запроса представлена ниже, поле "Отчество" необязательно.
// @ID create-person
// @Accept json
// @Produce  json
// @Param input body domain.InputPerson true "list info"
// @Success 200 {object} domain.PersonResponse
// @Failure 400,404,500 {object} api.ErrorResponse
// @Router /create_person [post]
func (h *Handler) CreatePerson(c *gin.Context) {
	applog.Logger.Info().Msg("Получен запрос на запись сущности в БД")
	if c.Request.Method != http.MethodPost {
		applog.Logger.Error().Msg("Требуется запрос POST")
		newErrorResponse(c, http.StatusBadRequest, "Неверный запрос")
		return
	}
	var input domain.InputPerson
	if err := c.ShouldBindJSON(&input); err != nil {
		applog.Logger.Error().Err(err).Msg(err.Error())
		newErrorResponse(c, http.StatusBadRequest, "Неверный запрос")
		return
	}
	applog.Logger.Debug().Msgf("Успешно прочитаны данные из запроса  %s, %s, %v", input.Name, input.Surname, input.Patronymic)
	var addedInput domain.Person
	addedInput.Name = input.Name
	addedInput.Surname = input.Surname
	addedInput.Patronymic = input.Patronymic
	err := getApiData(&addedInput)
	if err != nil {
		applog.Logger.Error().Err(err).Msg(err.Error())
		newErrorResponse(c, http.StatusBadRequest, "Ошибка запроса в API")
		return
	}
	result, err := h.Usecases.Person.CreatePerson(addedInput)
	if err != nil {
		applog.Logger.Error().Err(err).Msg("")
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка выполнения запроса "+err.Error())
		return
	}
	applog.Logger.Info().Msg("Получен ответ о добавлении данных с фио в БД")
	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
		"data":   result,
	})
}

// @Summary Delete person by ID
// @Tags people
// @Description Удаление сущности человека из БД по указнному ID. ID передается в адресе запроса.
// @ID delete-person
// @Param nameId query int true "ID"
// @Produce  json
// @Success 200 {string} message
// @Failure 400,404,500 {object} api.ErrorResponse
// @Router /delete_person [delete]
func (h *Handler) DeletePerson(c *gin.Context) {
	applog.Logger.Info().Msg("Получен запрос на удаление фио по идентификатору")
	if c.Request.Method != http.MethodDelete {
		applog.Logger.Error().Msg("Требуется запрос DELETE")
		newErrorResponse(c, http.StatusBadRequest, "Требуется запрос DELETE")
		return
	}
	nameId, err := strconv.Atoi(c.Query("nameId"))
	if err != nil {
		applog.Logger.Error().Err(err).Msg("")
		newErrorResponse(c, http.StatusInternalServerError, "Некорректный ввод ID")
		return
	}
	applog.Logger.Debug().Msgf("Успешно прочитан параметр из запроса %v", nameId)

	err = h.Usecases.Person.DeleteName(nameId)
	if err != nil {
		applog.Logger.Error().Err(err).Msg("")
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка выполнения запроса "+err.Error())
		return
	}
	applog.Logger.Info().Msg("Получен ответ на удаление фио")
	c.JSON(http.StatusOK, map[string]any{
		"message": "Фио удалено",
	})
}

// @Summary Update person by ID
// @Tags people
// @Description Обновление сущности человека из БД по указнному ID, применяя данные из запроса. Структура тела запроса представлена ниже, можно заполнить желаемое количество полей.
// @ID update-person
// @Param nameId query int false "ID"
// @Accept json
// @Produce  json
// @Param input body domain.UpdatePerson true "list info"
// @Success 200 {object} domain.PersonResponse
// @Failure 400,404,500 {object} api.ErrorResponse
// @Router /update_person [patch]
func (h *Handler) UpdateName(c *gin.Context) {
	applog.Logger.Info().Msg("Получен запрос на измнение сущности фио")
	if c.Request.Method != http.MethodPatch {
		newErrorResponse(c, http.StatusBadRequest, "Требуется запрос PATCH")
		return
	}
	nameId, err := strconv.Atoi(c.Query("nameId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "недопустимое значение идентификатора")
		return
	}
	applog.Logger.Debug().Int("id parameter", nameId).Msg("Успешно прочитан идентификатор сущности ФИО")

	var input domain.UpdatePerson
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	applog.Logger.Debug().
		Interface("binded_person", input).
		Msg("Успешная обработка JSON в структуру")

	res, err := h.Usecases.Person.UpdateName(nameId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	applog.Logger.Info().Msg("Получен ответ на обновление сущности")

	c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
		"data":   res,
	})
}

func getApiData(input *domain.Person) error {
	ageUrl := fmt.Sprintf("https://api.agify.io/?name=%s", input.Name)
	resp, err := http.Get(ageUrl)
	if err != nil {
		applog.Logger.Error().Err(err).Msg(err.Error())
		return fmt.Errorf("Неверный запрос в agefiy")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		applog.Logger.Error().Err(err).Msg(err.Error())
		return fmt.Errorf("Agefiy вернул ошибку")
	}
	if err := json.NewDecoder(resp.Body).Decode(&agify); err != nil {
		applog.Logger.Error().Err(err).Msg(err.Error())
		return fmt.Errorf("Ошибка обработки ответа Agify")
	}
	applog.Logger.Debug().Msgf("Успешно получен возраст %v", &agify.Age)

	input.Age = &agify.Age

	genderURL := fmt.Sprintf("https://api.genderize.io/?name=%s", input.Name)
	genderResp, err := http.Get(genderURL)
	if err != nil {
		applog.Logger.Error().Err(err).Msg(err.Error())
		return fmt.Errorf("Неверный запрос в genderize")
	}
	defer genderResp.Body.Close()

	if genderResp.StatusCode != http.StatusOK {
		applog.Logger.Error().Err(err).Msg(err.Error())
		return fmt.Errorf("genderize вернул ошибку")
	}
	if err := json.NewDecoder(genderResp.Body).Decode(&genderize); err != nil {
		applog.Logger.Error().Err(err).Msg(err.Error())
		return fmt.Errorf("Ошибка обработки ответа Genderize")
	}
	applog.Logger.Debug().Msgf("Успешно получен пол %v", &genderize.Gender)

	input.Gender = &genderize.Gender
	nationalUrl := fmt.Sprintf("https://api.nationalize.io/?name=%s", input.Name)
	nationalResp, err := http.Get(nationalUrl)
	if err != nil {
		applog.Logger.Error().Err(err).Msg(err.Error())
		return fmt.Errorf("Неверный запрос в nationalize")
	}
	defer nationalResp.Body.Close()

	if nationalResp.StatusCode != http.StatusOK {
		applog.Logger.Error().Err(err).Msg(err.Error())
		return fmt.Errorf("Nationalize вернул ошибку")
	}
	if err := json.NewDecoder(nationalResp.Body).Decode(&nationalizeData); err != nil {
		applog.Logger.Error().Err(err).Msg(err.Error())
		return fmt.Errorf("Ошибка обработки ответа nationalize")
	}

	if len(nationalizeData.Country) > 0 {
		max := nationalizeData.Country[0]
		for _, c := range nationalizeData.Country {
			if c.Probability > max.Probability {
				max = c
			}
		}
		applog.Logger.Debug().Msgf("Успешно получена национальность %v", &max.CountryID)

		input.Nationality = &max.CountryID
	}
	return nil
}

func strVal(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}

func intVal(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
