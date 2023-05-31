package infra

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/kapitan123/telegrofler/common/extensions"
	"github.com/kapitan123/telegrofler/service/bot/domain"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type FirestorePostsRepository struct {
	client *firestore.Client
}

func NewFirestorePostsRepository(client *firestore.Client) *FirestorePostsRepository {
	if client == nil {
		panic("missing firestoreClient")
	}

	return &FirestorePostsRepository{client}
}

func (r FirestorePostsRepository) postsCollection() *firestore.CollectionRef {
	return r.client.Collection("posts")
}

func (s FirestorePostsRepository) GetChatPosts(ctx context.Context, chatid int64) ([]domain.Post, error) {
	query := s.postsCollection().Where("chat_id", "==", chatid)
	iter := query.Documents(ctx)

	postModels, err := extensions.GetAll[PostModel](ctx, iter)

	if err != nil {
		return nil, errors.Wrap(err, "unable to get chat posts")
	}

	return mapModelsToPosts(postModels), nil
}

func mapModelsToPosts(models []PostModel) []domain.Post {
	posts := lo.Map(models, func(pm PostModel, _ int) domain.Post {
		return pm.toDomainModel()
	})
	return posts
}

func (s *FirestorePostsRepository) GetLastWeekPosts(ctx context.Context, chatid int64) ([]domain.Post, error) {
	query := s.postsCollection().Where("posted_on", ">=", time.Now().AddDate(0, 0, -7))
	query = query.Where("chat_id", "==", chatid)

	iter := query.Documents(ctx)

	postModels, err := extensions.GetAll[PostModel](ctx, iter)

	if err != nil {
		return nil, errors.Wrap(err, "unable to get last week posts")
	}

	return mapModelsToPosts(postModels), nil
}

func (s *FirestorePostsRepository) UpsertPost(ctx context.Context, p domain.Post) error {
	doc := s.postsCollection().Doc(p.Id)
	_, err := doc.Set(ctx, MapPostToModel(p))

	return err
}

func (s *FirestorePostsRepository) GetPostById(ctx context.Context, mediaId string) (domain.Post, bool, error) {
	var p PostModel
	doc := s.postsCollection().Doc(mediaId)
	snap, err := doc.Get(ctx)

	if err != nil {
		return p.toDomainModel(), false, errors.Wrap(err, "unable to get post by id")
	}

	if err := snap.DataTo(&p); err != nil {
		return p.toDomainModel(), false, errors.Wrap(err, "unable to deserialize post")
	}

	return p.toDomainModel(), true, nil
}

func (s *FirestorePostsRepository) GetByExternalSourceUrl(ctx context.Context, url string) (domain.Post, bool, error) {
	var p PostModel
	query := s.postsCollection().Where("external_source_url", "=", url).Limit(1)

	iter := query.Documents(ctx)

	postModels, err := extensions.GetAll[PostModel](ctx, iter)

	if err != nil {
		return p.toDomainModel(), false, errors.Wrap(err, "unable to get post by external_source_url")
	}

	if len(postModels) == 0 {
		return p.toDomainModel(), false, nil
	}

	return postModels[0].toDomainModel(), true, nil
}
