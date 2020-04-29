package eztv

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

// TestShowNotFound tests the ErrShowNotFound error
func TestShowNotFound(t *testing.T) {
	var requestURI string
	rawHTMLResponse := `{ "torrents_count": 280450, "limit": 30, "page": 1, "torrents": [] }`
	// Fake server with a fake answer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI = r.RequestURI
		fmt.Fprintln(w, rawHTMLResponse)
	}))
	defer ts.Close()
	endpoint = ts.URL

	_, err := GetShowTorrents("fakeID")
	if err == nil {
		t.Fatalf("got no error, expected %q", ErrShowNotFound)
	}
	expectedRequestURI := "/api/get-torrents?imdb_id=fakeID&limit=100&page=1"
	if requestURI != expectedRequestURI {
		t.Fatalf("expected URL %q, got %q", expectedRequestURI, requestURI)
	}

	if err != ErrShowNotFound {
		t.Fatalf("expected %q, got %q", ErrShowNotFound, err)
	}
}

// TestBadGetShowDetails tests GetShowTorrents with missing arguments
func TestBadGetShowDetails(t *testing.T) {
	_, err := GetShowTorrents("")
	if err == nil {
		t.Fatalf("should raise an error with missing argument")
	}
	if err != ErrMissingArgument {
		t.Errorf("should raise an error ErrMissingArgument")
	}
}

// TestGetShowTorrents tests GetShowTorrents
func TestGetShowTorrents(t *testing.T) {
	// Fake server with a fake answer
	var requestURI string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI = r.RequestURI
		fmt.Fprintln(w, rawResponseGetTorrents)
	}))
	defer ts.Close()
	endpoint = ts.URL

	s, err := GetShowTorrents("tt0383795")
	if err != nil {
		t.Fatalf("error while getting show torrents %+v", err)
	}

	expectedRequestURI := "/api/get-torrents?imdb_id=0383795&limit=100&page=1"
	if requestURI != expectedRequestURI {
		t.Fatalf("expected URL %q, got %q", expectedRequestURI, requestURI)
	}

	if reflect.DeepEqual(s, expectedTorrents) == false {
		t.Errorf("invalid torrent, expecting:\n %+v\nand got:\n%+v", expectedTorrents[0], s[0])
	}
}

// TestGetEpisodeTorrents tests GetEpisodeTorrents
func TestGetEpisodeTorrents(t *testing.T) {
	// Fake server with a fake answer
	var requestURI string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI = r.RequestURI
		fmt.Fprintln(w, rawResponseGetTorrents)
	}))
	defer ts.Close()
	endpoint = ts.URL

	e, err := GetEpisodeTorrents("tt2085059", 1, 1)
	if err == nil {
		t.Errorf("no errors while getting unknown show details %+v", err)
	}
	if e != nil {
		t.Error("episode should be nil", e)
	}
	if err != ErrEpisodeNotFound {
		t.Error("error should be ErrEpisodeNotFound", err)
	}

	e, err = GetEpisodeTorrents("tt0383795", 1, 10)
	if err != nil {
		t.Errorf("error while getting episode torrents %+v", err)
	}

	expectedRequestURI := "/api/get-torrents?imdb_id=0383795&limit=100&page=1"
	if requestURI != expectedRequestURI {
		t.Fatalf("expected URL %q, got %q", expectedRequestURI, requestURI)
	}

	if reflect.DeepEqual(e, expectedTorrents) == false {
		t.Errorf("invalid torrent, expecting:\n%+v\nand got:\n%+v", expectedTorrents, e)
	}
}

