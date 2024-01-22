package http

import (
	"github.com/gin-gonic/gin"
	"go-jun/internal/person/config"
	"go-jun/internal/usecase"
	"go.uber.org/zap"
)

type userRouter struct {
	p   usecase.PersonUC
	l   *zap.SugaredLogger
	cfg config.Config
}

func UserRouter(r *gin.Engine, p usecase.PersonUC, l *zap.SugaredLogger, cfg config.Config) {
	ur := userRouter{
		p:   p,
		l:   l,
		cfg: cfg,
	}
	v1 := r.Group("api/people/v1")
	{
		v1.GET("/people", ur.p.GetPersons())
		v1.POST("/people", ur.p.CreatePerson())
		v1.PUT("/people/:id", ur.p.UpdatePerson())
		v1.DELETE("/people/:id", ur.p.DeletePerson())
	}

}
