package events

import (
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

type UserPasswordChangedEvent struct {
	UserID uuid.UUID
}

func (UserPasswordChangedEvent) EventName() string { return "user.password_changed" }

type UserEmailChangedEvent struct {
	UserID   uuid.UUID
	OldEmail valueobjects.Email
	NewEmail valueobjects.Email
}

func (UserEmailChangedEvent) EventName() string { return "user.email_changed" }

type UserDeactivatedEvent struct {
	UserID uuid.UUID
}

func (UserDeactivatedEvent) EventName() string { return "user.deactivated" }

type UserActivatedEvent struct {
	UserID uuid.UUID
}

func (UserActivatedEvent) EventName() string { return "user.activated" }
