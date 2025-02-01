package lavender

type Snapshot interface {
	// AggregateID returns the identifier of the aggregate associated with this snapshot.
	AggregateID() Name

	// Version returns the version of the aggregate at the time the snapshot was taken.
	Version() Version
}
