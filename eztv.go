package eztv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var endpoint = "http://eztvapi.ml"

// Eztv errors
var (
	ErrEpisodeNotFound = errors.New("episode not found")
	ErrShowNotFound    = errors.New("show not found")
	ErrEmptyResponse   = errors.New("empty response from server")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrMissingArgument = errors.New("missing argument")
)

// Status reprensents the status response of a call
type Status struct {
	Status     string `json:"status"`
	Uptime     int64  `json:"uptime"`
	Server     string `json:"server"`
	Updated    int64  `json:"updated"`
	TotalShows int64  `json:"totalShows"`
	Version    string `json:"version"`
}

// Show represents a show
type Show struct {
	AirDay        string         `json:"air_day"`
	AirTime       string         `json:"air_time"`
	Country       string         `json:"country"`
	Episodes      []*ShowEpisode `json:"episodes"`
	Genres        []string       `json:"genres"`
	ID            string         `json:"_id"`
	Images        ShowImages     `json:"images"`
	ImdbID        string         `json:"imdb_id"`
	LastUpdated   int64          `json:"last_updated"`
	Network       string         `json:"network"`
	NumberSeasons int            `json:"num_seasons"`
	Rating        ShowRating     `json:"rating"`
	Runtime       string         `json:"runtime"`
	Slug          string         `json:"slug"`
	Status        string         `json:"status"`
	Synopsis      string         `json:"synopsis"`
	Title         string         `json:"title"`
	TvdbID        string         `json:"tvdb_id"`
	Year          string         `json:"year"`
}

// ShowEpisode represents a show episode
type ShowEpisode struct {
	Episode    int                     `json:"episode"`
	FirstAired int                     `json:"first_aired"`
	Overview   string                  `json:"overview"`
	Season     int                     `json:"season"`
	Title      string                  `json:"title"`
	TvdbID     int                     `json:"tvdb_id"`
	Torrents   map[string]*ShowTorrent `json:"torrents"`
}

// ShowImages represents the show images
type ShowImages struct {
	Poster string `json:"poster"`
	Fanart string `json:"fanart"`
	Banner string `json:"banner"`
}

// ShowRating respresents the show ratings
type ShowRating struct {
	Percentage float64 `json:"percentage"`
	Votes      int     `json:"votes"`
	Loved      int     `json:"loved"`
	Hated      int     `json:"hated"`
}

// ShowSeason represents a show season
type ShowSeason struct {
	Episodes map[string]*ShowEpisode
}

// ShowTorrent represents the show torrent
type ShowTorrent struct {
	Peers int    `json:"peers"`
	Seeds int    `json:"seeds"`
	URL   string `json:"url"`
}

// Ping will make a simple GET on the root path of the api to see if it's alive
func Ping() (*Status, error) {
	// Generate URL
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	// Make the request
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	status := &Status{}
	err = json.Unmarshal(body, status)
	if err != nil {
		return nil, err
	}
	return status, nil
}

// GetShowDetails will get the infos about a show from an ImdbID
func GetShowDetails(ImdbID string) (*Show, error) {
	if ImdbID == "" {
		return nil, ErrMissingArgument
	}
	// Generate URL
	u, err := url.Parse(fmt.Sprintf("%s/show/%s", endpoint, ImdbID))
	if err != nil {
		return nil, err
	}

	// Make the request
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Test for an empty body, a bit hacky but the API returns an empty
	// response sometimes...
	// 2 is for "{}" which is the smallest JSON response
	if len(body) < 2 {
		return nil, ErrEmptyResponse
	}

	s := &Show{}
	err = json.Unmarshal(body, s)
	if err != nil {
		return nil, err
	}

	// No title is considered as a failure
	if s.Title == "" {
		return nil, ErrShowNotFound
	}

	return s, nil
}

// GetEpisode will get the infos about a show's episode from an ImdbID,
// a season number and an episode number
func GetEpisode(ImdbID string, season, episode int) (*ShowEpisode, error) {
	s, err := GetShowDetails(ImdbID)
	if err != nil {
		return nil, err
	}
	for _, e := range s.Episodes {
		if e.Episode == episode && e.Season == season {
			return e, nil
		}
	}

	return nil, ErrEpisodeNotFound
}

// GetSeason will get the infos about all the show's episode from an
// ImdbID and a season number
func GetSeason(ImdbID string, season int) ([]*ShowEpisode, error) {
	s, err := GetShowDetails(ImdbID)
	if err != nil {
		return nil, err
	}
	episodes := []*ShowEpisode{}
	for _, e := range s.Episodes {
		if e.Season == season {
			episodes = append(episodes, e)
		}
	}

	return episodes, nil
}

// ListShows returns a list of 50 shows with metadata
func ListShows(page int) ([]*Show, error) {
	return getShows("", page)
}

// SearchShow will get the infos about a show from a string
func SearchShow(keyword string) ([]*Show, error) {
	if keyword == "" {
		return nil, ErrMissingArgument
	}
	return getShows(keyword, 1)
}

func getShows(keyword string, page int) ([]*Show, error) {
	if page <= 0 {
		return nil, ErrInvalidArgument
	}
	// Generate URL
	u, err := url.Parse(fmt.Sprintf("%s/shows/%d", endpoint, page))
	if err != nil {
		return nil, err
	}
	urlValues := &url.Values{}
	if keyword != "" {
		urlValues.Add("keywords", keyword)
	}
	u.RawQuery = urlValues.Encode()

	// Make the request
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s := []*Show{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
