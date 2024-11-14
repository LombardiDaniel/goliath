package models

type Token struct {
	UserId string `json:"userId" bson:"userId" binding:"required"`
	Token  string `json:"token" bson:"token" binding:"required"`
}
