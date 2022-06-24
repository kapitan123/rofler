package bot

import "fmt"

type (
	SourceVideoPost struct {
		Sender            string
		ChatId            int64
		OriginalMessageId int
		Url               string
		VideoData         VideoData
	}
	VideoData struct {
		Id         string
		Duration   int
		Title      string
		Payload    []byte
		LikesCount int
	}
)

type (
	ReplyToMediaPost struct {
		VideoId string
		Details Details
	}

	Details struct {
		MessageId int // RepllyToMessage.ID not the update.Message.ID
		Sender    string
		Text      string
	}
)

// AK TODO this method should be extracted
func (tp *SourceVideoPost) GetCaption() string {
	return fmt.Sprintf("<b>Rofler:</b> 🔥@%s🔥\n<b>Title</b>: %s", tp.Sender, tp.VideoData.Title)
}
