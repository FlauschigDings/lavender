package migration

import (
	"github.com/FlauschigDings/lavender"
	"github.com/FlauschigDings/lavender/repo"
)

type AccountAggregateV2 struct {
	Users []User
}

func NewV2() *AccountAggregateV2 {
	return &AccountAggregateV2{
		Users: make([]User, 0),
	}
}

func LoadV2(repo *repo.Repository) (*AccountAggregateV2, error) {
	aggregate := NewV2()
	if err := repo.LoadAggregate(aggregate); err != nil {
		return nil, err
	}
	return aggregate, nil
}

// Name implements lavender.CustomAggregate.
func (a *AccountAggregateV2) Name() lavender.Name {
	return "account"
}

// Version implements lavender.CustomAggregate.
func (a *AccountAggregateV2) Version() lavender.Version {
	return "2.0.0"
}

// ApplyEvent implements lavender.CustomAggregate.
func (a *AccountAggregateV2) ApplyEvent(event lavender.Event) {
	event.Apply(a)
}

// Events implements lavender.CustomAggregate.
func (a *AccountAggregateV2) Events() []lavender.Event {
	return []lavender.Event{
		new(Create),
	}
}

// ApplySnapshot implements lavender.CustomAggregate.
func (a *AccountAggregateV2) ApplySnapshot(snapshot lavender.Snapshot) {
	accountSnapshot, ok := snapshot.(*AccountSnapshot)
	if !ok {
		return
	}
	a.Users = accountSnapshot.Users
}

// TakeSnapshot implements lavender.CustomAggregate.
func (a *AccountAggregateV2) TakeSnapshot() lavender.Snapshot {
	return &AccountSnapshot{
		Users: a.Users,
	}
}

var _ lavender.Aggregate = new(AccountAggregateV2)
