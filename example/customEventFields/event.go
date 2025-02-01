package customeventfields

import (
	"github.com/FlauschigDings/lavender"
	"github.com/google/uuid"
)

// Add the custom operator field
type Operator interface {
	Operator() *uuid.UUID
}

// Create a custom event
type Event interface {
	lavender.Event
	Operator
}
