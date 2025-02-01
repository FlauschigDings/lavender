package example

import "github.com/FlauschigDings/lavender"

type Create struct {
	User
}

var _ lavender.Event = new(Create)

// Name implements lavender.Event.
func (c *Create) Name() lavender.Name {
	return "create"
}

// Apply implements lavender.Event.
func (c *Create) Apply(aggregate lavender.Aggregate) {
	accountAggregate := lavender.Parse[*AccountAggregate](aggregate)
	accountAggregate.Users[c.Id] = &c.User
	accountAggregate.Emails[c.Email] = &c.User
}
