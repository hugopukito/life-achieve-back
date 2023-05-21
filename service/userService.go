package service

import (
	"database/sql"
	"lifeAchieve/cors"
	"lifeAchieve/entity"
	"lifeAchieve/repository"
	"log"
	"net/http"

	"github.com/clarketm/json"

	"github.com/gorilla/mux"
)

// GetUser godoc
// @Tags users
// @Summary Get user
// @Produce  json
// @Param Authorization header string true "JWT"
// @Param id path int true "User id"
// @Success 200 {array} entity.User
// @Failure 400
// @Router /users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	if cors.EnableCors(&w, r) == "options" {
		return
	}

	userId := parseAuthorization(w, r)
	id := mux.Vars(r)["id"]

	if userId != id {
		log.Println("userId from token doesn't match id in parameters")
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	user, err := repository.FindUserById(id)
	if err == sql.ErrNoRows {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dto, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte(string(dto)))
}

func PostUser(w http.ResponseWriter, user entity.PostUser) {
	err := repository.InsertUser(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

// PatchUser godoc
// @Tags users
// @Summary Patch user
// @Accept  json
// @Param Authorization header string true "JWT"
// @Param id path int true "User id"
// @Param user body entity.User true "User object"
// @Success 204
// @Failure 400
// @Router /users/{id} [patch]
func PatchUser(w http.ResponseWriter, r *http.Request) {
	if cors.EnableCors(&w, r) == "options" {
		return
	}

	userId := parseAuthorization(w, r)
	id := mux.Vars(r)["id"]

	if userId != id {
		log.Println("userId from token doesn't match id in parameters")
		w.WriteHeader(http.StatusBadRequest)
	}

	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = repository.UpdateUser(id, user)
	if err != nil {
		if err.Error() == "no changes" {
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetUserWithEmail(w http.ResponseWriter, email string) (*entity.User, error) {
	user, err := repository.FindUserByEmail(email)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
		}
		return nil, err
	}

	return user, nil
}

func GetUserTypeById(w http.ResponseWriter, id string) (string, error) {
	userType, err := repository.FindUserTypeById(id)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
		}
		return "", err
	}

	return userType, nil
}

func EmailExist(email string) (bool, error) {
	_, err := repository.FindUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
