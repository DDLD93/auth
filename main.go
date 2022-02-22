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
	r.HandleFunc("/new/admin", route.CreateAdmin).Methods("POST")
	r.HandleFunc("/login", route.Login).Methods("POST")
	r.HandleFunc("/signup", route.CreateUser).Methods("POST")
	r.HandleFunc("/users", route.GetUsers).Methods("GET")
	r.HandleFunc("/user/{email}", route.GetUser).Methods("GET")
	r.HandleFunc("/analytics", route.GetUsersAnalytics).Methods("GET")
	r.HandleFunc("/paystack/{reference}/{email}", route.Verify).Methods("GET")
	r.HandleFunc("/formflag/{email}", route.FormFlag).Methods("GET")
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
