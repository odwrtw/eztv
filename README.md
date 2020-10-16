EZTV API client
=========

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/odwrtw/eztv)
[![Build Status](https://travis-ci.org/odwrtw/eztv.svg?branch=master)](https://travis-ci.org/odwrtw/eztv)
[![Coverage Status](https://coveralls.io/repos/odwrtw/eztv/badge.svg?branch=master&service=github)](https://coveralls.io/github/odwrtw/eztv?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/odwrtw/eztv)](https://goreportcard.com/report/github.com/odwrtw/eztv)

This is a wrapper around the EZTV API written in go

## Get details (episodes and torrents) from an imdb ID

```
        // Get all the episodes of tt2085059
        show, err := eztv.GetShowTorrents("tt2149175")

        // Get the torrents of the first episode of the second season of tt2085059
        e, err := GetEpisodeTorrents("tt2085059", 2, 1)
```

## List Shows

```
        // List popular shows (with pagination)
        list, err := eztv.GetTorrents(20, 2)
```
