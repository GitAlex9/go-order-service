package events

import (
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

type UserPasswordChangedEvent struct {
	UserID uuid.UUID
}

type UserEmailChangedEvent struct {
	UserID   uuid.UUID
	OldEmail valueobjects.Email
	NewEmail valueobjects.Email
}

type UserDeactivatedEvent struct {
	UserID uuid.UUID
}

type UserActivatedEvent struct {
	UserID uuid.UUID
}
