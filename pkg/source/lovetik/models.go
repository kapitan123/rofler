package lovetik

type SearchResult struct {
	Author string `json:"author"`
	Links  []Link `json:"links"`
	Vid    string `json:"vid"`
	Desc   string `json:"desc"`
}

type Link struct {
	DownloadAddr string `json:"a"`
	Type         string `json:"t"`
}
