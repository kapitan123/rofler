package roflers

// AK TODO rename package
import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/kapitan123/telegrofler/config"
)

const collectionName = "roflers"

type RoflersStore struct {
	ctx    *context.Context
	client *firestore.Client
	col    *firestore.CollectionRef
}

func New() *RoflersStore {
	ctx := context.Background()
	client := createClient(ctx)
	return &RoflersStore{
		&ctx,
		client,
		client.Collection(collectionName),
	}
}

func (rs *RoflersStore) GetAll() ([]Rofler, error) {
	docs, err := rs.col.Documents(*rs.ctx).GetAll()

	if err != nil {
		return nil, err
	}

	roflers := []Rofler{}
	for _, doc := range docs {
		r := Rofler{}
		doc.DataTo(&Rofler{})
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
	doc := rs.col.Doc(username)
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
	doc := rs.col.Doc(r.UserName)
	_, err := doc.Set(*rs.ctx, r)

	if err != nil {
		return err
	}

	return nil
}

func (rs *RoflersStore) Close() {
	rs.client.Close()
}

func createClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, config.ProjectId)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}
