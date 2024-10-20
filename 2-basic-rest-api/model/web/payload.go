package web

type RegisterUser struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `password:"password" validate:"required,min=3,max=255"`
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
