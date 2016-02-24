package eztv

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// TestPing tests the Ping method
func TestPing(t *testing.T) {
	rawHTMLResponse := `{"status":"online","uptime":2381142,"server":"serv01","totalShows":1426,"updated":1454655986,"version":"1.0.3"}`

	// Fake server with a fake answer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rawHTMLResponse)
	}))
	defer ts.Close()
	endpoint = ts.URL

	status, err := Ping()
	if err != nil {
		t.Errorf("Error when pinging the endpoint")
	}

	// Expected result
	expectedResult := &Status{
		Status:     "online",
		Uptime:     2381142,
		Server:     "serv01",
		TotalShows: 1426,
		Updated:    1454655986,
		Version:    "1.0.3",
	}

	if reflect.DeepEqual(status, expectedResult) == false {
		t.Errorf("Response not properly set")
	}
}

// TestShowNotFoundEmptyResponse tests the ErrShowNotFound error
func TestShowNotFoundEmptyResponse(t *testing.T) {
	// Fake server with a fake answer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
	}))
	defer ts.Close()
	endpoint = ts.URL

	_, err := GetShowDetails("fakeID")
	if err == nil {
		t.Fatalf("Got no error, expected %q", ErrEmptyResponse)
	}

	if err != ErrEmptyResponse {
		t.Fatalf("Expected %q, got %q", ErrEmptyResponse, err)
	}
}

// TestShowNotFound tests the ErrShowNotFound error
func TestShowNotFound(t *testing.T) {
	// Fake server with a fake answer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{}")
	}))
	defer ts.Close()
	endpoint = ts.URL

	_, err := GetShowDetails("fakeID")
	if err == nil {
		t.Fatalf("Got no error, expected %q", ErrShowNotFound)
	}

	if err != ErrShowNotFound {
		t.Fatalf("Expected %q, got %q", ErrShowNotFound, err)
	}
}

// TestBadGetShowDetails tests GetShowDetails with missing arguments
func TestBadGetShowDetails(t *testing.T) {
	_, err := GetShowDetails("")
	if err == nil {
		t.Errorf("Should raise an error with missing argument")
	}
	if err != ErrMissingArgument {
		t.Errorf("Should raise an error ErrMissingArgument")
	}
}

// TestGetShowDetails tests GetShowDetails with missing arguments
func TestGetShowDetails(t *testing.T) {
	// Fake server with a fake answer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rawHTMLResponseGetShowDetails)
	}))
	defer ts.Close()
	endpoint = ts.URL

	s, err := GetShowDetails("tt2085059")
	if err != nil {
		t.Errorf("Error while getting show details %+v", err)
	}

	expectedResult := &Show{
		AirDay:        "Tuesday",
		AirTime:       "21:00",
		Country:       "gb",
		ID:            "tt2085059",
		ImdbID:        "tt2085059",
		LastUpdated:   1430654860060,
		Network:       "Channel 4",
		NumberSeasons: 2,
		Runtime:       "60",
		Slug:          "black-mirror",
		Status:        "returning series",
		Synopsis:      "Over the last ten years, ...",
		Title:         "Black Mirror",
		TvdbID:        "253463",
		Year:          "2011",
		Genres:        []string{"Drama"},
		Rating: ShowRating{
			Percentage: 87.1234,
			Votes:      2662,
			Loved:      0,
			Hated:      0,
		},
		Images: ShowImages{
			Poster: "https://walter.trakt.us/images/shows/000/041/793/posters/original/bc21969723.jpg",
			Fanart: "https://walter.trakt.us/images/shows/000/041/793/fanarts/original/546a01eb67.jpg",
			Banner: "placeholders/banner.png",
		},
		Episodes: []*ShowEpisode{
			expectedEpisode,
		},
	}
	if reflect.DeepEqual(s, expectedResult) == false {
		t.Errorf("Response not properly set\n Expecting %+v\nAnd got %+v", expectedResult, s)
	}
}

func TestGetEpisode(t *testing.T) {
	// Fake server with a fake answer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rawHTMLResponseGetShowDetails)
	}))
	defer ts.Close()
	endpoint = ts.URL

	e, err := GetEpisode("tt2085059", 1, 1)
	if err == nil {
		t.Errorf("No errors while getting unknown show details %+v", err)
	}
	if e != nil {
		t.Error("Episode should be nil", e)
	}
	if err != ErrEpisodeNotFound {
		t.Error("Error should Be ErrEpisodeNotFound", err)
	}

	e, err = GetEpisode("tt2085059", 2, 1)
	if err != nil {
		t.Errorf("Error while getting show details %+v", err)
	}
	if reflect.DeepEqual(e, expectedEpisode) == false {
		t.Errorf("Response not properly set\n Expecting %+v\nAnd got %+v", expectedEpisode, e)
	}
}

