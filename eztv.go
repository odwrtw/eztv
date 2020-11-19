package eztv

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var endpoint = "https://eztv.re"

const (
	torrentEndpoint = "/api/get-torrents"
)

// Eztv errors
var (
	ErrEpisodeNotFound = errors.New("episode not found")
	ErrShowNotFound    = errors.New("show not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrMissingArgument = errors.New("missing argument")
)

// response reprensents the response structure of every call
type response struct {
	TorrentsCount int               `json:"torrents_count"`
	Limit         int               `json:"limit"`
	Page          int               `json:"page"`
	ImdbID        string            `json:"imdb_id"`
	Torrents      []*EpisodeTorrent `json:"torrents"`
}

// EpisodeTorrent represents a torrent for an episode
type EpisodeTorrent struct {
	ID              int       `json:"id"`
	Hash            string    `json:"hash"`
	Filename        string    `json:"filename"`
	EpisodeURL      string    `json:"episode_url"`
	TorrentURL      string    `json:"torrent_url"`
	MagnetURL       string    `json:"magnet_url"`
	Title           string    `json:"title"`
	ImdbID          string    `json:"imdb_id"`
	Season          int       `json:"season,string"`
	Episode         int       `json:"episode,string"`
	SmallScreenshot string    `json:"small_screenshot"`
	LargeScreenshot string    `json:"large_screenshot"`
	Seeds           int       `json:"seeds"`
	Peers           int       `json:"peers"`
	DateReleased    time.Time `json:"date_released_unix"`
	Size            uint      `json:"size_bytes,string"`
}

// UnmarshalJSON implements the Unmarshaller interface
func (t *EpisodeTorrent) UnmarshalJSON(data []byte) error {
	type Alias EpisodeTorrent
	// Unmarshal everything except the fields we need to edit
	aux := &struct {
		*Alias
		ImdbID          string `json:"imdb_id"`
		DateReleased    int64  `json:"date_released_unix"`
		SmallScreenshot string `json:"small_screenshot"`
		LargeScreenshot string `json:"large_screenshot"`
	}{}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	// Copy the whole thing
	*t = EpisodeTorrent(*(aux.Alias))
	// Set the proper DateReleased time / imdbID / screenshots URLs
	t.DateReleased = time.Unix(aux.DateReleased, 0)
	t.ImdbID = "tt" + aux.ImdbID
	t.SmallScreenshot = "https:" + aux.SmallScreenshot
	t.LargeScreenshot = "https:" + aux.LargeScreenshot

	return nil
}

func get(path string, v *url.Values) (*response, error) {
	// Make the request
	resp, err := http.Get(endpoint + path + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status %d when making the request", resp.StatusCode)
	}

	r := &response{}
	return r, json.NewDecoder(resp.Body).Decode(&r)
}

// getAll makes requests on all the pages until we have all the torrents
func getAll(path string, v *url.Values) ([]*EpisodeTorrent, error) {
	var maxPages = 20

	var torrents []*EpisodeTorrent
	for i := 1; i < maxPages; i++ {
		v.Add("page", strconv.Itoa(i))
		v.Add("limit", "100")
		resp, err := get(path, v)
		if err != nil {
			return nil, err
		}
		// No ImdbID or no torrents are considered as a failure
		if resp.ImdbID == "" || resp.TorrentsCount == 0 {
			return nil, ErrShowNotFound
		}

		torrents = append(torrents, resp.Torrents...)
		// We asked for a 100 torrents, so if we got less, then the search is
		// over
		if len(resp.Torrents) < 100 {
			break
		}
	}
	return torrents, nil
}

// GetShowTorrents will get the torrents of a show from an ImdbID
func GetShowTorrents(imdbID string) ([]*EpisodeTorrent, error) {
	if imdbID == "" {
		return nil, ErrMissingArgument
	}
	// The API only accepts imdbIDs without their "tt"
	imdbID = strings.Replace(imdbID, "t", "", -1)

	v := &url.Values{}
	v.Add("imdb_id", imdbID)

	// Make the request
	torrents, err := getAll(torrentEndpoint, v)
	if err != nil {
		return nil, err
	}

	return torrents, nil
}

// GetEpisodeTorrents will get the torrents of show's episode from an ImdbID,
// a season number and an episode number
func GetEpisodeTorrents(ImdbID string, season, episode int) ([]*EpisodeTorrent, error) {
	episodeTorrents := []*EpisodeTorrent{}
	torrents, err := GetShowTorrents(ImdbID)
	if err != nil {
		return nil, err
	}

	for _, e := range torrents {
		if e.Episode == episode && e.Season == season {
			episodeTorrents = append(episodeTorrents, e)
		}
	}

	if len(episodeTorrents) == 0 {
		return nil, ErrEpisodeNotFound
	}
	return episodeTorrents, nil
}

// GetTorrents will get torrents, given a limit and a page number
func GetTorrents(limit, page int) ([]*EpisodeTorrent, error) {
	// Generate URL
	v := &url.Values{}
	v.Add("limit", strconv.Itoa(limit))
	v.Add("page", strconv.Itoa(page))

	// Make the request
	r, err := get(torrentEndpoint, v)
	if err != nil {
		return nil, err
	}

	return r.Torrents, nil
}
