package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ddld93/auth/controller"
	"github.com/ddld93/auth/model"
	"github.com/ddld93/auth/utilities"
	"github.com/gorilla/mux"
)

type UserRoute struct {
	UserCtrl *controller.DB_Connect
}
type CustomResponse struct {
	Message     string `json:"message"`
	Description string `json:"description"`
}
type UserResponse struct {
	Status int        `json:"status"`
	Token  string     `json:"token"`
	User   model.User `json:"user"`
}

func (ur *UserRoute) CreateUser(w http.ResponseWriter, r *http.Request) {

	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "Error Decoding request body"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	validatedUserModel, err := utilities.UserModelValidate(&user)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "invalid form input"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	// hashing password using Bcrypt
	passwordHash, err := utilities.HashPassword(validatedUserModel.Password)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "internal server error"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	validatedUserModel.Password = passwordHash
	resp, err := ur.UserCtrl.CreateUser(validatedUserModel)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "error adding user to database"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (ur *UserRoute) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "Error Decoding request body"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	regUser, err := ur.UserCtrl.GetUser(user.Email)

	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "A user with that email dont exist"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	//comparing password
	isValid := utilities.CheckPasswordHash(user.Password, regUser.Password)
	if !isValid {
		resp := CustomResponse{Message: "password did not match", Description: "wrong password input"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}
	token, err := utilities.TokenMaker.CreateToken(regUser.Email, regUser.Role, time.Hour)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "An error occured generating token"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	regUser.Password = ""
	regUser.Role = ""
	response := UserResponse{Status: http.StatusCreated, Token: token, User: *regUser}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (ur *UserRoute) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	reqToken := r.Header.Get("Authorization")
	// checking if request carries a valid token
	if reqToken == "" {
		resp := CustomResponse{
			Message:     "Token not Found",
			Description: "Bearer token not included in request"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	token := splitToken[1]
	payload, err := utilities.TokenMaker.VerifyToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	// checking to token has admin previllages
	if payload.AccoutType != "client" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "not authorize to make such request"})
	}

	regUser, err := ur.UserCtrl.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(regUser)

}

func (ur *UserRoute) Verify(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	reqToken := r.Header.Get("Authorization")
	// checking if request carries a valid token
	if reqToken == "" {
		resp := CustomResponse{
			Message:     "Token not Found",
			Description: "Bearer token not included in request"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	token := splitToken[1]
	payload, err := utilities.TokenMaker.VerifyToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	refNo := mux.Vars(r)
	reference := refNo["reference"]
	fmt.Println(reference)
	err = utilities.VerifyPayment(reference)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Println(err)
		return
	}

	err = ur.UserCtrl.UpdatePayment(payload.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := CustomResponse{
		Message:     "payment ok",
		Description: "payment verified succesifully",
	}
	json.NewEncoder(w).Encode(resp)

}
