package entities

import "testing"

func TestRole_IsValid(t *testing.T) {
	tests := []struct {
		name string
		role Role
		want bool
	}{
		{"admin é válido", RoleAdmin, true},
		{"customer é válido", RoleCustomer, true},
		{"manager é válido", RoleManager, true},
		{"string vazia é inválida", Role(""), false},
		{"role desconhecida é inválida", Role("superadmin"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.role.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRole_String(t *testing.T) {
	tests := []struct {
		name string
		role Role
		want string
	}{
		{"admin", RoleAdmin, "admin"},
		{"customer", RoleCustomer, "customer"},
		{"manager", RoleManager, "manager"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.role.String()
			if got != tt.want {
				t.Errorf("String() got = %q, want %q", got, tt.want)
			}
		})
	}
}
