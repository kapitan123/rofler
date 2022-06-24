package roflers

import (
	"context"

	"cloud.google.com/go/firestore"
)

const roflersColName = "roflers"

// actully I rarely need it. Only to store Wins and shit
func GetAll(ctx context.Context, client *firestore.Client) ([]Rofler, error) {
	docs, err := client.Collection(roflersColName).Documents(ctx).GetAll()

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
func GetByUserName(ctx context.Context, client *firestore.Client, username string) (Rofler, bool, error) {
	var r Rofler
	doc := client.Collection(roflersColName).Doc(username)
	snap, err := doc.Get(ctx)

	if err != nil {
		return r, false, nil
	}

	if err := snap.DataTo(&r); err != nil {
		return r, false, err
	}

	return r, true, nil
}

func Upsert(ctx context.Context, client *firestore.Client, r Rofler) error {
	doc := client.Collection(roflersColName).Doc(r.UserName)
	_, err := doc.Set(ctx, r)

	return err
}

func Close(client *firestore.Client) {
	client.Close()
}
