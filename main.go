package main

import (
	"net/http"
	"github.com/rs/cors"

	"github.com/ddld93/auth/controller"
	"github.com/ddld93/auth/routes"
	"github.com/gorilla/mux"
)


func main() {
	port := "5000"
	userCtrl := controller.NewUserCtrl("mongo", 27017)
	route := routes.UserRoute{UserCtrl: userCtrl}
	r := mux.NewRouter()

	
	// router handlers
	r.HandleFunc("/api/v1/auth/login", route.Login).Methods("POST")
	r.HandleFunc("/api/v1/auth/signup", route.CreateUser).Methods("POST")
	//r.HandleFunc("/user", route.GetUser).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("www")))
	

	// fmt.Printf("Server listening on port %v", port)
	// if err := http.ListenAndServe(":"+port, r); err != nil {
	// 	log.Fatal("Error starting server !! ", err)
	// }
	
 	// cors.Default() setup the middleware with default options being
    // all origins accepted with simple methods (GET, POST). See
    // documentation below for more options.
    handler := cors.Default().Handler(r)
    http.ListenAndServe(":"+port, handler)
}
