package repo

import "github.com/FlauschigDings/lavender"

// AutoSnapshotHook is a hook that can be used to automatically snapshot an aggregate when it has been modified.
type AutoSnapshotHook[E lavender.Event, S lavender.Snapshot] func(aggregate lavender.CustomAggregate[E, S], items []E) bool
