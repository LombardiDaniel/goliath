package oauth

type User struct {
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Picture      string `json:"picture"`
	Provider     string `json:"provider"`
	RefreshToken string `json:"refreshToken"`
}
