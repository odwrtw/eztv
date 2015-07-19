EZTV API client
=========

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/odwrtw/eztv)
[![Build Status](https://travis-ci.org/odwrtw/eztv.svg?branch=master)](https://travis-ci.org/odwrtw/eztv)
[![Coverage Status](https://coveralls.io/repos/odwrtw/eztv/badge.svg?branch=master&service=github)](https://coveralls.io/github/odwrtw/eztv?branch=master)

This is a wrapper around the EZTV API written in go

## Search Shows

```
  // Search Shows
	list, err := eztv.SearchShow("black mirror")
```

## Get details (episodes and torrents) from an imdb ID

```
  // Get all the episodes of tt2085059
	show, err := eztv.GetShowDetails("tt2149175")

  // Get the first episode of the second season of tt2085059
	e, err := GetEpisode("tt2085059", 2, 1)

  // Get all the episodes of the second season of tt2085059
	showList, err := GetSeason("tt2085059", 2)
```

## List Shows

```
  // List popular shows (with pagination)
	list, err := eztv.ListShows(1)
```

## Ping the API to see if it's up

```
	// Test if the API is up
	status, err := eztv.Ping()
```
