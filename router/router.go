package router

import (
	"log"
	"net/http"

	"lifeAchieve/service"

	_ "lifeAchieve/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRouter() {
	router := mux.NewRouter()

	// USERS
	router.HandleFunc("/users/{id}", service.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/{id}", service.PatchUser).Methods("PATCH", "OPTIONS")

	// AUTH
	router.HandleFunc("/signup", service.SignUp).Methods("POST", "OPTIONS")
	router.HandleFunc("/signin", service.SignIn).Methods("POST", "OPTIONS")

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
