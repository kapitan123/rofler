package bot

import "fmt"

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
