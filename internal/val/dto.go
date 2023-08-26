package val

type RegisterUserRequest struct {
	Email                string `json:"email" validate:"required,email"`
	Name                 string `json:"name" validate:"required"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=6,eqfield=Password"`
}
