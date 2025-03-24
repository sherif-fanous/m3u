package m3u_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sherif-fanous/m3u"
)

func makePointer[T any](t T) *T {
	return &t
}

func makeURL(t *testing.T, s string) *url.URL {
	t.Helper()

	u, err := url.Parse(s)
	if err != nil {
		t.Fatal(err)
	}

	return u
}

func TestDecodeM3U(t *testing.T) {
	t.Parallel()

	input := `#EXTM3U
#EXTINF:-1,Channel 1 🙃
http://127.0.0.1/stream_1
#EXTINF:-1,Channel 2 🙃
http://127.0.0.1/stream_2
`

	playlist := &m3u.Playlist{}
	err := m3u.NewDecoder(strings.NewReader(input)).Decode(playlist)
	if err != nil {
		t.Fatalf("Failed to decode M3U: %v", err)
	}

	expectedPlaylist := &m3u.Playlist{
		Tracks: []m3u.Track{
			{
				Length: -1,
				Name:   "Channel 1 🙃",
				URL:    makeURL(t, "http://127.0.0.1/stream_1"),
			},
			{
				Length: -1,
				Name:   "Channel 2 🙃",
				URL:    makeURL(t, "http://127.0.0.1/stream_2"),
			},
		},
	}

	if diff := cmp.Diff(playlist, expectedPlaylist); diff != "" {
		t.Error(diff)
	}
}

func TestDecodeM3UPlus(t *testing.T) {
	t.Parallel()

	input := `#EXTM3U url-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml"
#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1 🙃",Channel 1 🙃
http://127.0.0.1/stream_1
#EXTINF:-1 tvg-id="channel-2" tvg-name="Channel 2" tvg-language="French" tvg-logo="http://127.0.0.1/logos/live_stream_2.png" group-title="Group 2 🙃",Channel 2 🙃
http://127.0.0.1/stream_2
`

	playlist := &m3u.Playlist{}
	err := m3u.NewDecoder(strings.NewReader(input)).Decode(playlist)
	if err != nil {
		t.Fatalf("Failed to decode M3U Plus: %v", err)
	}

	expectedPlaylist := &m3u.Playlist{
		TVGURL:  makeURL(t, "http://127.0.0.1/epg.xml"),
		XTVGURL: makeURL(t, "http://127.0.0.1/epg.xml"),
		Tracks: []m3u.Track{
			{
				Length:      -1,
				Name:        "Channel 1 🙃",
				TVGID:       makePointer("channel-1"),
				TVGName:     makePointer("Channel 1"),
				TVGLanguage: makePointer("English"),
				TVGLogo:     makeURL(t, "http://127.0.0.1/logos/live_stream_1.png"),
				GroupTitle:  makePointer("Group 1 🙃"),
				URL:         makeURL(t, "http://127.0.0.1/stream_1"),
			},
			{
				Length:      -1,
				Name:        "Channel 2 🙃",
				TVGID:       makePointer("channel-2"),
				TVGName:     makePointer("Channel 2"),
				TVGLanguage: makePointer("French"),
				TVGLogo:     makeURL(t, "http://127.0.0.1/logos/live_stream_2.png"),
				GroupTitle:  makePointer("Group 2 🙃"),
				URL:         makeURL(t, "http://127.0.0.1/stream_2"),
			},
		},
	}

	if diff := cmp.Diff(playlist, expectedPlaylist); diff != "" {
		t.Error(diff)
	}
}

func TestEncodeM3U(t *testing.T) {
	t.Parallel()

	playlist := &m3u.Playlist{
		Tracks: []m3u.Track{
			{
				Length: -1,
				Name:   "Channel 1 🙃",
				URL:    makeURL(t, "http://127.0.0.1/stream_1"),
			},
			{
				Length: -1,
				Name:   "Channel 2 🙃",
				URL:    makeURL(t, "http://127.0.0.1/stream_2"),
			},
		},
	}

	data, err := m3u.Marshal(playlist, m3u.M3U)
	if err != nil {
		t.Fatalf("Failed to marshal M3U: %v", err)
	}

	expected := `#EXTM3U
#EXTINF:-1,Channel 1 🙃
http://127.0.0.1/stream_1
#EXTINF:-1,Channel 2 🙃
http://127.0.0.1/stream_2
`
	if string(data) != expected {
		t.Fatalf("Expected:\n%s\nGot:\n%s", expected, string(data))
	}
}

