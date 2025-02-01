package store

import (
	"fmt"
	"time"

	"github.com/FlauschigDings/lavender"
	encoders "github.com/FlauschigDings/lavender/encoder"
	"github.com/thesyncim/go-clone"
	"gorm.io/gorm"
)

// Snapshot represents a stored snapshot of an aggregate.
type Snapshot struct {
	CreatedAt time.Time        // Timestamp when the snapshot was created.
	Version   lavender.Version // Aggregate version at the time of snapshot
	Name      lavender.Name    // Aggregate name
	Snapshot  string           // Serialized snapshot data
}

// Generate a table name for snapshot.
func SnapshotTableName(name lavender.Name) string {
	return fmt.Sprintf("snapshot_%s", name)
}

// Assign the dynamic table name for snapshot.
func (s Snapshot) TableName() string {
	return SnapshotTableName(s.Name)
}

// Event represents a stored event for an aggregate.
type Event struct {
	CreatedAt time.Time        // Timestamp when the event was created.
	Name      lavender.Name    // Aggregate name
	Version   lavender.Version // Aggregate version
	Topic     lavender.Name    // Event name
	Event     string           // Serialized event data
}

// Generate a table name for events.
func EventTableName(name lavender.Name) string {
	return fmt.Sprintf("event_%s", name)
}

// Assign the dynamic table name for events.
func (e Event) TableName() string {
	return EventTableName(e.Topic)
}

// Ensure GormStore implements both EventStore and SnapshotStore interfaces.
var _ SnapshotStore[lavender.Event, lavender.Snapshot] = new(GormStore[lavender.Event, lavender.Snapshot])
var _ EventStore[lavender.Event, lavender.Snapshot] = new(GormStore[lavender.Event, lavender.Snapshot])

// GormStore provides database-backed event and snapshot storage.
type GormStore[E lavender.Event, S lavender.Snapshot] struct {
	encoder          encoders.Encoder
	db               *gorm.DB
	eventRegister    map[lavender.EventIdentifier]E
	snapshotRegister map[lavender.Name]S
}

// NewGormStore initializes a GormStore with default CBOR encoding.
func NewGormStore(db *gorm.DB) *GormStore[lavender.Event, lavender.Snapshot] {
	return NewGormCustomStore[lavender.Event, lavender.Snapshot](db, encoders.NewCBorEncoder())
}

// NewGormCustomStore initializes a GormStore with a custom encoder.
func NewGormCustomStore[E lavender.Event, S lavender.Snapshot](db *gorm.DB, encoder encoders.Encoder) *GormStore[E, S] {
	return &GormStore[E, S]{
		encoder:          encoder,
		db:               db,
		eventRegister:    make(map[lavender.EventIdentifier]E),
		snapshotRegister: make(map[lavender.Name]S),
	}
}

// RegisterAggregates registers multiple aggregates for event and snapshot tracking.
func (store *GormStore[E, S]) RegisterAggregates(aggregates ...lavender.CustomAggregate[E, S]) *GormStore[E, S] {
	for _, aggregate := range aggregates {
		store.RegisterEvent(aggregate, aggregate.Events()...)
		store.RegisterSnapshot(aggregate.TakeSnapshot())
	}
	return store
}

// RegisterEvent registers event types for an aggregate and auto-migrates the event table.
func (store *GormStore[E, S]) RegisterEvent(aggregate lavender.CustomAggregate[E, S], events ...E) *GormStore[E, S] {
	for _, event := range events {
		store.eventRegister[lavender.EventId(aggregate.Name(), event.Name())] = event
		store.db.Table(EventTableName(aggregate.Name())).AutoMigrate(new(Event))
	}
	return store
}

// RegisterSnapshot registers snapshot types for an aggregate and auto-migrates the snapshot table.
func (store *GormStore[E, S]) RegisterSnapshot(snapshots ...S) *GormStore[E, S] {
	for _, snapshot := range snapshots {
		store.snapshotRegister[snapshot.AggregateID()] = snapshot
		store.db.Table(SnapshotTableName(snapshot.AggregateID())).AutoMigrate(new(Snapshot))
	}
	return store
}

// ClearEvents removes all events for an aggregate from the database.
func (store *GormStore[E, S]) ClearEvents(aggregate lavender.CustomAggregate[E, S]) error {
	return store.db.Table(EventTableName(aggregate.Name())).Where("name = ? AND version = ?", aggregate.Name(), aggregate.Version()).Delete(&Event{}).Error
}

// LoadEvents retrieves all events for an aggregate from the database.
func (store *GormStore[E, S]) LoadEvents(aggregate lavender.CustomAggregate[E, S]) (events []E, err error) {
	var readedEvents []Event

	if err := store.db.Table(EventTableName(aggregate.Name())).Where("name = ? AND version = ?", aggregate.Name(), aggregate.Version()).Find(&readedEvents).Error; err != nil {
		return nil, err
	}
	for _, eventData := range readedEvents {
		event, ok := store.eventRegister[lavender.EventId(aggregate.Name(), eventData.Topic)]
		if !ok {
			return nil, fmt.Errorf("invalid event type %s", eventData.Topic)
		}

		eventcp := clone.Clone(event).(E)
		if err := store.encoder.Unmarshal([]byte(eventData.Event), eventcp); err != nil {
			return nil, err
		}
		events = append(events, eventcp)
	}
	return events, nil
}

// SaveEvents stores multiple events for an aggregate within a database transaction.
func (store *GormStore[E, S]) SaveEvents(aggregate lavender.CustomAggregate[E, S], events []E) error {
	return store.db.Transaction(func(tx *gorm.DB) error {
		for _, event := range events {
			encodedData, err := store.encoder.Marshal(event)
			if err != nil {
				return err
			}

			err = tx.Table(EventTableName(aggregate.Name())).Create(Event{
				CreatedAt: time.Now(),
				Name:      aggregate.Name(),
				Version:   aggregate.Version(),
				Topic:     event.Name(),
				Event:     string(encodedData),
			}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// LoadSnapshot retrieves the latest snapshot for an aggregate.
func (store *GormStore[E, S]) LoadSnapshot(aggregate lavender.CustomAggregate[E, S]) (*S, error) {
	var snapshotData Snapshot

	tx := store.db.Table(SnapshotTableName(aggregate.Name())).Where("name = ? AND version = ?", aggregate.Name(), aggregate.Version()).Order("created_at DESC").First(&snapshotData)

	if tx.RowsAffected == 0 {
		return nil, nil
	}
	if err := tx.Error; err != nil {
		return nil, err
	}
	snapshot, ok := store.snapshotRegister[aggregate.Name()]
	if !ok {
		return nil, fmt.Errorf("invalid snapshot type %s", aggregate.Name())
	}

	// copy := clone.Clone(snapshot).(*S)
	if err := store.encoder.Unmarshal([]byte(snapshotData.Snapshot), snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}

// SaveSnapshot stores a snapshot of an aggregate's state.
func (store *GormStore[E, S]) SaveSnapshot(aggregate lavender.CustomAggregate[E, S], snapshot S) error {
	encodedData, err := store.encoder.Marshal(snapshot)
	if err != nil {
		return err
	}
	return store.db.Table(SnapshotTableName(aggregate.Name())).Create(&Snapshot{
		CreatedAt: time.Now(),
		Version:   aggregate.Version(),
		Name:      aggregate.Name(),
		Snapshot:  string(encodedData),
	}).Error
}
