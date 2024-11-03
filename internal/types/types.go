package types

type Student struct {
	ID    int
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"required,min=18,max=100"`
}
