package problem

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Problem in RFC-7807 format
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func New(url, title string, status int, detail, instance string) *Problem {
	return &Problem{url, title, status, detail, instance}
}

// Writes problem as HTTP response.
func (p *Problem) WriteTo(resp http.ResponseWriter) {
	log.Info("API error %s", p.Error())

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(p.Status)
	json.NewEncoder(resp).Encode(p)
}

// Implements error interface
func (p Problem) Error() string {
	return fmt.Sprintf("Problem: Type: '%s', Title: '%s', Status: '%d', Detail: '%s', Instance: '%s'",
		p.Type, p.Title, p.Status, p.Detail, p.Instance)
}
