package dto

type CreateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	CPF   string `json:"cpf"`
}

type UpdateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
