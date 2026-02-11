package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/handler"
)

func Register(r *gin.Engine, userHandler *handler.UserHandler) {
	api := r.Group("/api/users")
	{
		api.POST("/", userHandler.Create)
		api.GET("/", userHandler.GetAll)
	}
}
