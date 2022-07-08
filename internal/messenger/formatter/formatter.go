package formatter

import (
	"fmt"

	"github.com/kapitan123/telegrofler/internal/sortedmap"
)

// AK TODO makes sense to remove this struct
type Formatter struct {
}

func New() *Formatter {
	return &Formatter{}
}

func (f *Formatter) FormatAsDescendingList(m map[string]int, format string) string {
	listMeassge := ""

	sm := sortedmap.Sort(m)

	for _, pair := range sm {
		listMeassge += formatLine(format, pair.Key, pair.Value)
	}
	return listMeassge
}

func formatLine(format string, username string, score int) string {
	return fmt.Sprintf(format, username, score)
}
