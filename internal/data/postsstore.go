package data

// AK TODO rename package
import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/kapitan123/telegrofler/internal/data/model"

	log "github.com/sirupsen/logrus"
)

const postsColName = "posts"

// AK TODO need to separate in different data stores
type PostsStore struct {
	ctx      *context.Context
	client   *firestore.Client
	postsCol *firestore.CollectionRef
}

func NewPostsStore() *PostsStore {
	ctx := context.Background()
	client := createClient(ctx)
	return &PostsStore{
		&ctx,
		client,
		client.Collection(postsColName),
	}
}

func (rs *PostsStore) GetAllPosts() ([]model.Post, error) {
	docs, err := rs.postsCol.Documents(*rs.ctx).GetAll()

	if err != nil {
		return nil, err
	}

	posts := []model.Post{}
	for _, doc := range docs {
		p := model.Post{}
		doc.DataTo(&p)
		posts = append(posts, p)
	}

	return posts, nil
}

func (rs *PostsStore) Upsert(p model.Post) error {
	doc := rs.postsCol.Doc(p.VideoId)
	_, err := doc.Set(*rs.ctx, p)

	return err
}

func (rs *PostsStore) Create(p model.Post) error {
	doc := rs.postsCol.Doc(p.VideoId)
	_, err := doc.Create(*rs.ctx, p)

	return err
}

func (rs *PostsStore) AddReactionToPost(vr model.VideoReaction) error {
	posts, err := rs.GetAllPosts()
	if err != nil {
		return err
	}

	for _, p := range posts {
		if p.VideoId == vr.VideoId {
			p.AddReaction(vr.Sender, vr.Text, vr.MessageId)

			err = rs.Upsert(p)

			log.Infof("Reaction was saved for %s from %s: ", vr.VideoId, vr.Sender)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (rs *PostsStore) GetTopRoflerFromPosts() (string, int, error) {
	posts, err := rs.GetAllPosts()
	if err != nil {
		return "", 0, err
	}

	roflerScores := map[string]int{}

	for _, p := range posts {
		roflerScores[p.RoflerUserName] += len(p.Reactions)
	}

	maxUserName, max := "", 0
	for username, score := range roflerScores {
		if max < score {
			max = score
			maxUserName = username
		}
	}

	return maxUserName, max, nil
}

func (rs *PostsStore) Close() {
	rs.client.Close()
}
