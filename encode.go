package m3u

import (
	"bytes"
	"fmt"
	"io"
	"maps"
	"net/url"
	"slices"
	"strconv"
)

// Encoder writes M3U playlists to an output stream.
type Encoder struct {
	w   io.Writer
	err error
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes the M3U encoding of p to the stream.
func (e *Encoder) Encode(playlist *Playlist, playlistType PlaylistType) error {
	e.write("#EXTM3U")
	e.writeURLAttr("url-tvg", playlist.TVGURL)
	e.writeURLAttr("x-tvg-url", playlist.XTVGURL)

	// Write extra attributes in a deterministic (sorted) order
	for _, key := range slices.Sorted(maps.Keys(playlist.ExtraAttributes)) {
		e.write(fmt.Sprintf(" %s=\"%s\"", key, playlist.ExtraAttributes[key]))
	}

	e.write("\n")

	// Write tracks
	for _, track := range playlist.Tracks {
		e.write("#EXTINF:" + strconv.FormatFloat(track.Length, 'f', -1, 64))

		if playlistType == M3UPlus {
			e.writeAttr("tvg-id", track.TVGID)
			e.writeAttr("tvg-name", track.TVGName)
			e.writeAttr("tvg-language", track.TVGLanguage)
			e.writeURLAttr("tvg-logo", track.TVGLogo)
			e.writeAttr("group-title", track.GroupTitle)

			// Write extra attributes in a deterministic (sorted) order
			for _, key := range slices.Sorted(maps.Keys(track.ExtraAttributes)) {
				e.write(fmt.Sprintf(" %s=\"%s\"", key, track.ExtraAttributes[key]))
			}
		}

		e.write(fmt.Sprintf(",%s\n", track.Name))

		// Write extra directives
		for _, directive := range track.ExtraDirectives {
			e.write(fmt.Sprintf("%s\n", directive))
		}

		// Write URL
		if track.URL != nil {
			e.write(fmt.Sprintf("%s\n", track.URL.String()))
		}
	}

	return e.err
}

// Marshal returns the M3U encoding of p.
func Marshal(p *Playlist, playlistType PlaylistType) ([]byte, error) {
	var buf bytes.Buffer

	if err := NewEncoder(&buf).Encode(p, playlistType); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// writeAttr writes a quoted key="value" attribute when value is non-nil.
func (e *Encoder) writeAttr(key string, value *string) {
	if value != nil {
		e.write(fmt.Sprintf(" %s=\"%s\"", key, *value))
	}
}

// writeURLAttr writes a quoted key="url" attribute when u is non-nil.
func (e *Encoder) writeURLAttr(key string, u *url.URL) {
	if u != nil {
		e.write(fmt.Sprintf(" %s=\"%s\"", key, u.String()))
	}
}

// write appends s to the stream, retaining the first error encountered.
func (e *Encoder) write(s string) {
	if e.err != nil {
		return
	}

	if _, err := io.WriteString(e.w, s); err != nil {
		e.err = fmt.Errorf("failed to write string: %w", err)
	}
}
