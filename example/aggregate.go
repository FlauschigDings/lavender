package example

import (
	"github.com/FlauschigDings/lavender"
	"github.com/FlauschigDings/lavender/repo"
	"github.com/google/uuid"
)

type AccountAggregate struct {
	Users  map[uuid.UUID]*User
	Emails map[string]*User
}

func New() *AccountAggregate {
	return &AccountAggregate{
		Users:  make(map[uuid.UUID]*User),
		Emails: make(map[string]*User),
	}
}

func Load(repo *repo.Repository) (*AccountAggregate, error) {
	aggregate := New()
	if err := repo.LoadAggregate(aggregate); err != nil {
		return nil, err
	}
	return aggregate, nil
}

// Name implements lavender.CustomAggregate.
func (a *AccountAggregate) Name() lavender.Name {
	return "account"
}

// Version implements lavender.CustomAggregate.
func (a *AccountAggregate) Version() lavender.Version {
	return "0.0.1"
}

// ApplyEvent implements lavender.CustomAggregate.
func (a *AccountAggregate) ApplyEvent(event lavender.Event) {
	event.Apply(a)
}

// Events implements lavender.CustomAggregate.
func (a *AccountAggregate) Events() []lavender.Event {
	return []lavender.Event{
		new(Create),
	}
}

// ApplySnapshot implements lavender.CustomAggregate.
func (a *AccountAggregate) ApplySnapshot(snapshot lavender.Snapshot) {
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
func (a *AccountAggregate) TakeSnapshot() lavender.Snapshot {
	var users []User

	for _, fileData := range a.Users {
		users = append(users, *fileData)
	}

	return &AccountSnapshot{
		Users: users,
	}
}

var _ lavender.Aggregate = new(AccountAggregate)
