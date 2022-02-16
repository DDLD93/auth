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
	Status     string `json:"status"`
	Message string `json:"message"`
}
type PaymentResponse struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	User   model.User `json:"user"`
}
type UserCount struct {
	Total int       `json:"total"`
	TotalPaid  int  `json:"totalPaid"`
}
type UserResponse struct {
	Status string      `json:"status"`
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
	w.Header().Set("Content-Type", "application/json")
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp := CustomResponse{Status: err.Error(), Message: "Error Decoding request body"}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	validatedUserModel, err := utilities.UserModelValidate(&user)
	if err != nil {
		resp := CustomResponse{Status:"failed" , Message:err.Error()}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// hashing password using Bcrypt
	passwordHash, err := utilities.HashPassword(validatedUserModel.Password)
	if err != nil {
		resp := CustomResponse{Status: "failed", Message: "internal server error"}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	validatedUserModel.Password = passwordHash
	err = ur.UserCtrl.CreateUser(validatedUserModel)
	if err != nil {
		resp := CustomResponse{Status: "failed", Message: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	w.WriteHeader(http.StatusOK)
	
	
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
		resp := CustomResponse{Status: err.Error(), Message: "Error Decoding request body"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	regUser, err := ur.UserCtrl.GetUser(user.Email)

	if err != nil {
		 resp := CustomResponse{Status:"failed", Message: "A user with that email dont exist"}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("A user with that email dont exist")
		json.NewEncoder(w).Encode(resp)
		return
	}
	//comparing password
	isValid := utilities.CheckPasswordHash(user.Password, regUser.Password)
	if !isValid {
		resp := CustomResponse{Status: "failed", Message: "wrong password input"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}
	//testing socket connection

	token, err := utilities.TokenMaker.CreateToken(regUser.Email, regUser.Role, time.Hour)
	if err != nil {
		resp := CustomResponse{Status: "failed", Message: "An error occured generating token"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	regUser.Password = ""
	regUser.Role = ""
	response := UserResponse{Status: "Login Success", Token: token, User: *regUser}
	
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	

 }

func (ur *UserRoute) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	reqToken := r.Header.Get("Authorization")
	// checking if request carries a valid token
	if reqToken == "" {
		resp := CustomResponse{
			Status:     "failed",
			Message: "Bearer token not included in request"}
		json.NewEncoder(w).Encode(resp)
		
		return
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	token := splitToken[1]
	payload, err := utilities.TokenMaker.VerifyToken(token)
	if err != nil {
		resp := CustomResponse{Status: "failed", Message: "invalid token"}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// checking to token has admin previllages

	if payload.AccoutType != "admin"{
		w.WriteHeader(http.StatusUnauthorized)
		resp := CustomResponse{Status: "failed", Message: "Not authorize to make such request"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	regUser, err := ur.UserCtrl.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := CustomResponse{Status: "failed", Message: "Not authorize to make such request"}
		json.NewEncoder(w).Encode(resp)
	}
	json.NewEncoder(w).Encode(regUser)

}

func (ur *UserRoute) Verify(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	reqToken := r.Header.Get("Authorization")
	// checking if request carries a valid token
	if reqToken == "" {
		resp := CustomResponse{
			Status:     "failed",
			Message: "Bearer token not included in request",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	token := splitToken[1]
	payload, err := utilities.TokenMaker.VerifyToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		resp := CustomResponse{
			Status:     "failed",
			Message: "Invalid token",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	refNo := mux.Vars(r)
	reference := refNo["reference"]
	fmt.Println(reference)
	err = utilities.VerifyPayment(reference)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		resp := CustomResponse{
			Status:     "failed",
			Message: "Payment not verfied",
		}
		json.NewEncoder(w).Encode(resp)
	return
	}

	err = ur.UserCtrl.UpdatePayment(payload.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := CustomResponse{
			Status:     "failed",
			Message: "Something went wrong",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	user,err := ur.UserCtrl.GetUser(payload.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := CustomResponse{
			Status:     "failed",
			Message: "Something went wrong",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	user.Password = ""
	user.Role = ""

	json.NewEncoder(w).Encode(PaymentResponse{
		Message: "Success",
		User: *user,
	})

}
func (ur *UserRoute) GetUsersCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	reqToken := r.Header.Get("Authorization")
	// checking if request carries a valid token
	if reqToken == "" {
		resp := CustomResponse{
			Status:     "failed",
			Message: "Bearer token not included in request"}
		json.NewEncoder(w).Encode(resp)
		
		return
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	token := splitToken[1]
	_, err := utilities.TokenMaker.VerifyToken(token)
	if err != nil {
		resp := CustomResponse{Status: "failed", Message: "invalid token"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}
	// checking to token has admin previllages

	// if payload.AccoutType != "client"{
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	resp := CustomResponse{Status: "failed", Message: "Not authorize to make such request"}
	// 	json.NewEncoder(w).Encode(resp)
	// 	return
	// }
	User, _ := ur.UserCtrl.GetUsers()
	paidUser, err := ur.UserCtrl.GetPaidUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := CustomResponse{Status: "failed", Message: "Not authorize to make such request"}
		json.NewEncoder(w).Encode(resp)
	}
	resp:= UserCount{
		Total: len(*User),
		TotalPaid: len(*paidUser),
	}
	json.NewEncoder(w).Encode(resp)

}

func (ur *UserRoute) FormFlag(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	params := mux.Vars(r)
	email := params["useremail"]

	err := ur.UserCtrl.UpdateForm(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := CustomResponse{Status: "failed", Message: "Error toggle form flag"}
		json.NewEncoder(w).Encode(resp)
	}
	resp := CustomResponse{Status: "success", Message: "Error toggle form flag"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}
