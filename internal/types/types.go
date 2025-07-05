package types


type Student struct {
	ID        int    `validate:"required"`
	Name      string `validate:"required"`
	Age       int    `validate:"required"`
	Email     string `validate:"required"`
	City	  string `validate:"required"`
}

