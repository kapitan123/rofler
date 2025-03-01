package tgcommand

// import (
// 	"context"
// 	"testing"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// 	"github.com/kapitan123/telegrofler/internal/storage"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// type mockPostsStorage struct {
// 	posts []storage.Post
// 	err   error
// }

// func (s *mockPostsStorage) GetAllPosts(_ context.Context) ([]storage.Post, error) {
// 	if s.err != nil {
// 		return nil, s.err
// 	}
// 	return s.posts, nil
// }

// type mockMessenger struct {
// 	sendText func(chatId int64, message string) error
// }

// func (m *mockMessenger) SendText(chatId int64, text string) error {
// 	return m.sendText(chatId, text)
// }

// func TestTopRofler_Handle(t *testing.T) {
// 	t.Run("should send message with top posts", func(t *testing.T) {
// 		posts := []storage.Post{
// 			{RoflerUserName: "Bizon", Reactions: []storage.Reaction{{}, {}, {}}},
// 			{RoflerUserName: "Gleb", Reactions: []storage.Reaction{{}, {}, {}, {}, {}}},
// 			{RoflerUserName: "Klim"},
// 		}
// 		s := &mockPostsStorage{posts: posts}
// 		m := &mockMessenger{sendText: func(chatId int64, message string) error {
// 			assert.Equal(t, chatId, int64(228))
// 			// assert.Equal(t, message, formatTopRofler("Gleb", 5))
// 			return nil
// 		}}
// 		cmd := New(m, s)
// 		err := cmd.Handle(context.Background(), &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 228}})
// 		require.NoError(t, err)
// 	})

// 	t.Run("should not send message if there are no posts", func(t *testing.T) {
// 		s := &mockPostsStorage{posts: []storage.Post{}}
// 		m := &mockMessenger{sendText: func(chatId int64, message string) error {
// 			assert.Fail(t, "should not send message")
// 			return nil
// 		}}
// 		cmd := New(m, s)
// 		err := cmd.Handle(context.Background(), &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 228}})
// 		require.NoError(t, err)
// 	})
// }
