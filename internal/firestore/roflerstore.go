package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

const roflersColName = "roflers"

type RoflersStore struct {
	ctx        *context.Context
	client     *firestore.Client
	roflersCol *firestore.CollectionRef
}

func NewRoflerStore() *RoflersStore {
	ctx := context.Background()
	client := createClient(ctx)
	return &RoflersStore{
		&ctx,
		client,
		client.Collection(roflersColName),
	}
}

// actully I rarely need it. Only to store Wins and shit
func (rs *RoflersStore) GetAllRoflers() ([]Rofler, error) {
	docs, err := rs.roflersCol.Documents(*rs.ctx).GetAll()

	if err != nil {
		return nil, err
	}

	roflers := []Rofler{}
	for _, doc := range docs {
		r := Rofler{}
		doc.DataTo(&r)
		roflers = append(roflers, r)
	}

	return roflers, nil
}

// Tries to fetch a document by id
// Internally firestore throws an error if the document does not exist.
// We treat all errors on fetching as not found.
// Errors on convertion are treated in a regural way
func (rs *RoflersStore) GetByUserName(username string) (Rofler, bool, error) {
	var r Rofler
	doc := rs.roflersCol.Doc(username)
	snap, err := doc.Get(*rs.ctx)

	if err != nil {
		return r, false, nil
	}

	if err := snap.DataTo(&r); err != nil {
		return r, false, err
	}

	return r, true, nil
}

func (rs *RoflersStore) Upsert(r Rofler) error {
	doc := rs.roflersCol.Doc(r.UserName)
	_, err := doc.Set(*rs.ctx, r)

	return err
}

func (rs *RoflersStore) Close() {
	rs.client.Close()
}
