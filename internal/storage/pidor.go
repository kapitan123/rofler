package storage

import (
	"context"
	"time"

	"google.golang.org/api/iterator"
)

type Pidor struct {
	ChosenOn time.Time `firestore:"chosen_on"`
	UserName string    `firestore:"user_name"`
	ChatId   int64     `firestore:"chat_id"`
}

const pidorsCollection = "pidors"

// AK TODO maybe split the packages?
func (s *Storage) GetAllPidors(ctx context.Context) ([]Pidor, error) {
	docs := s.client.Collection(pidorsCollection).Documents(ctx)
	var pidors []Pidor
	for {
		doc, err := docs.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, err
		}
		var p Pidor
		if err := doc.DataTo(&p); err != nil {
			return nil, err
		}
		pidors = append(pidors, p)
	}
	return pidors, nil
}

func (s *Storage) GetPidorForDate(ctx context.Context, chatid int64, date time.Time) (Pidor, bool, error) {
	var p Pidor
	// AK TODO add filter for result fetch and not an ID based on date
	doc := s.client.Collection(pidorsCollection).Doc(roundToDateOnlyString(p.ChosenOn))
	snap, err := doc.Get(ctx)

	if err != nil {
		return p, false, nil
	}

	if err := snap.DataTo(&p); err != nil {
		return p, false, err
	}

	return p, true, nil
}

func (s *Storage) CreatePidor(ctx context.Context, p Pidor) error {
	// AK TODO add filter for result fetch
	doc := s.client.Collection(pidorsCollection).Doc(roundToDateOnlyString(p.ChosenOn))
	_, err := doc.Create(ctx, p)

	return err
}

func roundToDateOnlyString(date time.Time) string {
	// AK TODO should be converted to UTC date only
	return date.Format("yyyy-MM-d")
}
