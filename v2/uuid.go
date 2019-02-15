package eventhus

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

// GenerateUUID returns an ULID id
func GenerateUUID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
