package customeventfields_test

import (
	"testing"

	"github.com/FlauschigDings/lavender"
	customeventfields "github.com/FlauschigDings/lavender/example/customEventFields"
	"github.com/FlauschigDings/lavender/repo"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestCustomEventFields(t *testing.T) {
	store := customeventfields.NewCustomMemoryStore()

	user := *customeventfields.NewUser("duck@ducky.com", "iL0v3Duc7s")

	repo := repo.NewRepositoryConstructor(false, store, store)
	// Create each time a snapshot for snapshot testing.
	repo.AutoSnapshotHook = func(aggregate lavender.CustomAggregate[customeventfields.Event, lavender.Snapshot], items []customeventfields.Event) bool {
		return true
	}

	// Add event to the repository
	if err := repo.AddEvent(customeventfields.New(), &customeventfields.Create{
		User: user,
	}); err != nil {
		t.Fatalf("failed to add event: %v", err)
	}

	// Load V1 aggregate and check for errors
	aggV1, err := customeventfields.Load(repo)
	if err != nil {
		t.Fatalf("failed to load V1 aggregate: %v", err)
	}

	// Assertions
	assert.Contains(t, maps.Keys(aggV1.Users), user.Id, "V1 should contain the user")
}
