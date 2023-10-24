package http

import (
	"effective_mobile_test/internal/user/repository"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, rep *repository.Repo) {

	v1 := r.Group("/v1")
	{
		v1.POST("/user", rep.CreateUser())
		v1.GET("/user/all", rep.GetUsers())
		v1.PUT("/user/:id", rep.UpdateUser())
		v1.DELETE("/user/:id", rep.DeleteUser())
		v1.GET("/user/:id", rep.EnrichById())
	}
}
