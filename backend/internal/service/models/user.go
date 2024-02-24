package models

import (

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    UserName     string      `json:"userName"      bson:"user_name" `       
    PasswordHash string      ` json:"passwordHash" bson:"password_hash"`
}
