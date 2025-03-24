# M3U - Go Library for M3U/IPTV Playlist Parsing

A Go library for parsing and generating M3U playlists with support for IPTV-specific extensions (M3U Plus format).

## Installation

```bash
go get github.com/sherif-fanous/m3u
```

## Usage

### Parsing an M3U Playlist

```go
package main

import (
    "log"
    "strings"

    "github.com/sherif-fanous/m3u"
)

func main() {
    // Parse a playlist from a string
    playlistData := `#EXTM3U url-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml" tvg-url="http://127.0.0.1/epg.xml"
#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1" tvg-country="USA",Channel 1
#EXTVLCOPT:http-referrer=http://example.com/
#EXTVLCOPT:http-user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36
http://127.0.0.1/stream_1
#EXTINF:-1 tvg-id="channel-2" tvg-name="Channel 2" tvg-language="French" tvg-logo="http://127.0.0.1/logos/live_stream_2.png" group-title="Group 2",Channel 2
http://127.0.0.1/stream_2
`

    // Using Unmarshal
    playlist, err := m3u.Unmarshal([]byte(playlistData))
    if err != nil {
        log.Fatalf("Error parsing playlist: %v\n", err)
    }

    // Using a Decoder for streaming
    playlist = &m3u.Playlist{}
    decoder := m3u.NewDecoder(strings.NewReader(playlistData))
    if err := decoder.Decode(playlist); err != nil {
        log.Fatalf("Decoder error: %v\n", err)
    }

    log.Printf("TVG URL: %s\n", playlist.TVGURL.String())
    log.Printf("X TVG URL: %s\n", playlist.XTVGURL.String())

    // Access extra attributes
    for k, v := range playlist.ExtraAttributes {
        log.Printf("Extra Attribute: %s = %s\n", k, v)
    }

    log.Printf("Found %d tracks\n", len(playlist.Tracks))

    // Access track information
    for i, track := range playlist.Tracks {
        log.Printf("Track %d: %s\n", i+1, track.Name)
        if track.TVGID != nil {
            log.Printf("  ID: %s\n", *track.TVGID)
        }
        if track.TVGName != nil {
            log.Printf("  Name: %s\n", *track.TVGName)
        }
        if track.TVGLanguage != nil {
            log.Printf("  Language: %s\n", *track.TVGLanguage)
        }
        if track.TVGLogo != nil {
            log.Printf("  Logo: %s\n", track.TVGLogo.String())
        }
        if track.GroupTitle != nil {
            log.Printf("  Group: %s\n", *track.GroupTitle)
        }
        log.Printf("  URL: %s\n", track.URL)

        // Access extra attributes
        for k, v := range track.ExtraAttributes {
            log.Printf("  Extra Attribute: %s = %s\n", k, v)
        }

        // Access extra directives
        for _, d := range track.ExtraDirectives {
            log.Printf("  Extra Directive: %s\n", d)
        }
    }
}

```

### Creating and Writing an M3U Playlist

```go
package main

import (
    "log"
    "net/url"
    "os"

    "github.com/sherif-fanous/m3u"
)

func makePointer(s string) *string {
    return &s
}

func main() {
    epgURL, _ := url.Parse("http://127.0.0.1/epg.xml")
    tvgLogo, _ := url.Parse("http://127.0.0.1/logos/live_stream_1.png")
    trackURL, _ := url.Parse("http://127.0.0.1/stream_1")

    playlist := &m3u.Playlist{
        TVGURL:  epgURL,
        XTVGURL: epgURL,
        ExtraAttributes: map[string]string{
            "tvg-url": "http://127.0.0.1/epg.xml",
        },
        Tracks: []m3u.Track{
            {
                Length:      -1,
                Name:        "Channel 1",
                TVGID:       makePointer("channel-1"),
                TVGName:     makePointer("Channel 1"),
                TVGLanguage: makePointer("English"),
                TVGLogo:     tvgLogo,
                GroupTitle:  makePointer("Group 1"),
                URL:         trackURL,
                ExtraAttributes: map[string]string{
                    "tvg-country": "USA",
                },
                ExtraDirectives: []string{
                    "#EXTVLCOPT:http-referrer=http://example.com/",
                    "#EXTVLCOPT:http-user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
                },
            },
        },
    }

    // Marshal as M3U Plus format
    data, err := m3u.Marshal(playlist, m3u.M3UPlus)
    if err != nil {
        log.Fatalf("Error marshaling playlist: %v\n", err)
    }

    log.Println(string(data))

    // Write to a file using an Encoder
    file, err := os.Create("playlist.m3u")
    if err != nil {
        log.Fatalf("Error creating file: %v\n", err)
    }
    defer file.Close()

    encoder := m3u.NewEncoder(file)
    if err := encoder.Encode(playlist, m3u.M3UPlus); err != nil {
        log.Fatalf("Error encoding playlist: %v\n", err)
    }
}

```

## M3U Format Support

This library supports two M3U playlist formats:

### Basic M3U Format

```bash
#EXTM3U
#EXTINF:-1,Channel 1
http://127.0.0.1/stream_1
#EXTINF:-1,Channel 2
http://127.0.0.1/stream_2
```

### M3U Plus Format (with Extended Attributes)

```bash
#EXTM3U url-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml"
#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1",Channel 1
http://127.0.0.1/stream_1
#EXTINF:-1 tvg-id="channel-2" tvg-name="Channel 2" tvg-language="French" tvg-logo="http://127.0.0.1/logos/live_stream_2.png" group-title="Group 2",Channel 2
http://127.0.0.1/stream_2
```

### Choosing the Output Format

When generating M3U playlists, you can specify which format to use by setting the `playlistType` parameter in the `Marshal` or `Encode` functions:

```go
// Generate basic M3U format
basicData, err := m3u.Marshal(playlist, m3u.M3U)

// Generate M3U Plus format with extended attributes
extendedData, err := m3u.Marshal(playlist, m3u.M3UPlus)
```

The available format types are:

- `m3u.M3U`: Basic format that only includes track names and URLs
- `m3u.M3UPlus`: Extended format that includes all attributes like tvg-id, tvg-name, tvg-logo, etc.
