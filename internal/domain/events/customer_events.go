package events

import (
	"github.com/google/uuid"
)

type CustomerLinkedToUserEvent struct {
	CustomerID uuid.UUID
	UserID     uuid.UUID
}

type CustomerRenamedEvent struct {
	CustomerID uuid.UUID
	OldName    string
	NewName    string
}
