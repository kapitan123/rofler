package roflers

type Rofler struct {
	UserName string `firestore:"user_name"`
	Posts    []Post `firestore:"posts"`
}

// Video post stats for future important analytics
type Post struct {
	TiktokId      string   `firestore:"tiktok_id"`
	Url           string   `firestore:"url"`
	ChatLikeCount int      `firestore:"chat_like_count"`
	KeyWords      []string `firestore:"key_words"`
}

func (r *Rofler) AddPost(p Post) {
	r.Posts = append(r.Posts, p)
}

type RoflerStore interface {
	GetAll() ([]Rofler, error)
	Upsert(Rofler) error
	GetByUserName(string) (Rofler, error)
	//GetTop(string, string) (Rofler, error)
}
