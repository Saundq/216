package Controllers

import (
	"216/internal/orchestrator/Database"
	"216/internal/orchestrator/Entities"
	"216/internal/orchestrator/Response"
	"216/internal/orchestrator/Services"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		Response.ErrorResponse(w, http.StatusUnprocessableEntity, err)
	}
	user := Entities.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		Response.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		Response.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(Database.Instance)

	if err != nil {

		formattedError := Entities.FormatError(err.Error())

		Response.ErrorResponse(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	Response.JsonResponse(w, http.StatusCreated, userCreated)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	user := Entities.User{}

	users, err := user.FindAllUsers(Database.Instance)
	if err != nil {
		Response.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	Response.JsonResponse(w, http.StatusOK, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		Response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	user := Entities.User{}
	userGotten, err := user.FindUserByID(Database.Instance, uint32(uid))
	if err != nil {
		Response.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	Response.JsonResponse(w, http.StatusOK, userGotten)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	//uid, err := strconv.ParseUint(vars["id"], 10, 32)
	uid := vars["id"]
	// if err != nil {
	// 	Response.ErrorResponse(w, http.StatusBadRequest, err)
	// 	return
	// }
	body, err := io.ReadAll(r.Body)
	if err != nil {
		Response.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := Entities.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		Response.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := Services.ExtractTokenID(r)
	if err != nil {
		Response.ErrorResponse(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	uuid_r, err := uuid.FromString(uid)
	if tokenID != &uuid_r {
		Response.ErrorResponse(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		Response.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateAUser(Database.Instance, uuid_r)
	if err != nil {
		formattedError := Entities.FormatError(err.Error())
		Response.ErrorResponse(w, http.StatusInternalServerError, formattedError)
		return
	}
	Response.JsonResponse(w, http.StatusOK, updatedUser)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	user := Entities.User{}

	//uid, err := strconv.ParseUint(vars["id"], 10, 32)
	uid := vars["id"]
	// if err != nil {
	// 	Response.ErrorResponse(w, http.StatusBadRequest, err)
	// 	return
	// }
	//var tokenID *uuid.UUID
	tokenID, err := Services.ExtractTokenID(r)
	uuid_r, err := uuid.FromString(uid)
	if err != nil {
		Response.ErrorResponse(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != nil && tokenID != &uuid_r {
		Response.ErrorResponse(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteUser(Database.Instance, uuid_r)
	if err != nil {
		Response.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	Response.JsonResponse(w, http.StatusNoContent, "")
}
