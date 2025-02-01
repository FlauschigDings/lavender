package repo

import "github.com/FlauschigDings/lavender"

// LoadCache loads an aggregate from the cache.
func (r *CustomRepository[E, S]) LoadCache(aggregate lavender.CustomAggregate[E, S]) *lavender.CustomAggregate[E, S] {
	if !r.aggregateCacheActive {
		return nil
	}
	if cache, ok := r.aggregateCache.Load(aggregate.Name()); ok {
		item := cache.(lavender.CustomAggregate[E, S])
		return &item
	}
	return nil
}

// SaveCache saves an aggregate to the cache.
func (r *CustomRepository[E, S]) saveCache(aggregate lavender.CustomAggregate[E, S]) {
	if !r.aggregateCacheActive {
		return
	}
	r.aggregateCache.Store(aggregate.Name(), aggregate)
}
