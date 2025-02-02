package customeventfields

import (
	"github.com/FlauschigDings/lavender"
	"github.com/google/uuid"
)

type Create struct {
	User
}

// Apply implements Event.
func (c *Create) Apply(aggregate lavender.CustomAggregate[lavender.Event, lavender.Snapshot]) {
	agg := lavender.UnwrapAggregateWrapper[AccountAggregate](aggregate)
	if agg == nil {
		return
	}
	agg.Users[c.Id] = &c.User
	agg.Emails[c.Email] = &c.User
}

var _ Event = new(Create)

// Name implements lavender.Event.
func (c *Create) Name() lavender.Name {
	return "create"
}

// Operator implements Event.
func (c *Create) Operator() *uuid.UUID {
	id := uuid.New()
	return &id
}
