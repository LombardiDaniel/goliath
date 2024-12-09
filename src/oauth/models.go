package oauth

type User struct {
	Email        string  `json:"email"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	PictureUrl   *string `json:"pictureUrl"`
	Provider     string  `json:"provider"`
	RefreshToken string  `json:"refreshToken"`
}
