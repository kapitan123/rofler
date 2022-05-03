package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/kapitan123/telegrofler/config"

	log "github.com/sirupsen/logrus"
)

const postsColName = "posts"

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

func (rs *PostsStore) GetAllPosts() ([]Post, error) {
	docs, err := rs.postsCol.Documents(*rs.ctx).GetAll()

	if err != nil {
		return nil, err
	}

	posts := []Post{}
	for _, doc := range docs {
		p := Post{}
		doc.DataTo(&p)
		posts = append(posts, p)
	}

	return posts, nil
}

func (rs *PostsStore) Upsert(p Post) error {
	doc := rs.postsCol.Doc(p.VideoId)
	_, err := doc.Set(*rs.ctx, p)

	return err
}

func (rs *PostsStore) GetById(videoId string) (Post, bool, error) {
	var p Post
	doc := rs.postsCol.Doc(videoId)
	snap, err := doc.Get(*rs.ctx)

	if err != nil {
		return p, false, nil
	}

	if err := snap.DataTo(&p); err != nil {
		return p, false, err
	}

	return p, true, nil
}

func (rs *PostsStore) Create(p Post) error {
	doc := rs.postsCol.Doc(p.VideoId)
	_, err := doc.Create(*rs.ctx, p)

	return err
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

func createClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, config.ProjectId)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err) // AK TODO shouldn't be fatal
	}
	return client
}
