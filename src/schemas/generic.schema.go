package schemas

type IdString struct {
	Id string `json:"id" binding:"required"`
}
