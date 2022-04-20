package bot

import "fmt"

type TikTokVideoPost struct {
	Sender            string
	ChatId            int64
	OriginalMessageId int
	Url               string
	VideoData         VideoData
}

type VideoData struct {
	Id         string
	Duration   int
	Title      string
	Payload    []byte
	LikesCount int
}

func (tp *TikTokVideoPost) GetCaption() string {
	return fmt.Sprintf("<b>Rofler</b>: @%s \n<b>Title</b>: %s \n<b>:heart:</b>: %d", tp.Sender, tp.VideoData.Title, tp.VideoData.LikesCount)
}
