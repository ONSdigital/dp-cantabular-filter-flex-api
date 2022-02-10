package model

import (
	"crypto/sha1"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// Hash generates a SHA-1 hash of the filter struct. SHA-1 is not cryptographically safe,
// but it has been selected for performance as we are only interested in uniqueness.
// ETag field value is ignored when generating a hash.
// An optional byte array can be provided to append to the hash.
// This can be used, for example, to calculate a hash of this instance and an update applied to it.
func (f *Filter) Hash(extraBytes []byte) (string, error) {
	h := sha1.New()

	// copy by value to ignore ETag without affecting i
	f2 := *f
	f2.ETag = ""

	b, err := bson.Marshal(f2)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal as bson")
	}

	if _, err := h.Write(append(b, extraBytes...)); err != nil {
		return "", errors.Wrap(err, "failed to write to sha1 hash body")
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
