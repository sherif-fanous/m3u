package m3u

import "fmt"

type InvalidPlaylistError struct {
	Message    string
	LineNumber int
	Line       string
}

func (e InvalidPlaylistError) Error() string {
	return fmt.Sprintf("invalid m3u playlist: line %d: `%s`: %s", e.LineNumber, e.Line, e.Message)
}
