package service

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/clarketm/json"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"lifeAchieve/cors"
	"lifeAchieve/entity"
)

// SignUp godoc
// @Tags auth
// @Summary Post user
// @Accept  json
// @Param user body entity.PostUser true "User object"
// @Success 201
// @Failure 400
// @Router /signup [post]
func SignUp(w http.ResponseWriter, r *http.Request) {
	if cors.EnableCors(&w, r) == "options" {
		return
	}

	var user entity.PostUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	emailExist, err := EmailExist(user.Email)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	if emailExist {
		w.WriteHeader(http.StatusConflict)
		return
	}

	pwd, err := generateHashPassword(user.Password)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	user.Password = pwd
	PostUser(w, user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusCreated)
}

// SignIn godoc
// @Tags auth
// @Summary Get jwt token
// @Accept  json
// @Param authentication body entity.Authentication true "Authentication object"
// @Produce  json
// @Success 200 {string} string
// @Failure 400
// @Router /signin [post]
func SignIn(w http.ResponseWriter, r *http.Request) {
	if cors.EnableCors(&w, r) == "options" {
		return
	}

	var authdetails entity.Authentication
	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authUser, err := GetUserWithEmail(w, authdetails.Email)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	check := checkPasswordHash(authdetails.Password, authUser.Password)

	if !check {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	validToken, err := generateJWT(*authUser)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(validToken)
}

func parseAuthorization(w http.ResponseWriter, r *http.Request) string {
	token := r.Header["Authorization"]
	var userId string
	if len(token) > 0 {
		claims, err := ParseJwt(w, token[0])
		if err != nil {
			log.Println("Problem with jwt parsing")
			w.WriteHeader(http.StatusBadRequest)
		}
		userId = claims["userId"].(string)
	}
	if userId == "" {
		log.Println("No userId in token after parsing")
		w.WriteHeader(http.StatusBadRequest)
	}
	return userId
}

func generateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJWT(user entity.User) (string, error) {

	path, err := os.Getwd()
	if err != nil {
		log.Println("error in retrieving path dir")
		return "", err
	}

	secret_jwt, err := os.ReadFile(path + "/secret_jwt.txt")
	if err != nil {
		log.Println("error in read secret file")
		return "", err
	}

	var mySigningKey = []byte(secret_jwt)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["userId"] = user.Id
	claims["firstName"] = user.FirstName
	claims["lastName"] = user.LastName
	claims["exp"] = time.Now().Add(time.Minute * 60 * 24 * 365).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", errors.New("something Went Wrong: " + err.Error())
	}
	return tokenString, nil
}

func ParseJwt(w http.ResponseWriter, bearerToken string) (jwt.MapClaims, error) {
	bearerToken = strings.Replace(bearerToken, "Bearer ", "", 1)

	path, err := os.Getwd()
	if err != nil {
		log.Println("error in retrieving path dir")
		return nil, err
	}

	secret_jwt, err := os.ReadFile(path + "/secret_jwt.txt")
	if err != nil {
		log.Println("error in read secret file")
		return nil, err
	}

	var mySigningKey = []byte(secret_jwt)

	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		return mySigningKey, nil
	})

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	w.WriteHeader(http.StatusBadRequest)
	return nil, err
}
