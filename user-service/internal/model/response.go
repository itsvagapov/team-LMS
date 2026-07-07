package model

type AuthResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	ID    uint     `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Role  UserRole `json:"role"`
}
