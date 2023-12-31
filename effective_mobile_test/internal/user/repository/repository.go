package repository

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	CreateUser() gin.HandlerFunc
	GetUsers() gin.HandlerFunc
	UpdateUser() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc
	EnrichById() gin.HandlerFunc
}

type Repository interface {
	UserRepository
}

type Repo struct {
	main    sql.DB
	replica sql.DB
}

func NewRepository(main *sql.DB, replica *sql.DB) *Repo {
	return &Repo{
		main:    *main,
		replica: *replica,
	}
}
