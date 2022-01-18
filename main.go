package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ddld93/auth/controller"
	"github.com/ddld93/auth/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)
type MiddlewareFunc func(http.Handler) http.Handler

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
}
}

func main()  {
	port := os.Getenv("PORT")
	userCtrl := controller.NewUserCtrl("localhost", 27017)
	route := routes.UserRoute{UserCtrl: userCtrl}

	r := mux.NewRouter()
	
    r.HandleFunc("/login",route.Login ).Methods("POST")
    r.HandleFunc("/signup",route.CreateUser ).Methods("POST")
	r.HandleFunc("/user", route.GetUser).Methods("GET") 
	
	r.Use(mux.CORSMethodMiddleware(r))



    http.Handle("/", r)

	fmt.Printf("Server listening on port %v", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Error starting server !! ", err)
	}


}