package api

import (
	"time"

	_ "github.com/bllooop/nameservice/docs"
	"github.com/bllooop/nameservice/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Usecases *usecase.Usecase
	Now      func() time.Time
}

func NewHandler(usecases *usecase.Usecase) *Handler {
	return &Handler{Usecases: usecases, Now: func() time.Time { return time.Now() }}
}
func NewHandlerWithFixedTime(usecases *usecase.Usecase, fixedTime time.Time) *Handler {
	return &Handler{
		Usecases: usecases,
		Now:      func() time.Time { return fixedTime },
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/create_person", h.CreatePerson)
	router.DELETE("/delete_person", h.DeletePerson)
	router.PATCH("/update_person", h.UpdateName)
	router.GET("/get_people", h.GetPeople)
	return router
}
