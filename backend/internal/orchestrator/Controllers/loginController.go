package Controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"216/internal/orchestrator/Database"
	"216/internal/orchestrator/Entities"
	"216/internal/orchestrator/Response"
	"216/internal/orchestrator/Services"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
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

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		Response.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := Entities.FormatError(err.Error())
		Response.ErrorResponse(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	type Result struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}
	tr := Result{Token: token, Email: user.Email}
	Response.JsonResponse(w, http.StatusOK, tr)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {

	if !Database.Instance.Migrator().HasTable("arithmetic_expressions") ||
		!Database.Instance.Migrator().HasTable("arithmetic_operations") ||
		!Database.Instance.Migrator().HasTable("computing_resources") ||
		!Database.Instance.Migrator().HasTable("users") {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func SignIn(email, password string) (string, error) {

	var err error

	user := Entities.User{}

	err = Database.Instance.Debug().Model(Entities.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = Entities.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	log.Println(user)
	return Services.CreateToken(user.ID)
}

func Profile(w http.ResponseWriter, r *http.Request) {

	user := Entities.User{}
	uid, _ := Services.ExtractTokenID(r)

	err := Database.Instance.Debug().Model(Entities.User{}).Where("id = ?", uid).Take(&user).Error

	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
