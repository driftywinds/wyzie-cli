# wyzie-subs

A cross-platform CLI tool for searching and downloading subtitles via [Wyzie Subs](https://sub.wyzie.io).

## Features

- Search movies and TV shows by title
- Displays TMDB and IMDB IDs for each result
- Filter subtitles by language, source, and SDH/hearing-impaired
- Defaults to English + SDH, all sources
- Downloads subtitle to the current directory
- Config stored securely in `~/.wyziesubs/config.json`
- Zero runtime dependencies — single binary per platform

## Prerequisites

Two free API keys are required (no account, no credit card needed):

| Key | Where to get it |
|-----|----------------|
| **Wyzie API key** | https://sub.wyzie.io/redeem |
| **TMDB API key** | https://www.themoviedb.org/settings/api |

Both are free and take under a minute to obtain. Keys are saved locally on first run.

## Quick start

```bash
# Run directly (Go installed)
go run .

# Or build for your platform
go build -o wyzie-subs .
./wyzie-subs
```

## Build all platform binaries

```bash
make build
```

Produces binaries in `./dist/`:

```
dist/
├── wyzie-subs-linux-amd64
├── wyzie-subs-linux-arm64
├── wyzie-subs-macos-amd64
├── wyzie-subs-macos-arm64
└── wyzie-subs-windows-amd64.exe
```

No installation needed — just copy the right binary to any machine and run it.

## Usage walkthrough

```
  Title: The Martian

→ Searching TMDB...

    1.  The Martian (2015)  [Movie]  TMDB: 286217  IMDB: tt3659388
    2.  The Martian Chronicles (1980)  [TV Show]  TMDB: 26491  IMDB: tt0079822

  Enter number (1): 1

  Language(s) (en): en
  Include hearing-impaired / SDH subtitles? [y/n] (y):

→ Fetching available sources...

    1.  all  (every enabled source)
    2.  opensubtitles
    3.  subdl
    4.  podnapisi

  Enter number (1): 1

→ Searching for subtitles...
✓ Found 12 subtitle(s)

    1.  English  SRT  opensubtitles  The.Martian.2015.1080p.BluRay  [BluRay]  ↓5234
    2.  English [SDH]  SRT  subdl  The.Martian.2015.WEB-DL  [WEB]  ↓3211
    ...

  Enter number (1): 1

  Save as (The.Martian.2015.en.srt):

→ Downloading...
✓ Saved: The.Martian.2015.en.srt
```

## Config

Stored at `~/.wyziesubs/config.json`. Delete this file to re-enter your API keys.

```json
{
  "wyzie_key": "wyzie-abc123xyz",
  "tmdb_key": "your-tmdb-key-here"
}
```
