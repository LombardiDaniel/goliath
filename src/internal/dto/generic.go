package dto

type Id struct {
	Id string `json:"id" binding:"required"`
}

type Url struct {
	Url string `json:"url" binding:"required"`
}
