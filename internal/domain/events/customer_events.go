package events

import "github.com/google/uuid"

type CustomerLinkedToUserEvent struct {
	CustomerID uuid.UUID
	UserID     uuid.UUID
}

func (CustomerLinkedToUserEvent) EventName() string { return "customer.linked_to_user" }

type CustomerRenamedEvent struct {
	CustomerID uuid.UUID
	OldName    string
	NewName    string
}

func (CustomerRenamedEvent) EventName() string { return "customer.renamed" }
