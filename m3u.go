package m3u

import (
	"net/url"
)

// PlaylistType defines the format for encoding M3U playlists.
type PlaylistType string

const (
	// M3U is the basic M3U format.
	M3U PlaylistType = "M3U"
	// M3UPlus is the extended M3U format with additional attributes.
	M3UPlus PlaylistType = "M3UPlus"
)

// Playlist represents an M3U playlist.
type Playlist struct {
	Tracks []Track
}

// Track represents a single entry in an M3U playlist.
type Track struct {
	Length          float64
	Name            string
	TVGID           *string
	TVGName         *string
	TVGLanguage     *string
	TVGLogo         *url.URL
	GroupTitle      *string
	URL             *url.URL
	ExtraAttributes map[string]string
	ExtraDirectives []string
}
