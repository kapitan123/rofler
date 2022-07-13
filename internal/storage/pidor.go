package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
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
	dateOnly := getOnlyDate(p.ChosenOn)

	// the index is built for this specific order of fields
	query := s.client.Collection(pidorsCollection).Where("chosen_on", "==", dateOnly).Where("chat_id", "==", chatid).Limit(1)
	iter := query.Documents(ctx)

	for {
		snap, err := iter.Next()
		if err == iterator.Done {
			return p, false, nil
		}
		if err != nil {
			return p, false, err
		}
		if err := snap.DataTo(&p); err != nil {
			return p, false, err
		}
	}
}

func (s *Storage) CreatePidor(ctx context.Context, chatid int64, username string, date time.Time) error {
	uuid := uuid.New()

	p := &Pidor{
		ChosenOn: getOnlyDate(date),
		UserName: username,
		ChatId:   chatid,
	}

	doc := s.client.Collection(pidorsCollection).Doc(uuid.String())
	_, err := doc.Create(ctx, p)

	return err
}

func getOnlyDate(t time.Time) time.Time {
	d := (24 * time.Hour)
	return t.Truncate(d)
}
