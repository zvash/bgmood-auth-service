package val

type RegisterUserRequest struct {
	Email                string `json:"email" validate:"required,email"`
	Name                 string `json:"name" validate:"required,min=2"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=6,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword         string `json:"current_password" validate:"required"`
	NewPassword             string `json:"new_password" validate:"required,min=6"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required,min=6,eqfield=NewPassword"`
	TerminateOtherSessions  bool   `json:"terminate_other_sessions" validate:"omitempty,boolean"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email                string `json:"email" validate:"required,email"`
	Token                string `json:"token" validate:"required"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"new_password_confirmation" validate:"required,min=6,eqfield=Password"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"omitempty,min=2"`
}

type GetUsersInfoRequest struct {
	UserIds []string `json:"user_ids" validate:"required,min=1,dive,uuid"`
}
