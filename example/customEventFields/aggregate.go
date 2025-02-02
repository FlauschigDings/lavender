package customeventfields

import (
	"github.com/FlauschigDings/lavender"
	"github.com/FlauschigDings/lavender/repo"
	"github.com/google/uuid"
)

type Repository = repo.CustomRepository[Event, lavender.Snapshot]

var _ lavender.CustomAggregate[Event, lavender.Snapshot] = new(AccountAggregate)

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

func Load(repo *Repository) (*AccountAggregate, error) {
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
func (a *AccountAggregate) ApplyEvent(event Event) {
	event.Apply(a.Convert())
}

// Events implements lavender.CustomAggregate.
func (a *AccountAggregate) Events() []Event {
	return []Event{
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

// Version implements lavender.CustomAggregate.
func (a *AccountAggregate) Convert() lavender.CustomAggregate[lavender.Event, lavender.Snapshot] {
	return &lavender.AggregateWrapper[AccountAggregate]{
		HookApplyEvent: func(event lavender.Event) {
			switch event := event.(type) {
			case Event:
				a.ApplyEvent(event)
			default:
				return
			}
		},
		HookApplySnapshot: a.ApplySnapshot,
		HookEvents: func() (events []lavender.Event) {
			for _, event := range a.Events() {
				events = append(events, event)
			}
			return
		},
		HookName:         a.Name,
		HookTakeSnapshot: a.TakeSnapshot,
		HookVersion:      a.Version,
		Data:             a,
	}
}
