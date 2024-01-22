package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-jun/internal/person/entity"
	"go-jun/internal/person/personapi"
)

func (r *Repo) CreatePerson(ctx context.Context, person *entity.Person) (id int, err error) {
	age, gender, nationality, err := personapi.CallExternalAPI(person.Name)
	if err != nil {
		return 0, fmt.Errorf("cannot take data from site: %v", err)
	}
	err = r.main.QueryRow("insert into person(name,surname,patronymic, age,gender,nationality) values($1, $2, $3,$4,$5,$6) returning ID", person.Name, person.Surname, person.Patronymic, age, gender, nationality).Scan(&person.Id) // $1 and $2 is prepared statement
	if err != nil {
		return 0, fmt.Errorf("cannot query with error: %v", err)
	}
	return person.Id, nil
}

func (r *Repo) UpdatePerson(ctx context.Context, id int, person *entity.Person) error {
	_, err := r.main.Exec("UPDATE person SET name=$1, surname=$2, patronymic=$3 WHERE id=$4", person.Name, person.Surname, person.Patronymic, id)
	if err != nil {
		return fmt.Errorf("cannot query with err:%v", err)
	}
	return nil
}

func (r *Repo) DeletePerson(ctx context.Context, id int) error {
	_, err := r.main.Exec("DELETE from person WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("cannot delete person with err:%v", err)
	}
	return nil
}

func (r *Repo) GetPersons(ctx context.Context, sortKey, sortBy string) ([]entity.Person, error) {
	persons := make([]entity.Person, 0)
	rows, err := r.replica.Query("SELECT * from person")
	if err != nil {
		return nil, fmt.Errorf("cannot query with error: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
	for rows.Next() {
		var person entity.Person
		if err = rows.Scan(&person.Id, &person.Name, &person.Surname, &person.Patronymic, &person.Age, &person.Gender, &person.Nationality); err != nil {
			return nil, fmt.Errorf("cannot scan query with error: %v", err)
		}
		persons = append(persons, person)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows with error: %v", err)
	}
	return persons, nil
}
