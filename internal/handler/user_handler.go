package handler

import (
	"github.com/core-go/core"
	"github.com/core-go/core/search"
	s "github.com/core-go/search"
	"net/http"
	"reflect"

	. "go-service/internal/filter"
	. "go-service/internal/model"
	. "go-service/internal/service"
)

func NewUserHandler(service UserService, logError core.Log, validate core.Validate, action *core.ActionConfig) *UserHandler {
	modelType := reflect.TypeOf(User{})
	params := core.CreateParams(modelType, logError, validate, action)
	filterType := reflect.TypeOf(UserFilter{})
	paramIndex, filterIndex := s.BuildParams(filterType)
	return &UserHandler{service: service, Params: params, paramIndex: paramIndex, filterIndex: filterIndex}
}

type UserHandler struct {
	service UserService
	*search.SearchHandler
	*core.Params
	paramIndex  map[string]int
	filterIndex int
}

func (h *UserHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := core.GetRequiredParam(w, r)
	if len(id) > 0 {
		user, err := h.service.Load(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		core.JSON(w, core.IsFound(user), user)
	}
}
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := core.Decode(w, r, &user)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !core.HasError(w, r, errors, er2, h.Error, h.Log, h.Resource, h.Action.Create) {
			res, er3 := h.service.Create(r.Context(), &user)
			core.AfterCreated(w, r, &user, res, er3, h.Error, h.Log, h.Resource, h.Action.Create)
		}
	}
}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := core.DecodeAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !core.HasError(w, r, errors, er2, h.Error, h.Log, h.Resource, h.Action.Update) {
			res, er3 := h.service.Update(r.Context(), &user)
			core.HandleResult(w, r, &user, res, er3, h.Error, h.Log, h.Resource, h.Action.Update)
		}
	}
}
func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var user User
	r, json, er1 := core.BuildMapAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !core.HasError(w, r, errors, er2, h.Error, h.Log, h.Resource, h.Action.Patch) {
			res, er3 := h.service.Patch(r.Context(), json)
			core.HandleResult(w, r, json, res, er3, h.Error, h.Log, h.Resource, h.Action.Patch)
		}
	}
}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := core.GetRequiredParam(w, r)
	if len(id) > 0 {
		res, err := h.service.Delete(r.Context(), id)
		core.HandleDelete(w, r, res, err, h.Error, h.Log, h.Resource, h.Action.Delete)
	}
}
func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	filter := UserFilter{Filter: &s.Filter{}}
	s.Decode(r, &filter, h.paramIndex, h.filterIndex)

	var users []User
	users, total, err := h.service.Search(r.Context(), &filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	core.JSON(w, http.StatusOK, &s.Result{List: &users, Total: total})
}
