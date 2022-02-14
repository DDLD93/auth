package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ddld93/auth/controller"
	"github.com/ddld93/auth/model"
	"github.com/ddld93/auth/utilities"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type UserRoute struct {
	UserCtrl *controller.DB_Connect
	Session *websocket.Conn
}
type SocketConn struct{
	
}
type CustomResponse struct {
	Message     string `json:"message"`
	Description string `json:"description"`
}
type UserResponse struct {
	Status string       `json:"status"`
	Token  string     `json:"token"`
	User   model.User `json:"user"`
}
func NewSocket(r *http.Request,w http.ResponseWriter) (UserRoute, error)  {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	session, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return UserRoute{},err
    }
	return UserRoute{Session: session}, nil
}
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
    for {
    // read in a message
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
    // print out that message for clarity
        fmt.Println(string(p))

        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
        }

    }
}


func (ur *UserRoute) WsEndpoint(w http.ResponseWriter, r *http.Request) {
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    // upgrade this connection to a WebSocket
    // connection
    conn, err := NewSocket(r, w)
    if err != nil {
        log.Println(err)
    }
	err =conn.Session.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
        log.Println(err)
    }
    log.Println("Client Connected")

    reader(conn.Session)
    // listen indefinitely for new messages coming
    // through on our WebSocket connection
}

func (ur *UserRoute) CreateUser(w http.ResponseWriter, r *http.Request) {

	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "Error Decoding request body"}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	validatedUserModel, err := utilities.UserModelValidate(&user)
	if err != nil {
		resp := CustomResponse{Message:"invalid input fields" , Description:err.Error()}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// hashing password using Bcrypt
	passwordHash, err := utilities.HashPassword(validatedUserModel.Password)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "internal server error"}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	validatedUserModel.Password = passwordHash
	_, err = ur.UserCtrl.CreateUser(validatedUserModel)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "error adding user to database"}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Succuess",
		"message": "new account ceated",
	})
}

func (ur *UserRoute) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "Error Decoding request body"}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	regUser, err := ur.UserCtrl.GetUser(user.Email)

	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "A user with that email dont exist"}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//comparing password
	isValid := utilities.CheckPasswordHash(user.Password, regUser.Password)
	if !isValid {
		resp := CustomResponse{Message: "password did not match", Description: "wrong password input"}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//testing socket connection

	token, err := utilities.TokenMaker.CreateToken(regUser.Email, regUser.Role, time.Hour)
	if err != nil {
		resp := CustomResponse{Message: err.Error(), Description: "An error occured generating token"}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	regUser.Password = ""
	regUser.Role = ""
	response := UserResponse{Status: "Login Success", Token: token, User: *regUser}
	
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// initializing web sockets
	
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
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// checking to token has admin previllages

	if payload.AccoutType != "client"{
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message":"not authorize to make such request"})
		return
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
		json.NewEncoder(w).Encode("payment verification failed")
	return
	}

	err = ur.UserCtrl.UpdatePayment(payload.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user,err := ur.UserCtrl.GetUser(payload.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = ""
	user.Role = ""

	resp := CustomResponse{
		Message:     "payment ok",
		Description: "payment verified succesifully",
	}
	err =ur.Session.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
        log.Println(err)
    }
	json.NewEncoder(w).Encode(resp)

}
