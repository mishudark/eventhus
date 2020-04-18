package utils

import (
	"bytes"
	"crypto/rand"
	"io"
	"time"

	"github.com/oklog/ulid"
)

// UUID returns an unique id based on ULID algorithm
func UUID() (string, error) {
	t := time.Unix(1000000, 0)
	entropy := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, entropy); err != nil {
		return "", err
	}
	uuid := ulid.MustNew(ulid.Timestamp(t), bytes.NewReader(entropy[:])).String()
	return uuid, nil
}
