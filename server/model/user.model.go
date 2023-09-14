package model

type User struct {
	Id                 int    `json:"id"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	IsVerified         bool   `json:"is_verified"`
	ResetPasswordToken string `json:"reset_password_token"`
	VerificationToken  string `json:"verification_token"`
}