// TestGetTorrents tests GetTorrents
func TestGetTorrents(t *testing.T) {
	// Fake server with a fake answer
	var requestURI string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI = r.RequestURI
		fmt.Fprintln(w, rawResponseGetTorrents)
	}))
	defer ts.Close()
	endpoint = ts.URL

	e, err := GetTorrents(10, 1)
	if err != nil {
		t.Errorf("error while getting torrents %+v", err)
	}

	expectedRequestURI := "/api/get-torrents?limit=10&page=1"
	if requestURI != expectedRequestURI {
		t.Fatalf("expected URL %q, got %q", expectedRequestURI, requestURI)
	}

	if reflect.DeepEqual(e, expectedTorrents) == false {
		t.Errorf("invalid torrent, expecting:\n%+v\nand got:\n%+v", expectedTorrents, e)
	}
}

var rawResponseGetTorrents = `{ "torrents_count": 280450, "limit": 30, "page": 1, "imdb_id": "0383795", "torrents": [{ "id": 1443504, "hash": "033712dc2a9c6b23edbcbea737ee6ecf02619c2e", "filename": "The.Joy.of.Painting.S01E10.INTERNAL.480p.x264-mSD[eztv].mkv", "episode_url": "https://eztv.io/ep/1443504/the-joy-of-painting-s01e10-internal-480p-x264-msd/", "torrent_url": "https://zoink.ch/torrent/The.Joy.of.Painting.S01E10.INTERNAL.480p.x264-mSD[eztv].mkv.torrent", "magnet_url": "magnet:?xt=urn:btih:033712dc2a9c6b23edbcbea737ee6ecf02619c2e&dn=The.Joy.of.Painting.S01E10.INTERNAL.480p.x264-mSD%5Beztv%5D&tr=udp://tracker.coppersurfer.tk:80&tr=udp://glotorrents.pw:6969/announce&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://tracker.opentrackr.org:1337/announce&tr=udp://exodus.desync.com:6969", "title": "The Joy of Painting S01E10 INTERNAL 480p x264-mSD EZTV", "imdb_id": "0383795", "season": "1", "episode": "10", "small_screenshot": "//ezimg.ch/thumbs/the-joy-of-painting-s01e10-internal-480p-x264-msd-small.jpg", "large_screenshot": "//ezimg.ch/thumbs/the-joy-of-painting-s01e10-internal-480p-x264-msd-large.jpg", "seeds": 64, "peers": 0,"date_released_unix": 1588784249, "size_bytes": "86643580"}]}`

var expectedTorrents = []*EpisodeTorrent{
	&EpisodeTorrent{
		ID:              1443504,
		Hash:            "033712dc2a9c6b23edbcbea737ee6ecf02619c2e",
		Filename:        "The.Joy.of.Painting.S01E10.INTERNAL.480p.x264-mSD[eztv].mkv",
		EpisodeURL:      "https://eztv.io/ep/1443504/the-joy-of-painting-s01e10-internal-480p-x264-msd/",
		TorrentURL:      "https://zoink.ch/torrent/The.Joy.of.Painting.S01E10.INTERNAL.480p.x264-mSD[eztv].mkv.torrent",
		MagnetURL:       "magnet:?xt=urn:btih:033712dc2a9c6b23edbcbea737ee6ecf02619c2e&dn=The.Joy.of.Painting.S01E10.INTERNAL.480p.x264-mSD%5Beztv%5D&tr=udp://tracker.coppersurfer.tk:80&tr=udp://glotorrents.pw:6969/announce&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://tracker.opentrackr.org:1337/announce&tr=udp://exodus.desync.com:6969",
		Title:           "The Joy of Painting S01E10 INTERNAL 480p x264-mSD EZTV",
		ImdbID:          "tt0383795",
		Season:          1,
		Episode:         10,
		SmallScreenshot: "https://ezimg.ch/thumbs/the-joy-of-painting-s01e10-internal-480p-x264-msd-small.jpg",
		LargeScreenshot: "https://ezimg.ch/thumbs/the-joy-of-painting-s01e10-internal-480p-x264-msd-large.jpg",
		Seeds:           64,
		Peers:           0,
		DateReleased:    time.Unix(1588784249, 0),
		Size:            86643580,
	},
}
