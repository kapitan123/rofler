package infra

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	url "github.com/kapitan123/telegrofler/service/bot/domain/media-url"
)

type FirestoreUrlsRepository struct {
	client *firestore.Client
}

func NewFirestoreUrlsRepository(client *firestore.Client) *FirestorePostsRepository {
	if client == nil {
		panic("missing firestoreClient")
	}

	return &FirestorePostsRepository{client}
}

func (r FirestorePostsRepository) urlsCollection() *firestore.CollectionRef {
	return r.client.Collection("urls")
}

func (s *FirestorePostsRepository) GetUrlByAddr(ctx context.Context, url string) (url.MediaUrl, bool, error) {
	var u UrlModel
	doc := s.postsCollection().Doc(url)
	snap, err := doc.Get(ctx)

	if err != nil {
		return u.toDomainModel(), false, err
	}

	if err := snap.DataTo(&u); err != nil {
		return u.toDomainModel(), false, err
	}

	return u.toDomainModel(), true, nil
}

func (s *FirestorePostsRepository) CreateUrl(ctx context.Context, u url.MediaUrl) error {
	doc := s.postsCollection().Doc(u.Url)

	_, err := doc.Create(ctx, MapUrlToModel(u))

	return err
}

type (
	UrlModel struct {
		Url      string       `firestore:"url"`
		PostedOn time.Time    `firestore:"posted_on"`
		ChatId   int64        `firestore:"chat_id"`
		Poster   UserRefModel `firestore:"poster"`
	}
)

func (u UrlModel) toDomainModel() url.MediaUrl {
	return url.MediaUrl{
		Url:      u.Url,
		Poster:   u.Poster.toDomainModel(),
		ChatId:   u.ChatId,
		PostedOn: u.PostedOn,
	}
}

func MapUrlToModel(u url.MediaUrl) UrlModel {
	return UrlModel{
		Url:      u.Url,
		Poster:   MapUserRefToModel(u.Poster),
		ChatId:   u.ChatId,
		PostedOn: u.PostedOn,
	}
}
