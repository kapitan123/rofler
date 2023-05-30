package infra

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/kapitan123/telegrofler/service/bot/domain/media"
)

var ttl = 7 * 24 * time.Hour

type FirestoreUrlsRepository struct {
	client *firestore.Client
}

func NewFirestoreUrlsRepository(client *firestore.Client) *FirestorePostsRepository {
	if client == nil {
		panic("missing firestoreClient")
	}

	return &FirestorePostsRepository{client}
}

func (r FirestoreUrlsRepository) urlsCollection() *firestore.CollectionRef {
	return r.client.Collection("media")
}

func (s *FirestoreUrlsRepository) GetByUrl(ctx context.Context, url string) (media.Media, bool, error) {
	var u MediaModel
	doc := s.urlsCollection().Doc(url)
	snap, err := doc.Get(ctx)

	if err != nil {
		return u.toDomainModel(), false, err
	}

	if err := snap.DataTo(&u); err != nil {
		return u.toDomainModel(), false, err
	}

	return u.toDomainModel(), true, nil
}

func (s *FirestoreUrlsRepository) CreateUrl(ctx context.Context, u media.Media) error {
	doc := s.urlsCollection().Doc(u.Url)

	um := MapUrlToModel(u)
	um.ExpireAt = um.PostedOn.Add(ttl)
	_, err := doc.Create(ctx, MapUrlToModel(u))

	return err
}

type (
	MediaModel struct {
		Url      string       `firestore:"url"`
		ExpireAt time.Time    `firestore:"expire_at"`
		PostedOn time.Time    `firestore:"posted_on"`
		ChatId   int64        `firestore:"chat_id"`
		Poster   UserRefModel `firestore:"poster"`
	}
)

func (u MediaModel) toDomainModel() media.Media {
	return media.Media{
		Url:      u.Url,
		Poster:   u.Poster.toDomainModel(),
		ChatId:   u.ChatId,
		PostedOn: u.PostedOn,
	}
}

func MapUrlToModel(u media.Media) MediaModel {
	return MediaModel{
		Url:      u.Url,
		Poster:   MapUserRefToModel(u.Poster),
		ChatId:   u.ChatId,
		PostedOn: u.PostedOn,
	}
}
