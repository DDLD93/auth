package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"

	"github.com/ddld93/auth/controller"
	"github.com/ddld93/auth/routes"
	"github.com/gorilla/mux"
)


func main() {
	port := "5000"
	userCtrl := controller.NewUserCtrl("localhost", 27017)
	route := routes.UserRoute{UserCtrl: userCtrl}
	r := mux.NewRouter() 

	
	// router handlers
	r.HandleFunc("/ws", route.WsEndpoint).Methods("GET")
	r.HandleFunc("/api/v1/auth/login", route.Login).Methods("POST")
	r.HandleFunc("/api/v1/auth/signup", route.CreateUser).Methods("POST")
	r.HandleFunc("/api/v1/auth/users", route.GetUsers).Methods("GET")
	r.HandleFunc("/api/v1/auth/users/analytics", route.GetUsersCount).Methods("GET")
	r.HandleFunc("/api/paystack/verify/{reference}", route.Verify).Methods("GET")
	r.HandleFunc("/api/auth/formflag/{useremail}", route.FormFlag).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("www")))
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: true,
		Debug: false,
		
	})

    handler := c.Handler(r)
    
	fmt.Printf("Server listening on port %v", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Error starting server !! ", err)
	}
}
