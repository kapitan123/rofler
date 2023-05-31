package extensions

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

type firebaseDocIdSetter[TT any] interface {
	SetId(id string)
	*TT
}

func GetAll[T any, PT firebaseDocIdSetter[T]](ctx context.Context, iter *firestore.DocumentIterator) ([]T, error) {
	snapshots, err := iter.GetAll()

	if err != nil {
		return nil, errors.Wrap(err, "unable to get firestore snapshots")
	}

	var documents []T

	for _, snap := range snapshots {
		var d T
		err := snap.DataTo(&d)

		if err != nil {
			return nil, errors.Wrap(err, "unable to deserialize snapshot")
		}
		pointerToD := PT(&d)
		pointerToD.SetId(snap.Ref.ID)
		documents = append(documents, d)
	}
	return documents, nil
}
