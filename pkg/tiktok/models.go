package tiktok

type Stats struct {
	DiggCount    int `json:"diggCount"`
	ShareCount   int `json:"shareCount"`
	CommentCount int `json:"commentCount"`
	PlayCount    int `json:"playCount"`
}

type Video struct {
	Id           string `json:"id"`
	Height       int    `json:"height"`
	Width        int    `json:"width"`
	Duration     int    `json:"duration"`
	DownloadAddr string `json:"downloadAddr"`
	Format       string `json:"format"`
	CodecType    string `json:"codecType"`
}

type Item struct {
	Id     string `json:"id"`
	Desc   string `json:"desc"`
	Author string `json:"author"`
	IsAd   bool   `json:"isAd"`
	Stats  Stats  `json:"stats"`
	Video  Video  `json:"video"`
}
