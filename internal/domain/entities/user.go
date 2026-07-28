package entities

import (
	"errors"
	"time"

	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

type User struct {
	id        uuid.UUID
	email     valueobjects.Email
	password  valueobjects.Password
	role      Role
	active    bool
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(email valueobjects.Email, plainPassword string, role Role) (*User, error) {
	if !role.IsValid() {
		return nil, errors.New("role inválida")
	}
	password, err := valueobjects.NewPassword(plainPassword)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &User{
		id:        uuid.New(),
		email:     email,
		password:  password,
		role:      role,
		active:    true,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func RestoreUser(
	id uuid.UUID,
	email valueobjects.Email,
	passwordHash string,
	role Role,
	active bool,
	createdAt,
	updatedAt time.Time,
) *User {
	return &User{
		id:        id,
		email:     email,
		password:  valueobjects.NewPasswordFromHash(passwordHash),
		role:      role,
		active:    active,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (u *User) ID() uuid.UUID             { return u.id }
func (u *User) Email() valueobjects.Email { return u.email }
func (u *User) Role() Role                { return u.role }
func (u *User) Active() bool              { return u.active }
func (u *User) PasswordHash() string      { return u.password.Hash() }
func (u *User) CreatedAt() time.Time      { return u.createdAt }
func (u *User) UpdatedAt() time.Time      { return u.updatedAt }

func (u *User) CheckPassword(plain string) bool {
	return u.password.Matches(plain)
}

func (u *User) ChangePassword(currentPlain, newPlain string) error {
	if !u.password.Matches(currentPlain) {
		return errors.New("senha atual incorreta")
	}
	newPassword, err := valueobjects.NewPassword(newPlain)
	if err != nil {
		return err
	}
	u.password = newPassword
	u.updatedAt = time.Now()
	return nil
}

func (u *User) ChangeEmail(newEmail valueobjects.Email) {
	u.email = newEmail
	u.updatedAt = time.Now()
}

func (u *User) Deactivate() { u.active = false; u.updatedAt = time.Now() }
func (u *User) Activate()   { u.active = true; u.updatedAt = time.Now() }
