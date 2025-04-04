package m3u

import (
	"fmt"
	"io"
	"strings"
)

// Encoder writes M3U playlists to an output stream.
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes the M3U encoding of p to the stream.
func (e *Encoder) Encode(playlist *Playlist, playlistType PlaylistType) error {
	if _, err := e.writeString("#EXTM3U"); err != nil {
		return err
	}

	// Write TVG-URL if present
	if playlist.TVGURL != nil {
		if _, err := e.writeString(fmt.Sprintf(" url-tvg=\"%s\"", playlist.TVGURL.String())); err != nil {
			return err
		}
	}

	// Write XTVG-URL if present
	if playlist.XTVGURL != nil {
		if _, err := e.writeString(fmt.Sprintf(" x-tvg-url=\"%s\"", playlist.XTVGURL.String())); err != nil {
			return err
		}
	}

	// Write extra attributes
	for key, value := range playlist.ExtraAttributes {
		if _, err := e.writeString(fmt.Sprintf(" %s=\"%s\"", key, value)); err != nil {
			return err
		}
	}

	// Write newline after #EXTM3U line
	if _, err := e.writeString("\n"); err != nil {
		return err
	}

	// Write tracks
	for _, track := range playlist.Tracks {
		// Write #EXTINF line
		if playlistType == M3UPlus {
			// M3UPlus format with attributes
			if _, err := e.writeString(fmt.Sprintf("#EXTINF:%.0f", track.Length)); err != nil {
				return err
			}

			// Write TVG-ID if present
			if track.TVGID != nil {
				if _, err := e.writeString(fmt.Sprintf(" tvg-id=\"%s\"", *track.TVGID)); err != nil {
					return err
				}
			}

			// Write TVG-Name if present
			if track.TVGName != nil {
				if _, err := e.writeString(fmt.Sprintf(" tvg-name=\"%s\"", *track.TVGName)); err != nil {
					return err
				}
			}

			// Write TVG-Language if present
			if track.TVGLanguage != nil {
				if _, err := e.writeString(fmt.Sprintf(" tvg-language=\"%s\"", *track.TVGLanguage)); err != nil {
					return err
				}
			}

			// Write TVG-Logo if present
			if track.TVGLogo != nil {
				if _, err := e.writeString(fmt.Sprintf(" tvg-logo=\"%s\"", track.TVGLogo.String())); err != nil {
					return err
				}
			}

			// Write Group-Title if present
			if track.GroupTitle != nil {
				if _, err := e.writeString(fmt.Sprintf(" group-title=\"%s\"", *track.GroupTitle)); err != nil {
					return err
				}
			}

			// Write extra attributes
			for key, value := range track.ExtraAttributes {
				if _, err := e.writeString(fmt.Sprintf(" %s=\"%s\"", key, value)); err != nil {
					return err
				}
			}

		} else {
			// M3U format without attributes
			if _, err := e.writeString(fmt.Sprintf("#EXTINF:%.0f", track.Length)); err != nil {
				return err
			}
		}

		// Add track name
		if _, err := e.writeString(fmt.Sprintf(",%s\n", track.Name)); err != nil {
			return err
		}

		// Write extra directives
		for _, directive := range track.ExtraDirectives {
			if _, err := e.writeString(fmt.Sprintf("%s\n", directive)); err != nil {
				return err
			}
		}

		// Write URL
		if track.URL != nil {
			if _, err := e.writeString(fmt.Sprintf("%s\n", track.URL.String())); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Encoder) writeString(s string) (string, error) {
	if _, err := io.WriteString(e.w, s); err != nil {
		return "", fmt.Errorf("failed to write string: %w", err)
	}
	return s, nil
}

// Marshal returns the M3U encoding of p.
func Marshal(p *Playlist, playlistType PlaylistType) ([]byte, error) {
	var buf strings.Builder

	err := NewEncoder(&buf).Encode(p, playlistType)
	if err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}
