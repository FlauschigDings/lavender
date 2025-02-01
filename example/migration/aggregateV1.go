package migration

import (
	"github.com/FlauschigDings/lavender"
	"github.com/FlauschigDings/lavender/repo"
	"github.com/google/uuid"
)

type AccountAggregateV1 struct {
	Users  map[uuid.UUID]*User
	Emails map[string]*User
}

func NewV1() *AccountAggregateV1 {
	return &AccountAggregateV1{
		Users:  make(map[uuid.UUID]*User),
		Emails: make(map[string]*User),
	}
}

func LoadV1(repo *repo.Repository) (*AccountAggregateV1, error) {
	aggregate := NewV1()
	if err := repo.LoadAggregate(aggregate); err != nil {
		return nil, err
	}
	return aggregate, nil
}

// Name implements lavender.CustomAggregate.
func (a *AccountAggregateV1) Name() lavender.Name {
	return "account"
}

// Version implements lavender.CustomAggregate.
func (a *AccountAggregateV1) Version() lavender.Version {
	return "1.0.0"
}

// ApplyEvent implements lavender.CustomAggregate.
func (a *AccountAggregateV1) ApplyEvent(event lavender.Event) {
	event.Apply(a)
}

// Events implements lavender.CustomAggregate.
func (a *AccountAggregateV1) Events() []lavender.Event {
	return []lavender.Event{
		new(Create),
	}
}

// ApplySnapshot implements lavender.CustomAggregate.
func (a *AccountAggregateV1) ApplySnapshot(snapshot lavender.Snapshot) {
	accountSnapshot, ok := snapshot.(*AccountSnapshot)
	if !ok {
		return
	}
	for _, user := range accountSnapshot.Users {
		a.Users[user.Id] = &user
		a.Emails[user.Email] = &user
	}
}

// TakeSnapshot implements lavender.CustomAggregate.
func (a *AccountAggregateV1) TakeSnapshot() lavender.Snapshot {
	var users []User

	for _, fileData := range a.Users {
		users = append(users, *fileData)
	}

	return &AccountSnapshot{
		Users: users,
	}
}

var _ lavender.Aggregate = new(AccountAggregateV1)
