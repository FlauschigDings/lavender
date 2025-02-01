# Lavender - Event Sourcing Library for Go
Lavender is a lightweight and flexible event sourcing library for Go (Golang). It enables you to manage event-driven systems, ensuring data consistency and persistence through events. Lavender makes it easy to implement the event sourcing pattern with support for aggregates, event storage, and event replay.

## Table of Contents
1. Installation
2. Features
3. Quick Start
4. Usage
    - Creating Aggregates
    - Persisting Events
    - Event
    - Snapshot
    - Event and Snapshot Extensions
5. Examples
6. Running Tests
8. Contributing
9. License

## 1. Installation
```bash
go get github.com/FlauschigDings/lavender
```

## 2. Features
- Event Extension: Easy extension of events for custom event types. You can easily define and handle your own events.
- Snapshot Extension: Simple integration of snapshots, allowing you to save and load aggregates at specific points in time.
- Event Storage and Retrieval: Support for persisting and retrieving events from various storage backends (e.g., SQL, NoSQL, in-memory).
- Event Replay: Replay past events to restore aggregates or reprocess events.
- Aggregate Pattern: Manage business entities with aggregates that apply events to modify their state.
- Event Handlers: Register event handlers to react to events and trigger side effects.
- Idempotency Handling: Ensure event processing is safe and repeatable.

## 3. Quick Start
Here's a simple example to get you started with Lavender:
```go
	store := store.NewInMemoryStore()

	user := *migration.NewUser("duck@ducky.com", "iL0v3Duc7s")

	repo := repo.NewRepository(store, store)

	// Add event to the repository
	if err := repo.AddEvent(migration.NewV1(), &migration.Create{
		User: user,
	}); err != nil {
		t.Fatalf("failed to add event: %v", err)
	}

	// Load  aggregate and check for errors
	agg, err := migration.LoadV1(repo)
	if err != nil {
		t.Fatalf("failed to load aggregate: %v", err)
	}
```

## 4. Usage

### 4.1 Creating Aggregates
Aggregates are core to event sourcing. They represent business entities that apply events to modify their state. Here's how you can define and use aggregates:
```go
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

var _ lavender.Aggregate = new(AccountAggregate)

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
	return "1.0.0"
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
```
### 4.2 Persisting Events
Lavender supports multiple storage backends. Here's an example of saving and loading events:
```go
    // Create event store (In-memory in this case)
	store := store.NewInMemoryStore()

	user := *migration.NewUser("duck@ducky.com", "iL0v3Duc7s")

	repo := repo.NewRepositoryConstructor(false, store, store)
	// Create each time a snapshot for snapshot testing.
	repo.AutoSnapshotHook = func(aggregate lavender.CustomAggregate[lavender.Event, lavender.Snapshot], items []lavender.Event) bool {
		return true
	}

	// Persist events to the store
	if err := repo.AddEvent(migration.NewV1(), &migration.Create{
		User: user,
	}); err != nil {
		t.Fatalf("failed to add event: %v", err)
	}
```
### 4.3 Event
Event handlers allow you to react to events. And change the current state of the Aggregate.
```go
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
```
### 4.3 Snapshot
Snapshots are compressed events. You lose the history data but it makes the eventSource calls faster.
```go
type AccountSnapshot struct {
	Users []User
}

var _ lavender.Snapshot = new(AccountSnapshot)

// AggregateID implements lavender.Snapshot.
func (s *AccountSnapshot) AggregateID() lavender.Name {
	return new(AccountAggregate).Name()
}

// Version implements lavender.Snapshot.
func (s *AccountSnapshot) Version() lavender.Version {
	return new(AccountAggregate).Version()
}
```
### 4.4 Event and Snapshot Extensions
```go
// Add the custom operator field
type Extension interface {
    Operator() *uuid.UUID
}

// Create a custom event
type Event interface {
	lavender.Event
	Extension
}


// Create a custom event
type Snapshot interface {
	lavender.Snapshot
	Extension
}
```
## 5. Example
All examples can be found in the /example directory. Some of the key examples include:
- [Event and Snapshot Extensions](https://github.com/FlauschigDings/lavender/tree/master/example/customEventFields)
- [Migration](https://github.com/FlauschigDings/lavender/tree/master/example/migration)

## 6. Running Tests
Lavender supports tests using Go’s built-in testing framework. To run tests for the entire project, simply use:
```bash
go test ./...
```

## 7. Contributing
We welcome contributions to Lavender! Here's how you can contribute:

1. Fork the repository.
2. Clone your fork and create a new branch.
3. Make your changes and add tests.
4. Open a pull request with a clear description of your changes.

## 8. License
Lavender is licensed under the MIT License. See the [LICENSE](https://github.com/FlauschigDings/lavender/blob/master/LICENSE) file for more details.

## 9. Contact
For any questions, feel free to reach out at discord: flauschig or create a issue.