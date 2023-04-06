package storage

type UserRef struct {
	DisplayName string `firestore:"user_name"`
	Id          int64  `firestore:"user_id"`
}
