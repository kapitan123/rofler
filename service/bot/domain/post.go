package domain

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/kapitan123/telegrofler/service/bot/domain/media"
)

type (
	Post struct {
		Id                string
		ExternalSourceUrl *url.URL
		Type              media.Type
		Reactions         []Reaction
		PostedOn          time.Time
		ChatId            int64
		Poster            UserRef
	}

	Reaction struct {
		ReactToMessageId int
		Text             string
		PostedOn         time.Time
		Reactor          UserRef
	}
)

func NewPost(mediaId string, poster UserRef, chatId int64) Post {
	return Post{
		Id:        mediaId,
		Poster:    poster,
		ChatId:    chatId,
		Reactions: []Reaction{},
		PostedOn:  time.Now(),
	}
}

func NewPostFromExternalSource(externalSourceUrl *url.URL, poster UserRef, chatId int64) Post {
	return Post{
		Id:                uuid.NewString(),
		ExternalSourceUrl: externalSourceUrl,
		Poster:            poster,
		ChatId:            chatId,
		Reactions:         []Reaction{},
		PostedOn:          time.Now(),
	}
}

func (p *Post) AddReaction(reaction Reaction) {
	p.Reactions = append(p.Reactions, reaction)
}

// type postsCollection []Post

// // AK TODO move counting ab list buildeng logic here
// // groups posts by user
// func (posts *postsCollection) GroupReactionsByUser() map[string]int {

// 	roflerScores := map[UserRef]int{}
// 	for _, p := range posts {
// 		roflerScores[p.Poster] += len(p.Reactions)
// 	}

// 	names := map[string]int{}
// 	for k, v := range roflerScores {
// 		names[k.DisplayName] = v
// 	}
// 	return names
// }
