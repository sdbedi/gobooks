package store

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redeam/gobooks/objects"
)

// IBookStore is the database interface for storing Books
type IBookStore interface {
	Get(ctx context.Context, in *objects.GetRequest) (*objects.Book, error)
	List(ctx context.Context, in *objects.ListRequest) ([]*objects.Book, error)
	Create(ctx context.Context, in *objects.CreateRequest) error
	UpdateDetails(ctx context.Context, in *objects.UpdateDetailsRequest) error
	Delete(ctx context.Context, in *objects.DeleteRequest) error
}

func init() {
	rand.Seed(time.Now().UTC().Unix())
}

// GenerateUniqueID will returns a time based sortable unique id
func GenerateUniqueID() string {
	word := []byte("0987654321")
	rand.Shuffle(len(word), func(i, j int) {
		word[i], word[j] = word[j], word[i]
	})
	now := time.Now().UTC()
	return fmt.Sprintf("%010v-%010v-%s", now.Unix(), now.Nanosecond(), string(word))
}
