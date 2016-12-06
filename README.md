# Manifest

Manifest is a library for parsing and creating video manifest files.

## Usage

See [godoc](https://godoc.org/github.com/ingest/manifest).

### Features

* Complete HLS compliance upto version 7, defined in the _April 4 2016_ [specification](https://tools.ietf.org/html/draft-pantos-http-live-streaming-19)

### In-progress

* DASH is currently not running in production, follow along and help us guide the creation!

## Motivation

Ingest as a organization makes use of lots of open-source software. We initially worked and used [grafov/m3u8](https://github.com/grafov/m3u8) for parsing HLS media playlists but found that we quickly were outpacing the scope of the library. As our roadmap grew and required future HLS version support, and other manifest formats such as DASH, we felt that it would be best to move the development in-house.