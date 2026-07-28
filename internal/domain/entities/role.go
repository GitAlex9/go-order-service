package entities

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
	RoleManager  Role = "manager"
)

// validRoles funciona como o "conjunto" de valores aceitos.
// Usar map em vez de switch facilita adicionar roles novas sem repetir lógica.
var validRoles = map[Role]bool{
	RoleAdmin:    true,
	RoleCustomer: true,
	RoleManager:  true,
}

func (r Role) IsValid() bool {
	return validRoles[r]
}

func (r Role) String() string {
	return string(r)
}
