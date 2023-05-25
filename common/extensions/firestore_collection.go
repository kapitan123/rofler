package extensions

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

func GetAll[T any](ctx context.Context, iter *firestore.DocumentIterator) ([]T, error) {
	snapsots, err := iter.GetAll()

	if err != nil {
		return nil, errors.Wrap(err, "unable to get firestor snapshots")
	}

	var documents []T

	for _, snap := range snapsots {
		var d T
		err := snap.DataTo(&d)

		if err != nil {
			return nil, errors.Wrap(err, "unable to deserialize snapshot")
		}
		documents = append(documents, d)
	}
	return documents, nil
}
