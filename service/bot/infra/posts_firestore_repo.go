package infra

import (
	"context"
	"net/url"
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

type (
	PostModel struct {
		Id                string          `firestore:"-"`
		ExternalSourceUrl string          `firestore:"external_source_url"`
		Reactions         []ReactionModel `firestore:"reactions"`
		PostedOn          time.Time       `firestore:"posted_on"`
		ChatId            int64           `firestore:"chat_id"`
		Poster            UserRefModel    `firestore:"poster"`
	}

	ReactionModel struct {
		ReactToMessageId int          `firestore:"react_to_message_id"`
		Text             string       `firestore:"text"`
		PostedOn         time.Time    `firestore:"posted_on"`
		Reactor          UserRefModel `firestore:"reactor"`
	}

	UserRefModel struct {
		DisplayName string `firestore:"user_name"`
		Id          int64  `firestore:"user_id"`
	}
)

func (pm *PostModel) SetId(id string) {
	pm.Id = id
}

func (urm UserRefModel) toDomainModel() domain.UserRef {
	return domain.UserRef{
		DisplayName: urm.DisplayName,
		Id:          urm.Id,
	}
}

func (rm ReactionModel) toDomainModel() domain.Reaction {
	return domain.Reaction{
		Text:     rm.Text,
		PostedOn: rm.PostedOn,
		Reactor:  rm.Reactor.toDomainModel(),
	}
}

func (pm PostModel) toDomainModel() domain.Post {
	reactions := lo.Map(pm.Reactions, func(rm ReactionModel, _ int) domain.Reaction {
		return rm.toDomainModel()
	})

	url, _ := url.Parse(pm.ExternalSourceUrl)
	return domain.Post{
		Id:                pm.Id,
		ExternalSourceUrl: url,
		Reactions:         reactions,
		PostedOn:          pm.PostedOn,
		ChatId:            pm.ChatId,
		Poster:            pm.Poster.toDomainModel(),
	}
}

func MapPostToModel(p domain.Post) PostModel {
	reactionModels := lo.Map(p.Reactions, func(r domain.Reaction, _ int) ReactionModel {
		return MapReactionToModel(r)
	})

	return PostModel{
		Id:                p.Id,
		ExternalSourceUrl: p.ExternalSourceUrl.String(),
		Reactions:         reactionModels,
		PostedOn:          p.PostedOn,
		ChatId:            p.ChatId,
		Poster:            MapUserRefToModel(p.Poster),
	}
}

func MapUserRefToModel(ur domain.UserRef) UserRefModel {
	return UserRefModel{
		DisplayName: ur.DisplayName,
		Id:          ur.Id,
	}
}

func MapReactionToModel(r domain.Reaction) ReactionModel {
	return ReactionModel{
		Text:     r.Text,
		PostedOn: r.PostedOn,
		Reactor:  MapUserRefToModel(r.Reactor),
	}
}
