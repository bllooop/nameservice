package api

import (
	applog "github.com/bllooop/nameservice/internal/log"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	applog.Logger.Error().Msg(message)
	c.AbortWithStatusJSON(statusCode, ErrorResponse{message})
}
