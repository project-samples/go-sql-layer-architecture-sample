package repository

import (
	"context"
	"database/sql"
	"fmt"
	q "github.com/core-go/sql"
	"reflect"

	. "go-service/internal/model"
)

type UserRepository interface {
	Load(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) (int64, error)
	Update(ctx context.Context, user *User) (int64, error)
	Patch(ctx context.Context, user map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

type userRepository struct {
	DB *sql.DB
}

func (r *userRepository) Load(ctx context.Context, id string) (*User, error) {
	var users []User
	query := fmt.Sprintf("select id, username, email, phone, date_of_birth from users where id = %s limit 1", q.BuildParam(1))
	err := q.Query(ctx, r.DB, nil, &users, query, id)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return &users[0], nil
	}
	return nil, nil
}

func (r *userRepository) Create(ctx context.Context, user *User) (int64, error) {
	query, args := q.BuildToInsert("users", user, q.BuildParam)
	res, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, nil
	}
	return res.RowsAffected()
}

func (r *userRepository) Update(ctx context.Context, user *User) (int64, error) {
	query, args := q.BuildToUpdate("users", user, q.BuildParam)
	res, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, nil
	}
	return res.RowsAffected()
}

func (r *userRepository) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	userType := reflect.TypeOf(User{})
	jsonColumnMap := q.MakeJsonColumnMap(userType)
	colMap := q.JSONToColumns(user, jsonColumnMap)
	keys, _ := q.FindPrimaryKeys(userType)
	query, args := q.BuildToPatch("users", colMap, keys, q.BuildParam)
	res, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *userRepository) Delete(ctx context.Context, id string) (int64, error) {
	query := "delete from users where id = ?"
	stmt, er0 := r.DB.Prepare(query)
	if er0 != nil {
		return -1, nil
	}
	res, er1 := stmt.ExecContext(ctx, id)
	if er1 != nil {
		return -1, er1
	}
	return res.RowsAffected()
}
