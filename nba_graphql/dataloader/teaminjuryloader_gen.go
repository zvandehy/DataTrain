// Code generated by github.com/vektah/dataloaden, DO NOT EDIT.

package dataloader

import (
	"sync"
	"time"

	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
)

// TeamInjuryLoaderConfig captures the config to create a new TeamInjuryLoader
type TeamInjuryLoaderConfig struct {
	// Fetch is a method that provides the data for the loader
	Fetch func(keys []int) ([][]*model.Injury, []error)

	// Wait is how long wait before sending a batch
	Wait time.Duration

	// MaxBatch will limit the maximum number of keys to send in one batch, 0 = not limit
	MaxBatch int
}

// NewTeamInjuryLoader creates a new TeamInjuryLoader given a fetch, wait, and maxBatch
func NewTeamInjuryLoader(config TeamInjuryLoaderConfig) *TeamInjuryLoader {
	return &TeamInjuryLoader{
		fetch:    config.Fetch,
		wait:     config.Wait,
		maxBatch: config.MaxBatch,
	}
}

// TeamInjuryLoader batches and caches requests
type TeamInjuryLoader struct {
	// this method provides the data for the loader
	fetch func(keys []int) ([][]*model.Injury, []error)

	// how long to done before sending a batch
	wait time.Duration

	// this will limit the maximum number of keys to send in one batch, 0 = no limit
	maxBatch int

	// INTERNAL

	// lazily created cache
	cache map[int][]*model.Injury

	// the current batch. keys will continue to be collected until timeout is hit,
	// then everything will be sent to the fetch method and out to the listeners
	batch *teamInjuryLoaderBatch

	// mutex to prevent races
	mu sync.Mutex
}

type teamInjuryLoaderBatch struct {
	keys    []int
	data    [][]*model.Injury
	error   []error
	closing bool
	done    chan struct{}
}

// Load a Injury by key, batching and caching will be applied automatically
func (l *TeamInjuryLoader) Load(key int) ([]*model.Injury, error) {
	return l.LoadThunk(key)()
}

// LoadThunk returns a function that when called will block waiting for a Injury.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
func (l *TeamInjuryLoader) LoadThunk(key int) func() ([]*model.Injury, error) {
	l.mu.Lock()
	if it, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() ([]*model.Injury, error) {
			return it, nil
		}
	}
	if l.batch == nil {
		l.batch = &teamInjuryLoaderBatch{done: make(chan struct{})}
	}
	batch := l.batch
	pos := batch.keyIndex(l, key)
	l.mu.Unlock()

	return func() ([]*model.Injury, error) {
		<-batch.done

		var data []*model.Injury
		if pos < len(batch.data) {
			data = batch.data[pos]
		}

		var err error
		// its convenient to be able to return a single error for everything
		if len(batch.error) == 1 {
			err = batch.error[0]
		} else if batch.error != nil {
			err = batch.error[pos]
		}

		if err == nil {
			l.mu.Lock()
			l.unsafeSet(key, data)
			l.mu.Unlock()
		}

		return data, err
	}
}

// LoadAll fetches many keys at once. It will be broken into appropriate sized
// sub batches depending on how the loader is configured
func (l *TeamInjuryLoader) LoadAll(keys []int) ([][]*model.Injury, []error) {
	results := make([]func() ([]*model.Injury, error), len(keys))

	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}

	injurys := make([][]*model.Injury, len(keys))
	errors := make([]error, len(keys))
	for i, thunk := range results {
		injurys[i], errors[i] = thunk()
	}
	return injurys, errors
}

// LoadAllThunk returns a function that when called will block waiting for a Injurys.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
func (l *TeamInjuryLoader) LoadAllThunk(keys []int) func() ([][]*model.Injury, []error) {
	results := make([]func() ([]*model.Injury, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([][]*model.Injury, []error) {
		injurys := make([][]*model.Injury, len(keys))
		errors := make([]error, len(keys))
		for i, thunk := range results {
			injurys[i], errors[i] = thunk()
		}
		return injurys, errors
	}
}

// Prime the cache with the provided key and value. If the key already exists, no change is made
// and false is returned.
// (To forcefully prime the cache, clear the key first with loader.clear(key).prime(key, value).)
func (l *TeamInjuryLoader) Prime(key int, value []*model.Injury) bool {
	l.mu.Lock()
	var found bool
	if _, found = l.cache[key]; !found {
		// make a copy when writing to the cache, its easy to pass a pointer in from a loop var
		// and end up with the whole cache pointing to the same value.
		cpy := make([]*model.Injury, len(value))
		copy(cpy, value)
		l.unsafeSet(key, cpy)
	}
	l.mu.Unlock()
	return !found
}

// Clear the value at key from the cache, if it exists
func (l *TeamInjuryLoader) Clear(key int) {
	l.mu.Lock()
	delete(l.cache, key)
	l.mu.Unlock()
}

func (l *TeamInjuryLoader) unsafeSet(key int, value []*model.Injury) {
	if l.cache == nil {
		l.cache = map[int][]*model.Injury{}
	}
	l.cache[key] = value
}

// keyIndex will return the location of the key in the batch, if its not found
// it will add the key to the batch
func (b *teamInjuryLoaderBatch) keyIndex(l *TeamInjuryLoader, key int) int {
	for i, existingKey := range b.keys {
		if key == existingKey {
			return i
		}
	}

	pos := len(b.keys)
	b.keys = append(b.keys, key)
	if pos == 0 {
		go b.startTimer(l)
	}

	if l.maxBatch != 0 && pos >= l.maxBatch-1 {
		if !b.closing {
			b.closing = true
			l.batch = nil
			go b.end(l)
		}
	}

	return pos
}

func (b *teamInjuryLoaderBatch) startTimer(l *TeamInjuryLoader) {
	time.Sleep(l.wait)
	l.mu.Lock()

	// we must have hit a batch limit and are already finalizing this batch
	if b.closing {
		l.mu.Unlock()
		return
	}

	l.batch = nil
	l.mu.Unlock()

	b.end(l)
}

func (b *teamInjuryLoaderBatch) end(l *TeamInjuryLoader) {
	b.data, b.error = l.fetch(b.keys)
	close(b.done)
}