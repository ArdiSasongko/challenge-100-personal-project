package payload

type CreateUser struct {
	Name     string `json:"name" form:"name" validate:"required,min=1,max=255"`
	Age      int    `json:"age" form:"age" validate:"required,number"`
	Username string `json:"username" form:"username" validate:"required,alphanum"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,alphanum"`
}

type UpdateUser struct {
	Name     string `json:"name" form:"name" validate:"required,min=1,max=255"`
	Age      int    `json:"age" form:"age" validate:"required,number"`
	Username string `json:"username" form:"username" validate:"required,alphanum"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,alphanum"`
}