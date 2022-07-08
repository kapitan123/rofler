package shortsget

type Response struct {
	VideoDetails VideoDetails `json:"videoDetails"`
	Formats      []Format     `json:"formats"`
}

type VideoDetails struct {
	Title   string `json:"title"`
	VideoId string `json:"videoId"`
}

type Format struct {
	Itag         int    `json:"itag"`
	MimeType     string `json:"mimeType"`     //"mimeType": "video/mp4; codecs=\"avc1.42001E, mp4a.40.2\""
	QualityLabel string `json:"qualityLabel"` //"qualityLabel": "240p"
	Codec        string `json:"codec"`
}
