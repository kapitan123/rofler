package roflers

import (
	"context"

	"cloud.google.com/go/firestore"
)

const roflersColName = "roflers"

type RoflerStorage struct {
	client *firestore.Client
}

func New(client *firestore.Client) *RoflerStorage {
	return &RoflerStorage{
		client: client,
	}
}

// actully I rarely need it. Only to store Wins and shit
func (s *RoflerStorage) GetAll(ctx context.Context) ([]Rofler, error) {
	docs, err := s.client.Collection(roflersColName).Documents(ctx).GetAll()

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
func (s *RoflerStorage) GetByUserName(ctx context.Context, username string) (Rofler, bool, error) {
	var r Rofler
	doc := s.client.Collection(roflersColName).Doc(username)
	snap, err := doc.Get(ctx)

	if err != nil {
		return r, false, nil
	}

	if err := snap.DataTo(&r); err != nil {
		return r, false, err
	}

	return r, true, nil
}

func (s *RoflerStorage) Upsert(ctx context.Context, r Rofler) error {
	doc := s.client.Collection(roflersColName).Doc(r.UserName)
	_, err := doc.Set(ctx, r)

	return err
}

func (s *RoflerStorage) Close() {
	s.client.Close()
}
