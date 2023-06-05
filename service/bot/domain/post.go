package domain

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

type (
	Post struct {
		Id                string
		ExternalSourceUrl *url.URL
		OriginalMessageId MessageId
		Type              MediaType
		Reactions         []Reaction
		PostedOn          time.Time
		ChatId            ChatId
		Poster            UserRef
	}

	Reaction struct {
		Text     string
		PostedOn time.Time
		Reactor  UserRef
	}
)

func NewPost(mediaId string, poster UserRef, chatId ChatId, mediaType MediaType) Post {
	return Post{
		Id:        mediaId,
		Type:      mediaType,
		Poster:    poster,
		ChatId:    chatId,
		Reactions: []Reaction{},
		PostedOn:  time.Now(),
	}
}

func NewPostFromExternalSource(externalSourceUrl *url.URL, poster UserRef, chatId ChatId, messageId MessageId) Post {
	return Post{
		Id:                uuid.NewString(),
		Type:              Video,
		OriginalMessageId: messageId,
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

// AK TODO should group scores
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
