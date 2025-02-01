package store_test

import (
	"sync"
	"testing"

	"github.com/FlauschigDings/lavender"
	"github.com/FlauschigDings/lavender/example"
	"github.com/FlauschigDings/lavender/store"
)

func TestSaveEvent(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		memStore := store.NewInMemoryStore()
		err := memStore.SaveEvents(example.New(), []lavender.Event{
			&example.Create{
				User: *example.NewUser("Nils0", "dasIstMeinPassword,Ja das ist toll"),
			},
		})
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("multi", func(t *testing.T) {
		memStore := store.NewInMemoryStore()

		err := memStore.SaveEvents(example.New(), []lavender.Event{
			&example.Create{
				User: *example.NewUser("Nils1", "dasIstMeinPassword,Ja das ist toll"),
			},
			&example.Create{
				User: *example.NewUser("Nils1", "dasIstMeinPassword,Ja das ist toll"),
			},
			&example.Create{
				User: *example.NewUser("Nils1", "dasIstMeinPassword,Ja das ist toll"),
			},
			&example.Create{
				User: *example.NewUser("Nils1", "dasIstMeinPassword,Ja das ist toll"),
			},
			&example.Create{
				User: *example.NewUser("Nils1", "dasIstMeinPassword,Ja das ist toll"),
			},
			&example.Create{
				User: *example.NewUser("Nils1", "dasIstMeinPassword,Ja das ist toll"),
			},
			&example.Create{
				User: *example.NewUser("Nils1", "dasIstMeinPassword,Ja das ist toll"),
			},

			&example.Create{
				User: *example.NewUser("Nils1", "dasIstMeinPassword,Ja das ist toll"),
			},
			&example.Create{
				User: *example.NewUser("Nils1", "dasIstMeinPassword,Ja das ist toll"),
			},
		})
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("multi-thread", func(t *testing.T) {
		memoryStore := store.NewInMemoryStore()
		var wg sync.WaitGroup

		for i := 0; i < 10000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				memoryStore.SaveEvents(example.New(), []lavender.Event{
					&example.Create{
						User: *example.NewUser("Nils2", "dasIstMeinPassword,Ja das ist toll"),
					},
				})
				events, err := memoryStore.LoadEvents(example.New())
				if err != nil {
					t.Error(err)
				}
				_ = events
			}()
		}

		wg.Wait()
	})
}

func TestLoadEvent(t *testing.T) {
	t.Run("single", func(t *testing.T) {
	})

	t.Run("multi", func(t *testing.T) {
	})

	t.Run("multi-thread", func(t *testing.T) {

	})
}

func TestSaveSnapshot(t *testing.T) {
	t.Run("single", func(t *testing.T) {
	})

	t.Run("multi", func(t *testing.T) {
	})

	t.Run("multi-thread", func(t *testing.T) {

	})
}

func TestLoadSnapshot(t *testing.T) {
	t.Run("single", func(t *testing.T) {
	})

	t.Run("multi", func(t *testing.T) {
	})

	t.Run("multi-thread", func(t *testing.T) {

	})
}
