package m3u

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Regular expressions for parsing EXTINF lines
var (
	extm3uLineRegex = regexp.MustCompile(`^#EXTM3U(?:\s+(.*))?$`)
	extinfLineRegex = regexp.MustCompile(`^#EXTINF:(-?\d+\.?\d*)(.*)?,(.*)$`)
	attributeRegex  = regexp.MustCompile(`([\p{L}\p{N}-]+)="([^"]*)"`)
)

// Decoder reads and decodes M3U playlists from an input stream.
type Decoder struct {
	r          *bufio.Reader
	lineNumber int
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r:          bufio.NewReader(r),
		lineNumber: 0,
	}
}

// Decode reads an M3U playlist from its input.
func (d *Decoder) Decode(playlist *Playlist) error {
	*playlist = Playlist{}

	// Read #EXTM3U header
	line, err := d.readLine()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(line, "#EXTM3U") {
		return ErrInvalidPlaylist{
			Message:    "playlist must start with the `#EXTM3U` directive",
			LineNumber: d.lineNumber,
			Line:       line,
		}
	}

	if err := d.parseEXTM3ULine(line, playlist); err != nil {
		return err
	}

	var currentTrack *Track

	for {
		line, err := d.readLine()

		// Handle empty lines with special EOF case
		if line == "" {
			if err == io.EOF {
				// If we reached EOF and there's a pending track, that's an error
				if currentTrack != nil {
					return ErrInvalidPlaylist{
						Message:    "`#EXTINF` directive block must end with a URL",
						LineNumber: d.lineNumber,
						Line:       line,
					}
				}

				break
			}

			if err != nil {
				return err
			}

			continue // Just an empty line, skip it
		}

		if strings.HasPrefix(line, "#EXTINF:") {
			if currentTrack != nil {
				return ErrInvalidPlaylist{
					Message:    "`#EXTINF` directive block must end with a URL",
					LineNumber: d.lineNumber,
					Line:       line,
				}
			}

			// Parse new track
			var track Track
			if err := d.parseEXTINFLine(line, &track); err != nil {
				return err
			}

			currentTrack = &track
			currentTrack.ExtraDirectives = nil
		} else if strings.HasPrefix(line, "#") {
			if currentTrack == nil {
				return ErrInvalidPlaylist{
					Message:    "`#EXTINF` directive must appear before any other directive",
					LineNumber: d.lineNumber,
					Line:       line,
				}
			}
			// It's a directive, add to extra directives
			currentTrack.ExtraDirectives = append(currentTrack.ExtraDirectives, line)
		} else if currentTrack != nil && currentTrack.URL == nil {
			// This should be the URL line for the current track
			parsedURL, err := url.Parse(line)
			if err != nil {
				return ErrInvalidPlaylist{
					Message:    fmt.Sprintf("invalid URL: %v", err),
					LineNumber: d.lineNumber,
					Line:       line,
				}
			}

			currentTrack.URL = parsedURL
			playlist.Tracks = append(playlist.Tracks, *currentTrack)

			// Reset for the next track
			currentTrack = nil
		} else {
			return ErrInvalidPlaylist{
				Message:    "unexpected content",
				LineNumber: d.lineNumber,
				Line:       line,
			}
		}

		// Check for EOF after processing the line
		if err == io.EOF {
			break
		}
	}

	return nil
}

func (d *Decoder) parseEXTINFLine(line string, track *Track) error {
	// Match the basic pattern first
	matches := extinfLineRegex.FindStringSubmatch(line)
	if matches == nil {
		return ErrInvalidPlaylist{
			Message: fmt.Sprintf(
				"malformed `#EXTINF` line: `#EXTNF` line failed to match regex %q",
				extinfLineRegex,
			),
			LineNumber: d.lineNumber,
			Line:       line,
		}
	}

	// Parse duration
	// Intentionally ignore errors here, as the regex should have matched a valid float
	length, _ := strconv.ParseFloat(matches[1], 64)
	track.Length = length

	// Get the attributes part (between duration and name)
	attributes := strings.TrimSpace(matches[2])

	// Extract all attributes
	matchedAttributes := attributeRegex.FindAllStringSubmatch(attributes, -1)
	for _, match := range matchedAttributes {
		key := match[1]
		value := match[2]

		switch key {
		case "tvg-id":
			track.TVGID = &value
		case "tvg-name":
			track.TVGName = &value
		case "tvg-language":
			track.TVGLanguage = &value
		case "tvg-logo":
			if logoURL, err := url.Parse(value); err == nil {
				track.TVGLogo = logoURL
			}
		case "group-title":
			track.GroupTitle = &value
		default:
			if track.ExtraAttributes == nil {
				track.ExtraAttributes = make(map[string]string)
			}
			track.ExtraAttributes[key] = value
		}
	}

	// Set name (after the last comma)
	track.Name = strings.TrimSpace(matches[3])

	return nil
}

func (d *Decoder) parseEXTM3ULine(line string, playlist *Playlist) error {
	// Match the basic pattern first
	matches := extm3uLineRegex.FindStringSubmatch(line)
	if matches == nil {
		return ErrInvalidPlaylist{
			Message: fmt.Sprintf(
				"malformed `#EXTM3U` line: `#EXTM3U` line failed to match regex %q",
				extm3uLineRegex,
			),
			LineNumber: d.lineNumber,
			Line:       line,
		}
	}

	// Get the attributes part
	attributes := strings.TrimSpace(matches[1])

	// Extract all attributes
	matchedAttributes := attributeRegex.FindAllStringSubmatch(attributes, -1)
	for _, match := range matchedAttributes {
		key := match[1]
		value := match[2]

		switch key {
		case "url-tvg":
			if tvgURL, err := url.Parse(value); err == nil {
				playlist.TVGURL = tvgURL
			}
		case "x-tvg-url":
			if xTVGURL, err := url.Parse(value); err == nil {
				playlist.XTVGURL = xTVGURL
			}
		default:
			if playlist.ExtraAttributes == nil {
				playlist.ExtraAttributes = make(map[string]string)
			}
			playlist.ExtraAttributes[key] = value
		}
	}

	return nil
}

func (d *Decoder) readLine() (string, error) {
	d.lineNumber++

	line, err := d.r.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error reading line: %w", err)
	}

	return strings.TrimSpace(line), err
}

// Unmarshal parses the M3U-encoded data and returns the playlist.
func Unmarshal(data []byte) (*Playlist, error) {
	playlist := &Playlist{}

	err := NewDecoder(strings.NewReader(string(data))).Decode(playlist)
	if err != nil {
		return nil, err
	}

	return playlist, nil
}
