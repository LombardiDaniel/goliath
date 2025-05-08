package dto

import "time"

type CreateUser struct {
	Email       string     `json:"email" binding:"email,required"`
	Password    string     `json:"password" binding:"required"`
	FirstName   string     `json:"firstName" binding:"required"`
	LastName    string     `json:"lastName" binding:"required"`
	DateOfBirth *time.Time `json:"dateOfBirth" example:"2006-01-02T15:04:05-07:00"`
}

type EditUser struct {
	FirstName   string     `json:"firstName" binding:"required"`
	LastName    string     `json:"lastName" binding:"required"`
	DateOfBirth *time.Time `json:"dateOfBirth" example:"2006-01-02T15:04:05-07:00"`
}

type UloadPicture struct {
	Content string `json:"content" binding:"required"`
}
