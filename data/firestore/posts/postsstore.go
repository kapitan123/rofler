package posts

import (
	"context"

	"cloud.google.com/go/firestore"
)

const postsColName = "posts"

func GetAll(ctx context.Context, client *firestore.Client) ([]Post, error) {
	docs, err := client.Collection(postsColName).Documents(ctx).GetAll()

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

func Upsert(ctx context.Context, client *firestore.Client, p Post) error {
	doc := client.Collection(postsColName).Doc(p.VideoId)
	_, err := doc.Set(ctx, p)

	return err
}

func GetById(ctx context.Context, client *firestore.Client, videoId string) (Post, bool, error) {
	var p Post
	doc := client.Collection(postsColName).Doc(videoId)
	snap, err := doc.Get(ctx)

	if err != nil {
		return p, false, nil
	}

	if err := snap.DataTo(&p); err != nil {
		return p, false, err
	}

	return p, true, nil
}

func Create(ctx context.Context, client *firestore.Client, p Post) error {
	doc := client.Collection(postsColName).Doc(p.VideoId)
	_, err := doc.Create(ctx, p)

	return err
}

func GetTopRoflerFromPosts(ctx context.Context, client *firestore.Client) (string, int, error) {
	posts, err := GetAll(ctx, client)
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

func Close(client *firestore.Client) {
	client.Close()
}