// TestGetSeason will test the function to get all the episodes of a season
func TestGetSeason(t *testing.T) {
	// Fake server with a fake answer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rawHTMLResponseGetShowDetails)
	}))
	defer ts.Close()
	endpoint = ts.URL

	// With a non-existing season, we should have an empty list
	emptyList, err := GetSeason("tt2085059", 5)
	if err != nil {
		t.Errorf("Error while getting show details %+v", err)
	}
	expectedEmptyResult := []*ShowEpisode{}
	if reflect.DeepEqual(emptyList, expectedEmptyResult) == false {
		t.Errorf("Response not properly set\n Expecting %+v\nAnd got %+v", expectedEmptyResult, emptyList)
	}

	// With an existing season, we should have a list of episodes
	e, err := GetSeason("tt2085059", 2)
	if err != nil {
		t.Errorf("Error while getting show details %+v", err)
	}
	expectedResult := []*ShowEpisode{
		expectedEpisode,
	}
	if reflect.DeepEqual(e, expectedResult) == false {
		t.Errorf("Response not properly set\n Expecting %+v\nAnd got %+v", expectedResult, e)
	}
}

// TestSearchShow will test the search of a show
func TestSearchShow(t *testing.T) {
	// Fake server with a fake answer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rawHTMLResponseSearchShow)
	}))
	defer ts.Close()
	endpoint = ts.URL
	s, err := SearchShow("black mirror")
	if err != nil {
		t.Errorf("Error while searching show %+v", err)
	}
	expectedResult := []*Show{
		&Show{
			ID:            "tt2085059",
			ImdbID:        "tt2085059",
			LastUpdated:   1430654860060,
			NumberSeasons: 2,
			Slug:          "black-mirror",
			Title:         "Black Mirror",
			TvdbID:        "253463",
			Year:          "2011",
			Rating: ShowRating{
				Percentage: 87,
				Votes:      2662,
				Loved:      0,
				Hated:      0,
			},
			Images: ShowImages{
				Poster: "https://walter.trakt.us/images/shows/000/041/793/posters/original/bc21969723.jpg",
				Fanart: "https://walter.trakt.us/images/shows/000/041/793/fanarts/original/546a01eb67.jpg",
				Banner: "placeholders/banner.png",
			},
		},
	}
	if reflect.DeepEqual(s, expectedResult) == false {
		t.Errorf("Response not properly set\n Expecting %+v\nAnd got %+v", expectedResult, s)
	}
}

var rawHTMLResponseSearchShow = `[{"_id":"tt2085059","images":{"poster":"https://walter.trakt.us/images/shows/000/041/793/posters/original/bc21969723.jpg","fanart":"https://walter.trakt.us/images/shows/000/041/793/fanarts/original/546a01eb67.jpg","banner":"placeholders/banner.png"},"imdb_id":"tt2085059","last_updated":1430654860060,"num_seasons":2,"rating":{"percentage":87,"votes":2662,"loved":0,"hated":0},"slug":"black-mirror","title":"Black Mirror","tvdb_id":"253463","year":"2011"}]`

var rawHTMLResponseGetShowDetails = `{"_id": "tt2085059","air_day": "Tuesday","air_time": "21:00","country": "gb","imdb_id": "tt2085059","last_updated": 1430654860060,"network": "Channel 4","num_seasons": 2,"rating": {"hated": 0,"loved": 0,"percentage": 87.1234,"votes": 2662},"runtime": "60","slug": "black-mirror","status": "returning series","synopsis": "Over the last ten years, ...","title": "Black Mirror","tvdb_id": "253463","year": "2011","episodes": [{"date_based": false,"episode": 1,"first_aired": 1360616400,"overview": "Martha and Ash are a young couple ....","season": 2,"title": "Be Right Back","torrents": {"480p": {"peers": 0,"seeds": 0,"url": "magnet:?xt=urn:btih:Z4I3PBZVO&dn=Black.Mirror.2x01&tr=udp://tracker.openent.com:80/"},"720p": {"peers": 0,"seeds": 0,"url": "magnet:?xt=urn:btih:LSVN23DMU&dn=Black.Mirror.2x01tr=udp://tracker.openent.com:80/"}},"tvdb_id": 4418485,"watched": {"watched": false}}],"genres": ["Drama"],"images": {"banner": "placeholders/banner.png","fanart": "https://walter.trakt.us/images/shows/000/041/793/fanarts/original/546a01eb67.jpg","poster": "https://walter.trakt.us/images/shows/000/041/793/posters/original/bc21969723.jpg"}}`

var expectedEpisode = &ShowEpisode{
	Episode:    1,
	FirstAired: 1360616400,
	Overview:   "Martha and Ash are a young couple ....",
	Season:     2,
	Title:      "Be Right Back",
	TvdbID:     4418485,
	Torrents: map[string]*ShowTorrent{
		"480p": &ShowTorrent{
			Peers: 0,
			Seeds: 0,
			URL:   "magnet:?xt=urn:btih:Z4I3PBZVO&dn=Black.Mirror.2x01&tr=udp://tracker.openent.com:80/",
		},
		"720p": &ShowTorrent{
			Peers: 0,
			Seeds: 0,
			URL:   "magnet:?xt=urn:btih:LSVN23DMU&dn=Black.Mirror.2x01tr=udp://tracker.openent.com:80/",
		},
	},
}
