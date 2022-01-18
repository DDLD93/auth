package routes

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ddld93/auth/controller"
	"github.com/ddld93/auth/model"
	"github.com/ddld93/auth/utilities"
)

type UserRoute struct {
	UserCtrl *controller.DB_Connect
}
type CustomResponse struct {
	Message     string `json:"message"`
	Description string `json:"description"`
}



func (ur *UserRoute) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "Error Decoding request body"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	validatedUserModel,err:= utilities.UserModelValidate(&user)
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
	resp,err := ur.UserCtrl.CreateUser(validatedUserModel)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "error adding user to database"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}


func (ur *UserRoute) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "Error Decoding request body"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
	}
	regUser, err :=	ur.UserCtrl.GetUser(user.Email)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "A user with that email dont exist"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	}
	//comparing password 
	isValid:= utilities.CheckPasswordHash(user.Password ,regUser.Password)
	if !isValid {
		resp := CustomResponse{Message: "password did not match", Description: "wrong password input"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
	}
	token,err:=  utilities.TokenMaker.CreateToken(regUser.Email,regUser.Role,time.Hour)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "An error occured generating token"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token":token})
}




func (ur *UserRoute) GetUser(w http.ResponseWriter, r *http.Request) {
	reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]
		payload,err := utilities.TokenMaker.VerifyToken(reqToken)
		if err != nil{
			w.WriteHeader(http.StatusUnauthorized)
		}
		regUser, err :=	ur.UserCtrl.GetUser(payload.Username)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(regUser)

}






// func (ur *UserRoute) Login(email, password string) (string,error){
// 	user, err :=	ur.UserCtrl.GetUser(email)
// 	if err != nil {
// 		return "", err
// 	}
// 	//comparing password 
// 	resp:= utilities.CheckPasswordHash(password ,user.Password)
// 	if !resp {
// 		return "", errors.New("invalid password")
// 	}
// 	token,err:=  utilities.TokenMaker.CreateToken(user.Email,user.Role,time.Hour)
// 	if err != nil {
// 		return "", err
// 	}
// 	return token,nil
// }

// func (ur *UserRoute) Register(user *model.User) (string,error){
// 	validatedUserModel,err:= utilities.UserModelValidate(user)
// 	if err != nil {
// 		return "", err
// 	}
// 	// hashing password using Bcrypt
// 	passwordHash, err2 := utilities.HashPassword(validatedUserModel.Password)
// 	if err2 != nil {
// 		return "", errors.New(" error harshing password")
// 	}
// 	validatedUserModel.Password = passwordHash
// 	resp,err3 := ur.UserCtrl.CreateUser(validatedUserModel)
// 	if err3 != nil {
// 		return "", err3
// 	}
// 	return resp,nil
// }
