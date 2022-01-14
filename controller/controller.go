package controller

import (
	"errors"
	"fmt"
	"log"

	"github.com/ddld93/auth/model"
	"github.com/ddld93/auth/utilities"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DB_Connect struct {
	Session *mgo.Session
}
var (
	database = "testdb"
	collection= "student"
)

func NewUserCtrl(host string, port int) *DB_Connect {
	url := fmt.Sprintf("mongodb://%s:%d", host, port)
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal("Error connecting to mogo database", err)
	}
	return &DB_Connect{Session: session}
}

func (u *DB_Connect) CreateUser(user *model.User) (string, error) {
	// validating user inputs
	user, err  := utilities.UserModelValidate(user)
	if err != nil{
		return "",err
	}
	//checking if user with same email exist
	resp,_:= u.Session.DB(database).C(collection).Find(bson.M{"email":user.Email}).Count()
	if resp >= 1 {
		return "", errors.New("an Account with this email already exist")
	}
	err2 := u.Session.DB(database).C(collection).Insert(user)
	if err2 != nil {
		fmt.Println("Error inserting new user ", err2)
		return "", err2
	}
	fmt.Println("User inserted successfully!")
	return "User added successifully", nil
}

func (u *DB_Connect) GetUser(email string) (*model.User, error) {
	user := model.User{}
	err := u.Session.DB(database).C(collection).Find(bson.M{"email":email}).One(&user)
	if err != nil {
		return &user, errors.New("error getting user by email ")
	}
	fmt.Println("Sigle User found !")
	return &user, nil
}