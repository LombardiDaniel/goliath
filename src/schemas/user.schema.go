package schemas

import "time"

type CreateUser struct {
	Email       string     `json:"email" binding:"email,required"`
	Password    string     `json:"password" binding:"required"`
	FirstName   string     `json:"firstName" binding:"required"`
	LastName    string     `json:"lastName" binding:"required"`
	DateOfBirth *time.Time `json:"dateOfBirth"`
}
