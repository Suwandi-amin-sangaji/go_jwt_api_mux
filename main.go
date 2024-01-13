package main

import (
	"go-jwt-api/controllers/authcontroller"
	"go-jwt-api/controllers/productcontroller"
	"go-jwt-api/middlewares"
	"go-jwt-api/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	models.ConnectDatabase()
	r := mux.NewRouter()

	r.HandleFunc("/login", authcontroller.Login).Methods("POST")
	r.HandleFunc("/Register", authcontroller.Register).Methods("POST")
	r.HandleFunc("/logout", authcontroller.Logout).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/products", productcontroller.Index).Methods("GET")
	api.Use(middlewares.JWTmiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}
