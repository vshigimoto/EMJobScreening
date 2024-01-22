package usecase

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-jun/internal/person/entity"
	"go-jun/internal/person/repository"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type PersonUC struct {
	l *zap.SugaredLogger
	r *repository.Repo
}

func New(l *zap.SugaredLogger, r *repository.Repo) *PersonUC {
	return &PersonUC{
		l: l,
		r: r,
	}
}

func (p *PersonUC) CreatePerson() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user entity.Person
		if err := ctx.ShouldBindJSON(&user); err != nil {
			p.l.Warnf("cannot parse person with err:%v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}
		id, err := p.r.CreatePerson(context.TODO(), &user)
		if err != nil {
			p.l.Warnf("cannot person user with err:%v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "person created", "Id": id})
	}
}

func (p *PersonUC) GetPersons() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sortBy := ctx.Query("sortBy")
		sortKey := ctx.Query("sortKey")
		if sortBy != "ASC" && sortBy != "DESC" && sortBy != "" {
			p.l.Warnf("cannot execute getPersons because sortBy is invalid")
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "sortBy is invalid"})
			return
		}
		if sortKey == "" {
			sortKey = "id"
		}
		users, err := p.r.GetPersons(context.TODO(), sortKey, sortBy)
		if err != nil {
			p.l.Warnf("cannot get person with err:%v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "cannot get person"})
			return
		}
		ctx.JSON(http.StatusOK, users)
	}
}

func (p *PersonUC) UpdatePerson() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := ctx.Param("id")
		if param == "" {
			p.l.Warnf("cannot update personapi, id is empty")
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "id must not be empty"})
			return
		}
		var user entity.Person
		if err := ctx.ShouldBindJSON(&user); err != nil {
			p.l.Warnf("cannot unmarshall personapi with err:%v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "cannot unmarshall personapi"})
			return
		}
		id, err := strconv.Atoi(param)
		if err != nil {
			p.l.Warnf("cannot parse id")
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "id is invalid"})
			return
		}
		err = p.r.UpdatePerson(context.TODO(), id, &user)
		if err != nil {
			p.l.Warnf("cannot update personapi with err:%v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "cannot update personapi"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func (p *PersonUC) DeletePerson() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := ctx.Param("id")
		if param == "" {
			p.l.Warnf("cannot update person, id is empty")
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "id must not be empty"})
			return
		}
		id, err := strconv.Atoi(param)
		if err != nil {
			p.l.Warnf("cannot parse id")
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "id is invalid"})
			return
		}
		err = p.r.DeletePerson(context.TODO(), id)
		if err != nil {
			p.l.Warnf("cannot delete personapi with err:%v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}
