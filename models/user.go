package models

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	FullName     string `json:"full_name,omitempty"`
	PasswordHash string `json:"-"`
	Role         string `json:"role,omitempty"`
}

type SignupInput struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
}
