package migration

import "github.com/FlauschigDings/lavender"

type AccountSnapshot struct {
	Users []User `json:"users"`
}

var _ lavender.Snapshot = new(AccountSnapshot)

// AggregateID implements lavender.Snapshot.
func (s *AccountSnapshot) AggregateID() lavender.Name {
	return new(AccountAggregateV1).Name()
}

// Version implements lavender.Snapshot.
func (s *AccountSnapshot) Version() lavender.Version {
	return new(AccountAggregateV1).Version()
}
