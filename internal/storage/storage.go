package storage

import "cloud.google.com/go/firestore"

type Storage struct {
	client *firestore.Client
}

func New(client *firestore.Client) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) Close() error {
	return s.client.Close()
}
