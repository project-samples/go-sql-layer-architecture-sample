package repository

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/core-go/search/query"
	q "github.com/core-go/sql"

	. "go-service/internal/filter"
	. "go-service/internal/model"
)

func NewUserAdapter(db *sql.DB, buildQuery func(*UserFilter) (string, []interface{})) (*UserAdapter, error) {
	userType := reflect.TypeOf(User{})
	if buildQuery == nil {
		userQueryBuilder := query.NewBuilder(db, "users", userType)
		buildQuery = func(filter *UserFilter) (s string, i []interface{}) {
			return userQueryBuilder.BuildQuery(filter)
		}
	}
	params, err := q.CreateParams(userType, db)
	if err != nil {
		return nil, err
	}
	return &UserAdapter{DB: db, Params: params, BuildQuery: buildQuery}, nil
}

type UserAdapter struct {
	DB         *sql.DB
	BuildQuery func(*UserFilter) (string, []interface{})
	*q.Params
}

func (r *UserAdapter) Load(ctx context.Context, id string) (*User, error) {
	var users []User
	query := fmt.Sprintf("select %s from users where id = %s limit 1", r.Fields, r.BuildParam(1))
	err := q.Query(ctx, r.DB, r.Map, &users, query, id)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return &users[0], nil
	}
	return nil, nil
}

func (r *UserAdapter) Create(ctx context.Context, user *User) (int64, error) {
	query, args := q.BuildToInsert("users", user, r.BuildParam)
	res, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, nil
	}
	return res.RowsAffected()
}

func (r *UserAdapter) Update(ctx context.Context, user *User) (int64, error) {
	query, args := q.BuildToUpdate("users", user, r.BuildParam)
	res, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, nil
	}
	return res.RowsAffected()
}

func (r *UserAdapter) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	colMap := q.JSONToColumns(user, r.JsonColumnMap)
	query, args := q.BuildToPatch("users", colMap, r.Keys, r.BuildParam)
	res, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *UserAdapter) Delete(ctx context.Context, id string) (int64, error) {
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

func (r *UserAdapter) Search(ctx context.Context, filter *UserFilter) ([]User, int64, error) {
	var users []User
	if filter.Limit <= 0 {
		return users, 0, nil
	}
	query, params := r.BuildQuery(filter)
	offset := q.GetOffset(filter.Limit, filter.Page)
	pagingQuery := q.BuildPagingQuery(query, filter.Limit, offset)
	countQuery := q.BuildCountQuery(query)

	total, err := q.Count(ctx, r.DB, countQuery, params...)
	if err != nil {
		return users, 0, err
	}
	err = q.Query(ctx, r.DB, r.Map, &users, pagingQuery, params...)
	return users, total, err
}

func BuildQuery(filter *UserFilter) (string, []interface{}) {
	query := "select * from users"
	where, params := BuildFilter(filter)
	if len(where) > 0 {
		query = query + " where " + where
	}
	return query, params
}
func BuildFilter(filter *UserFilter) (string, []interface{}) {
	buildParam := q.BuildDollarParam
	var where []string
	var params []interface{}
	i := 1
	if len(filter.Id) > 0 {
		params = append(params, filter.Id)
		where = append(where, fmt.Sprintf(`id = %s`, buildParam(i)))
		i++
	}
	if filter.DateOfBirth != nil {
		if filter.DateOfBirth.Min != nil {
			params = append(params, filter.DateOfBirth.Min)
			where = append(where, fmt.Sprintf(`date_of_birth >= %s`, buildParam(i)))
			i++
		}
		if filter.DateOfBirth.Max != nil {
			params = append(params, filter.DateOfBirth.Max)
			where = append(where, fmt.Sprintf(`date_of_birth <= %s`, buildParam(i)))
			i++
		}
	}
	if len(filter.Username) > 0 {
		q := filter.Username + "%"
		params = append(params, q)
		where = append(where, fmt.Sprintf(`username like %s`, buildParam(i)))
		i++
	}
	if len(filter.Email) > 0 {
		q := filter.Email + "%"
		params = append(params, q)
		where = append(where, fmt.Sprintf(`email like %s`, buildParam(i)))
		i++
	}
	if len(filter.Phone) > 0 {
		q := "%" + filter.Phone + "%"
		params = append(params, q)
		where = append(where, fmt.Sprintf(`phone like %s`, buildParam(i)))
		i++
	}
	if len(where) > 0 {
		return strings.Join(where, " and "), params
	}
	return "", params
}
