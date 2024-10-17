package web

type CreatePerson struct {
	Name string `validate:"required,min=1,max=255" json:"name"`
	Age  int    `validate:"required,number" json:"age"`
}

type UpdatePerson struct {
	Name string `validate:"min=1,max=255" json:"name"`
	Age  int    `validate:"number" json:"age"`
}
