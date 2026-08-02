package dto

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ChangeUserEmailRequest struct {
	Email string `json:"email"`
}
