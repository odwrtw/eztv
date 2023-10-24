EZTV API client
=========

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/odwrtw/eztv)
[![Coverage Status](https://coveralls.io/repos/odwrtw/eztv/badge.svg?branch=master&service=github)](https://coveralls.io/github/odwrtw/eztv?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/odwrtw/eztv)](https://goreportcard.com/report/github.com/odwrtw/eztv)

This is a wrapper around the EZTV API written in go.

## Usage

### Get all the torrents for a show from its imdb id

```go
torrents, err := eztv.GetShowTorrents("tt2149175")
```

### Get all the torrents of a specific episode

```go
// Get the torrents of the first episode of the second season of tt2085059
torrents, err := GetEpisodeTorrents("tt2085059", 2, 1)
```

### Get the last torrents available

```go
// Get the 20 torrents from the 2 first pages
torrents, err := eztv.GetTorrents(20, 2)
```
