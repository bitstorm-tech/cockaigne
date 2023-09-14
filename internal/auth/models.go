package auth

type CreateAccountRequest struct {
	Username       string
	Password       string
	PasswordRepeat string
	Email          string
	IsDealer       bool
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
