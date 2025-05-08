package dto

type LoginForm struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type Email struct {
	Email string `json:"email" binding:"required"`
}

type Password struct {
	Password string `form:"password" binding:"required"`
}
