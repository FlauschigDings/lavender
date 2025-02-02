package lavender

type AggregateWrapper[T any] struct {
	HookApplyEvent    func(event Event)
	HookApplySnapshot func(snapshot Snapshot)
	HookEvents        func() []Event
	HookName          func() Name
	HookTakeSnapshot  func() Snapshot
	HookVersion       func() Version
	Data              *T
}

// ApplyEvent implements CustomAggregate.
func (a *AggregateWrapper[any]) ApplyEvent(event Event) {
	a.HookApplyEvent(event)
}

// ApplySnapshot implements CustomAggregate.
func (a *AggregateWrapper[any]) ApplySnapshot(snapshot Snapshot) {
	a.HookApplySnapshot(snapshot)
}

// Events implements CustomAggregate.
func (a *AggregateWrapper[any]) Events() []Event {
	return a.HookEvents()
}

// Name implements CustomAggregate.
func (a *AggregateWrapper[any]) Name() Name {
	return a.HookName()
}

// TakeSnapshot implements CustomAggregate.
func (a *AggregateWrapper[any]) TakeSnapshot() Snapshot {
	return a.HookTakeSnapshot()
}

// Version implements CustomAggregate.
func (a *AggregateWrapper[any]) Version() Version {
	return a.HookVersion()
}

var _ Aggregate = new(AggregateWrapper[any])

func UnwrapAggregateWrapper[T any](aggreage CustomAggregate[Event, Snapshot]) *T {
	wrapper, ok := aggreage.(*AggregateWrapper[T])
	if !ok {
		return nil
	}
	return wrapper.Data
}