func TestEncodeM3UPlus(t *testing.T) {
	t.Parallel()

	playlist := &m3u.Playlist{
		TVGURL:  makeURL(t, "http://127.0.0.1/epg.xml"),
		XTVGURL: makeURL(t, "http://127.0.0.1/epg.xml"),
		Tracks: []m3u.Track{
			{
				Length:      -1,
				Name:        "Channel 1 🙃",
				TVGID:       makePointer("channel-1"),
				TVGName:     makePointer("Channel 1"),
				TVGLanguage: makePointer("English"),
				TVGLogo:     makeURL(t, "http://127.0.0.1/logos/live_stream_1.png"),
				GroupTitle:  makePointer("Group 1 🙃"),
				URL:         makeURL(t, "http://127.0.0.1/stream_1"),
				ExtraAttributes: map[string]string{
					"tvg-country": "USA",
				},
				ExtraDirectives: []string{
					"#EXTVLCOPT:http-referrer=http://example.com/",
					"#EXTVLCOPT:http-user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
				},
			},
			{
				Length:      -1,
				Name:        "Channel 2 🙃",
				TVGID:       makePointer("channel-2"),
				TVGName:     makePointer("Channel 2"),
				TVGLanguage: makePointer("French"),
				TVGLogo:     makeURL(t, "http://127.0.0.1/logos/live_stream_2.png"),
				GroupTitle:  makePointer("Group 2 🙃"),
				URL:         makeURL(t, "http://127.0.0.1/stream_2"),
			},
		},
	}

	data, err := m3u.Marshal(playlist, m3u.M3UPlus)
	if err != nil {
		t.Fatalf("Failed to marshal M3U Plus: %v", err)
	}

	expected := `#EXTM3U url-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml"
#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1 🙃" tvg-country="USA",Channel 1 🙃
#EXTVLCOPT:http-referrer=http://example.com/
#EXTVLCOPT:http-user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36
http://127.0.0.1/stream_1
#EXTINF:-1 tvg-id="channel-2" tvg-name="Channel 2" tvg-language="French" tvg-logo="http://127.0.0.1/logos/live_stream_2.png" group-title="Group 2 🙃",Channel 2 🙃
http://127.0.0.1/stream_2
`
	if string(data) != expected {
		t.Fatalf("Expected:\n%s\nGot:\n%s", expected, string(data))
	}
}

