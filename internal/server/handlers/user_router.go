package handlers

import (
	"encoding/json"
	"fmt"
	"go-postgres/pkg/response"
	"go-postgres/pkg/user"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type UserRouter struct {
	Repository user.Repository
}

/*
	CREATE USER
*/

func (ur *UserRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var u user.User

	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = ur.Repository.Create(ctx, &u)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	u.Password = ""
	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), u.ID))
	response.JSON(w, r, http.StatusCreated, response.Map{"user": u})
}

/*
	LIST USER
*/

func (ur *UserRouter) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := r.URL.Query().Get("username")

	if username != "" {
		user, err := ur.Repository.GetByUsername(ctx, username)

		if err != nil {
			response.HTTPError(w, r, http.StatusNotFound, err.Error())
			return
		}

		response.JSON(w, r, http.StatusOK, response.Map{"user": user})
		return

	}

	users, err := ur.Repository.GetAll(ctx)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"users": users})
}

/*
	GET USER BY ID
*/

func (ur *UserRouter) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	u, err := ur.Repository.GetOne(ctx, uint(id))
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"user": u})
}

/*
	UPDATE USER
*/

func (ur *UserRouter) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var u user.User
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	_user, err := ur.Repository.Update(ctx, uint(id), u)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"user": _user})
}

/*
	DELETE USERS

*/

func (ur *UserRouter) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	err = ur.Repository.Delete(ctx, uint(id))
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{})
}

//ROUTERS

func (ur *UserRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", ur.GetAllHandler)
	r.Post("/", ur.CreateHandler)
	r.Get("/{id}/", ur.GetOneHandler)
	r.Put("/{id}/", ur.UpdateHandler)
	r.Delete("/{id}/", ur.DeleteHandler)

	return r
}
