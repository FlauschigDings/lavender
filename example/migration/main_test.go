package migration_test

import (
	"testing"

	"github.com/FlauschigDings/lavender"
	"github.com/FlauschigDings/lavender/example/migration"
	"github.com/FlauschigDings/lavender/repo"
	"github.com/FlauschigDings/lavender/store"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestMigrate(t *testing.T) {
	store := store.NewInMemoryStore()

	user := *migration.NewUser("duck@ducky.com", "iL0v3Duc7s")

	repo := repo.NewRepositoryConstructor(false, store, store)
	// Create each time a snapshot for snapshot testing.
	repo.AutoSnapshotHook = func(aggregate lavender.CustomAggregate[lavender.Event, lavender.Snapshot], items []lavender.Event) bool {
		return true
	}

	// Add event to the repository
	if err := repo.AddEvent(migration.NewV1(), &migration.Create{
		User: user,
	}); err != nil {
		t.Fatalf("failed to add event: %v", err)
	}

	// Load V1 aggregate and check for errors
	aggV1, err := migration.LoadV1(repo)
	if err != nil {
		t.Fatalf("failed to load V1 aggregate: %v", err)
	}

	// Load V2 aggregate and check for errors
	aggV2, err := migration.LoadV2(repo)
	if err != nil {
		t.Fatalf("failed to load V2 aggregate: %v", err)
	}

	// Assertions
	assert.Contains(t, maps.Keys(aggV1.Users), user.Id, "V1 should contain the user")
	assert.Contains(t, aggV2.Users, user, "V2 should contain the user")
}
