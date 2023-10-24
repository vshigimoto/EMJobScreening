package repository

import (
	"database/sql"
	"effective_mobile_test/internal/user/entity"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
)

func (r *Repo) CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user entity.User
		err := ctx.ShouldBindJSON(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
		//user.Name and user.Email is arguments for prepared statement ($1 and $2)
		err = r.main.QueryRow("INSERT INTO people(name, surname, patronymic, age, sex, nationality) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", user.Name, user.Surname, user.Patronymic, user.Age, user.Sex, user.Nationality).Scan(&user.Id)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"id": user.Id})
	}
}

func (r *Repo) GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users := make([]entity.User, 0)

		rows, err := r.replica.Query("SELECT * from people")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {

			}
		}(rows)
		for rows.Next() {
			var user entity.User
			if err = rows.Scan(&user.Id, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Sex, &user.Nationality); err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				return
			}
			users = append(users, user)
		}
		if err = rows.Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, users)
	}
}

func (r *Repo) UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "id should not to be empty"})
			return
		}
		var user entity.User
		// Unmarshal json to a new user
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}
		_, err := r.main.Exec("UPDATE people SET name=$1, surname=$2, patronymic=$3, age=$4, sex=$5, nationality=$6 WHERE id=$7", user.Name, user.Surname, user.Patronymic, user.Age, user.Sex, user.Nationality, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}

func (r *Repo) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "id should not to be empty"})
			return
		}
		_, err := r.main.Exec("DELETE from people WHERE id=$1", id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}

func (r *Repo) EnrichById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "id should not to be empty"})
			return
		}
		rows, err := r.replica.Query("SELECT * FROM people WHERE id=$1", id)
		var user entity.User
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {

			}
		}(rows)
		rows.Next()
		if err = rows.Scan(&user.Id, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Sex, &user.Nationality); err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ageUrl := "https://api.agify.io/?name=" + user.Name
		sexUrl := "https://api.genderize.io/?name=" + user.Name
		nationUrl := "https://api.nationalize.io/?name=" + user.Name

		values := []string{ageUrl, sexUrl, nationUrl}
		var results []string
		for _, value := range values {
			res, err := http.Get(value)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(res.Body)
			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				return
			}
			results = append(results, string(resBody))
		}
		var age entity.AgifyResult
		var sex entity.GenderizeResult
		var nation entity.NationalizeResult
		if err := json.Unmarshal([]byte(results[0]), &age); err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		if err := json.Unmarshal([]byte(results[1]), &sex); err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		if err := json.Unmarshal([]byte(results[2]), &nation); err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		user.Age = age.Age
		user.Sex = sex.Gender
		user.Nationality = nation.Country[0].CountryID
		_, err = r.main.Exec("UPDATE people SET name=$1, surname=$2, patronymic=$3, age=$4, sex=$5, nationality=$6 WHERE id=$7", user.Name, user.Surname, user.Patronymic, user.Age, user.Sex, user.Nationality, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
		}

		ctx.JSON(http.StatusOK, user)
	}
}
