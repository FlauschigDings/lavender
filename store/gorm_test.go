package store_test

import (
	"fmt"
	"testing"

	"github.com/FlauschigDings/lavender"
	encoders "github.com/FlauschigDings/lavender/encoder"
	"github.com/FlauschigDings/lavender/example"
	"github.com/FlauschigDings/lavender/repo"
	"github.com/FlauschigDings/lavender/store"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseService struct {
	*gorm.DB
}

func New(db *gorm.DB) *DatabaseService {
	return &DatabaseService{
		db,
	}
}

func Sqlite(fileName string) (*DatabaseService, error) {
	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	return New(db), nil
}

var encoderList = []encoders.Encoder{
	encoders.NewGobEncoder(),
	encoders.NewCBorEncoder(),
	encoders.NewJsonEncoder(),
}
var accounts = []string{
	"a@t.de",
	"b@t.de",
	"c@t.de",
	"d@t.de",
}

func TestGorm(t *testing.T) {
	db, err := Sqlite(":memory:")
	if err != nil {
		t.Error(err)
	}
	store := store.NewGormStore(db.DB)
	store.RegisterAggregates(example.New())
	repo := repo.NewRepository(store, store)

	for _, user := range accounts {
		repo.AddEvent(example.New(), &example.Create{
			User: example.User{
				Id:       uuid.New(),
				Email:    user,
				Password: user,
			},
		})

		data, err := example.Load(repo)
		if err != nil {
			t.Error(err)
		}

		for user := range data.Emails {
			fmt.Println(user)
		}
	}

	events, err := store.LoadEvents(example.New())
	if err != nil {
		t.Error(err)
	}

	for _, e := range events {
		fmt.Println(e)
	}

}

func TestEncoder(t *testing.T) {
	for _, encoder := range encoderList {
		db, err := Sqlite(":memory:")
		if err != nil {
			t.Error(err)
		}

		store := store.NewGormCustomStore[lavender.Event, lavender.Snapshot](db.DB, encoder)
		store.RegisterAggregates(example.New())
		repo := repo.NewRepositoryConstructor(false, store, store)

		for _, user := range accounts {
			repo.AddEvent(example.New(), &example.Create{
				User: example.User{
					Id:       uuid.New(),
					Email:    user,
					Password: user,
				},
			})

		}

		data, err := example.Load(repo)
		if err != nil {
			t.Error(err)
		}

		for user := range data.Emails {
			assert.Contains(t, accounts, user, "user not exist")
		}

	}
}
