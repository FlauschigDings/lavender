package migration

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
	if aggreage, ok := aggregate.(*AccountAggregateV1); ok {
		aggreage.Users[c.Id] = &c.User
		aggreage.Emails[c.Email] = &c.User
	}
	if aggreage, ok := aggregate.(*AccountAggregateV2); ok {
		aggreage.Users = append(aggreage.Users, c.User)
	}
}
