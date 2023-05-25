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

const postsCollectionName = "posts"

type FirestorePostsRepository struct {
	client *firestore.Client
	//postFactory     post.Factory
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

	doc := s.postsCollection().Doc(p.VideoId)
	_, err := doc.Set(ctx, p)

	return err
}

type (
	PostModel struct {
		VideoId   string          `firestore:"video_id"`
		Source    string          `firestore:"source"`
		Url       string          `firestore:"url"`
		Reactions []ReactionModel `firestore:"reactions"`
		PostedOn  time.Time       `firestore:"posted_on"`
		ChatId    int64           `firestore:"chat_id"`
		UserRef   UserRefModel    `firestore:"user_ref"`
	}

	ReactionModel struct {
		ReactToMessageId int          `firestore:"react_to_message_id"` // RepllyToMessage.ID not the update.Message.ID
		Sender           UserRefModel `firestore:"sender"`
		Text             string       `firestore:"text"`
		PostedOn         time.Time    `firestore:"posted_on"`
		ReactorUserRef   UserRefModel `firestore:"reactor_user_ref"`
	}

	UserRefModel struct {
		DisplayName string `firestore:"user_name"`
		Id          int64  `firestore:"user_id"`
	}
)

func (urm UserRefModel) toDomainModel() domain.UserRef {
	return domain.UserRef{
		urm.DisplayName,
		urm.Id,
	}
}

func (rm ReactionModel) toDomainModel() domain.Reaction {
	return domain.Reaction{
		ReactToMessageId: rm.ReactToMessageId,
		Sender:           rm.Sender.toDomainModel(),
		Text:             rm.Text,
		PostedOn:         rm.PostedOn,
		ReactorUserRef:   rm.ReactorUserRef.toDomainModel(),
	}
}

func (pm PostModel) toDomainModel() domain.Post {

	reactions := lo.Map(pm.Reactions, func(rm ReactionModel, _ int) domain.Reaction {
		return rm.toDomainModel()
	})

	return domain.Post{
		VideoId:   pm.VideoId,
		Source:    pm.Source,
		Url:       pm.Url,
		Reactions: reactions,
		PostedOn:  pm.PostedOn,
		ChatId:    pm.ChatId,
		UserRef:   pm.UserRef.toDomainModel(),
	}
}

func MapPostToModel(p domain.Post) PostModel {
	reactionModels := lo.Map(p.Reactions, func(r domain.Reaction, _ int) ReactionModel {
		return MapReactionToModel(r)
	})

	return PostModel{
		VideoId:   p.VideoId,
		Source:    p.Source,
		Url:       p.Url,
		Reactions: reactionModels,
		PostedOn:  p.PostedOn,
		ChatId:    p.ChatId,
		UserRef:   MapUserRefToModel(p.UserRef),
	}
}

func MapUserRefToModel(ur domain.UserRef) UserRefModel {
	return UserRefModel{
		ur.DisplayName,
		ur.Id,
	}
}

func MapReactionToModel(r domain.Reaction) ReactionModel {
	return ReactionModel{
		ReactToMessageId: r.ReactToMessageId,
		Sender:           MapUserRefToModel(r.Sender),
		Text:             r.Text,
		PostedOn:         r.PostedOn,
		ReactorUserRef:   MapUserRefToModel(r.ReactorUserRef),
	}
}