func TestExtraAttributesAndDirectives(t *testing.T) {
	t.Parallel()

	input := `#EXTM3U url-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml" tvg-url="http://127.0.0.1/epg.xml"
#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-country="USA" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1",Channel 1
#EXTVLCOPT:http-referrer=http://example.com/
#EXTVLCOPT:http-user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36
http://127.0.0.1/stream_1
`

	playlist := &m3u.Playlist{}
	err := m3u.NewDecoder(strings.NewReader(input)).Decode(playlist)
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	expectedPlaylist := &m3u.Playlist{
		TVGURL:  makeURL(t, "http://127.0.0.1/epg.xml"),
		XTVGURL: makeURL(t, "http://127.0.0.1/epg.xml"),
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
				TVGLogo:     makeURL(t, "http://127.0.0.1/logos/live_stream_1.png"),
				GroupTitle:  makePointer("Group 1"),
				URL:         makeURL(t, "http://127.0.0.1/stream_1"),
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

	if diff := cmp.Diff(playlist, expectedPlaylist); diff != "" {
		t.Error(diff)
	}
}

func TestRoundTrip(t *testing.T) {
	t.Parallel()

	input := `#EXTM3U url-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml"
#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1 🙃",Channel 1 🙃
http://127.0.0.1/stream_1
#EXTINF:-1 tvg-id="channel-2" tvg-name="Channel 2" tvg-language="French" tvg-logo="http://127.0.0.1/logos/live_stream_2.png" group-title="Group 2 🙃",Channel 2 🙃
http://127.0.0.1/stream_2
`

	playlist, err := m3u.Unmarshal([]byte(input))
	if err != nil {
		t.Fatalf("Failed to unmarshal M3U: %v", err)
	}

	output, err := m3u.Marshal(playlist, m3u.M3UPlus)
	if err != nil {
		t.Fatalf("Failed to marshal M3U: %v", err)
	}

	playlist2, err := m3u.Unmarshal(output)
	if err != nil {
		t.Fatalf("Failed to unmarshal second M3U: %v", err)
	}

	if diff := cmp.Diff(playlist, playlist2); diff != "" {
		t.Error(diff)
	}
}

func TestErrInvalidPlaylistEXTM3U(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name: "missing `#EXTM3U` directive",
			input: `#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1",Channel 1
http://127.0.0.1/stream_1
#EXTINF:-1 tvg-id="channel-2" tvg-name="Channel 2" tvg-language="French" tvg-logo="http://127.0.0.1/logos/live_stream_2.png" group-title="Group 2",Channel 2
http://127.0.0.1/stream_2
`,
			expectedError: "playlist must start with the `#EXTM3U` directive",
		},
		{
			name: "malformed `#EXTM3U` line",
			input: `#EXTM3Uurl-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml"
#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1",Channel 1
http://127.0.0.1/stream_1
#EXTINF:-1 tvg-id="channel-2" tvg-name="Channel 2" tvg-language="French" tvg-logo="http://127.0.0.1/logos/live_stream_2.png" group-title="Group 2",Channel 2
http://127.0.0.1/stream_2
`,
			expectedError: "malformed `#EXTM3U` line: `#EXTM3U` line failed to match regex",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := m3u.Unmarshal([]byte(test.input))
			if err == nil {
				t.Fatal("Expected an ErrInvalidPlaylist error")
			}

			_, ok := err.(m3u.ErrInvalidPlaylist)
			if !ok {
				t.Fatalf("Expected an ErrInvalidPlaylist error, got: %v", err)
			}
			if !strings.Contains(err.Error(), test.expectedError) {
				t.Fatalf("Expected error message to contain %s, got: %v", test.expectedError, err)
			}
		})
	}
}

func TestErrInvalidPlaylistEXTINFFirst(t *testing.T) {
	t.Parallel()

	input := `#EXTM3U url-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml"
#EXTVLCOPT:http-referrer=http://example.com/
#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-country="USA" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1",Channel 1
#EXTVLCOPT:http-user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36
http://127.0.0.1/stream_1
`

	_, err := m3u.Unmarshal([]byte(input))
	if err == nil {
		t.Fatal("Expected an ErrInvalidPlaylist error")
	}

	_, ok := err.(m3u.ErrInvalidPlaylist)
	if !ok {
		t.Fatalf("Expected an ErrInvalidPlaylist error, got: %v", err)
	}
	if !strings.Contains(
		err.Error(),
		"`#EXTINF` directive must appear before any other directive",
	) {
		t.Fatalf(
			"Expected error message to contain `#EXTINF` directive must appear before any other directive, got: %v",
			err,
		)
	}
}

