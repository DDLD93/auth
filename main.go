package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ddld93/auth/controller"
	"github.com/ddld93/auth/routes"

	"github.com/gorilla/mux"
)

type MiddlewareFunc func(http.Handler) http.Handler



func main() {
	port := "5000"
	userCtrl := controller.NewUserCtrl("mongo", 27017)
	route := routes.UserRoute{UserCtrl: userCtrl}
	r := mux.NewRouter()

	// router handlers
	r.HandleFunc("/api/v1/auth/login", route.Login).Methods("POST")
	r.HandleFunc("/api/v1/auth/signup", route.CreateUser).Methods("POST")
	//r.HandleFunc("/user", route.GetUser).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./www/")))
	r.Use(mux.CORSMethodMiddleware(r))

	fmt.Printf("Server listening on port %v", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Error starting server !! ", err)
	}

}
