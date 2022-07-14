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
	UserId   int64     `firestore:"user_id"`
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

	// the index is built for this specific order of fields
	query := s.client.Collection(pidorsCollection).Where("chosen_on", ">=", getStartOfTheDay(date)).Where("chosen_on", "<", getEndOfTheDay(date))
	query = query.Where("chat_id", "==", chatid).Limit(1)
	iter := query.Documents(ctx)

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
	return p, true, nil

}

func (s *Storage) CreatePidor(ctx context.Context, chatid int64, username string, userid int64, date time.Time) error {
	uuid := uuid.New()

	p := &Pidor{
		ChosenOn: date,
		UserName: username,
		ChatId:   chatid,
		UserId:   userid,
	}

	doc := s.client.Collection(pidorsCollection).Doc(uuid.String())
	_, err := doc.Create(ctx, p)

	return err
}

func getStartOfTheDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func getEndOfTheDay(t time.Time) time.Time {
	oneDay := 24 * time.Hour

	rounded := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	rounded = rounded.Add(oneDay)

	return rounded
}