func TestErrInvalidPlaylistInvalidURL(t *testing.T) {
	t.Parallel()

	input := `#EXTM3U url-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml"
#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-country="USA" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1",Channel 1
#EXTVLCOPT:http-referrer=http://example.com/
#EXTVLCOPT:http-user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36
http://127.0.0.1/stream_1%
`

	_, err := m3u.Unmarshal([]byte(input))
	if err == nil {
		t.Fatal("Expected an ErrInvalidPlaylist error")
	}

	_, ok := err.(m3u.ErrInvalidPlaylist)
	if !ok {
		t.Fatalf("Expected an ErrInvalidPlaylist error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "invalid URL") {
		t.Fatalf("Expected error message to contain invalid URL, got: %v", err)
	}
}

func TestErrInvalidPlaylistMalformedEXTINFLine(t *testing.T) {
	t.Parallel()

	input := `#EXTM3U url-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml"
#EXTINF:NotANumber tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-country="USA" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1",Channel 1
http://127.0.0.1/stream_1
`

	_, err := m3u.Unmarshal([]byte(input))
	if err == nil {
		t.Fatal("Expected an ErrInvalidPlaylist error")
	}

	_, ok := err.(m3u.ErrInvalidPlaylist)
	if !ok {
		t.Fatalf("Expected an ErrInvalidPlaylist error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "malformed `#EXTINF` line") {
		t.Fatalf("Expected error message to contain malformed `#EXTINF` line, got: %v", err)
	}
}

func TestErrInvalidPlaylistMissingURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{
			name: "missing first URL",
			input: `#EXTM3U url-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml"
#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1 🙃",Channel 1 🙃
#EXTINF:-1 tvg-id="channel-2" tvg-name="Channel 2" tvg-language="French" tvg-logo="http://127.0.0.1/logos/live_stream_2.png" group-title="Group 2 🙃",Channel 2 🙃
http://127.0.0.1/stream_2
#EXTINF:-1 tvg-id="channel-3" tvg-name="Channel 3" tvg-language="French" tvg-logo="http://127.0.0.1/logos/live_stream_2.png" group-title="Group 3 🙃",Channel 3 🙃
http://127.0.0.1/stream_3
`,
		},
		{
			name: "missing intermediate URL",
			input: `#EXTM3U url-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml"
#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1 🙃",Channel 1 🙃
http://127.0.0.1/stream_1
#EXTINF:-1 tvg-id="channel-2" tvg-name="Channel 2" tvg-language="French" tvg-logo="http://127.0.0.1/logos/live_stream_2.png" group-title="Group 2 🙃",Channel 2 🙃
#EXTINF:-1 tvg-id="channel-3" tvg-name="Channel 3" tvg-language="French" tvg-logo="http://127.0.0.1/logos/live_stream_2.png" group-title="Group 3 🙃",Channel 3 🙃
http://127.0.0.1/stream_3
`,
		},
		{
			name: "missing final URL",
			input: `#EXTM3U url-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml"
#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1 🙃",Channel 1 🙃
http://127.0.0.1/stream_1
#EXTINF:-1 tvg-id="channel-2" tvg-name="Channel 2" tvg-language="French" tvg-logo="http://127.0.0.1/logos/live_stream_2.png" group-title="Group 2 🙃",Channel 2 🙃
http://127.0.0.1/stream_2
#EXTINF:-1 tvg-id="channel-3" tvg-name="Channel 3" tvg-language="French" tvg-logo="http://127.0.0.1/logos/live_stream_2.png" group-title="Group 3 🙃",Channel 3 🙃
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := m3u.Unmarshal([]byte(test.input))
			if err == nil {
				t.Fatal("Expected an ErrInvalidPlaylist error")
			}

			_, ok := err.(m3u.ErrInvalidPlaylist)
			if !ok {
				t.Fatalf("Expected an ErrInvalidPlaylist error, got: %v", err)
			}
			if !strings.Contains(err.Error(), "`#EXTINF` directive block must end with a URL") {
				t.Fatalf(
					"Expected error message to contain `#EXTINF` directive block must end with a URL, got: %v",
					err,
				)
			}
		})
	}
}

func TestErrInvalidPlaylistUnexpectedContent(t *testing.T) {
	t.Parallel()

	input := `#EXTM3U url-tvg="http://127.0.0.1/epg.xml" x-tvg-url="http://127.0.0.1/epg.xml"
#EXTINF:-1 tvg-id="channel-1" tvg-name="Channel 1" tvg-language="English" tvg-country="USA" tvg-logo="http://127.0.0.1/logos/live_stream_1.png" group-title="Group 1",Channel 1
http://127.0.0.1/stream_1
Unexpected content
`

	_, err := m3u.Unmarshal([]byte(input))
	if err == nil {
		t.Fatal("Expected an ErrInvalidPlaylist error")
	}

	_, ok := err.(m3u.ErrInvalidPlaylist)
	if !ok {
		t.Fatalf("Expected an ErrInvalidPlaylist error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "unexpected content") {
		t.Fatalf("Expected error message to contain unexpected content, got: %v", err)
	}
}
