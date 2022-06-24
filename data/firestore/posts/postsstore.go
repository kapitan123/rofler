package posts

import (
	"context"

	"cloud.google.com/go/firestore"
)

const postsColName = "posts"

type PostsStorage struct {
	client *firestore.Client
}

func New(client *firestore.Client) *PostsStorage {
	return &PostsStorage{
		client: client,
	}
}

func (s *PostsStorage) GetAll(ctx context.Context) ([]Post, error) {
	docs, err := s.client.Collection(postsColName).Documents(ctx).GetAll()

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

func (s *PostsStorage) Upsert(ctx context.Context, p Post) error {
	doc := s.client.Collection(postsColName).Doc(p.VideoId)
	_, err := doc.Set(ctx, p)

	return err
}

func (s *PostsStorage) GetById(ctx context.Context, videoId string) (Post, bool, error) {
	var p Post
	doc := s.client.Collection(postsColName).Doc(videoId)
	snap, err := doc.Get(ctx)

	if err != nil {
		return p, false, nil
	}

	if err := snap.DataTo(&p); err != nil {
		return p, false, err
	}

	return p, true, nil
}

func (s *PostsStorage) Create(ctx context.Context, client *firestore.Client, p Post) error {
	doc := client.Collection(postsColName).Doc(p.VideoId)
	_, err := doc.Create(ctx, p)

	return err
}

func (s *PostsStorage) GetTopRoflerFromPosts(ctx context.Context) (string, int, error) {
	posts, err := s.GetAll(ctx)
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

func (s *PostsStorage) Close() {
	s.client.Close()
}
