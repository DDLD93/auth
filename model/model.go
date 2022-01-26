package model

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type User struct{
	Id    bson.ObjectId 	`json:"id,omitempty" bson:"_id,omitempty"`
	FullName string        	`json:"fullName" bson:"fullName"`				
	Email string        	`json:"email" bson:"email"`	
	Password string        	`json:"password" bson:"password"`
	Phone string        	`json:"phone" bson:"phone"`
	IsPayment bool 			`json:"isPayment" bson:"isPayment"`
	PaymentInfo	PaymentInfo	`json:"paymentInfo" bson:"paymentInfo"`
	Role string        		`json:"role" bson:"role"` //client
	CreatedAt time.Time		`json:"createAt" bson:"CreatAt"`
}
type PaymentInfo struct{
	Channel string `json:"channel" bson:"channel"`
	Refrence string `json:"refrence" bson:"refrence"`
}