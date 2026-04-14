package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const wyzieBase = "https://sub.wyzie.io"

type Subtitle struct {
	ID                string `json:"id"`
	URL               string `json:"url"`
	Format            string `json:"format"`
	Encoding          string `json:"encoding"`
	Display           string `json:"display"`
	Language          string `json:"language"`
	Media             string `json:"media"`
	IsHearingImpaired bool   `json:"isHearingImpaired"`
	Source            string `json:"source"`
	Release           string `json:"release"`
	FileName          string `json:"fileName"`
	DownloadCount     *int   `json:"downloadCount"`
	Origin            string `json:"origin"`
}

type SubtitleSearchParams struct {
	TMDBID   int
	Season   int
	Episode  int
	Language []string
	Source   string
	Hi       bool
	WyzieKey string
}

func searchSubtitles(p SubtitleSearchParams) ([]Subtitle, error) {
	params := url.Values{}
	params.Set("id", strconv.Itoa(p.TMDBID))
	params.Set("key", p.WyzieKey)

	if p.Season > 0 && p.Episode > 0 {
		params.Set("season", strconv.Itoa(p.Season))
		params.Set("episode", strconv.Itoa(p.Episode))
	}
	if len(p.Language) > 0 {
		params.Set("language", strings.Join(p.Language, ","))
	}
	if p.Source == "all" {
		params.Set("source", "all")
	} else if p.Source != "" {
		params.Set("source", p.Source)
	}
	if p.Hi {
		params.Set("hi", "true")
	}

	u := fmt.Sprintf("%s/search?%s", wyzieBase, params.Encode())
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("wyzie request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("invalid Wyzie API key — visit https://sub.wyzie.io/redeem")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("wyzie error %d: %s", resp.StatusCode, string(body))
	}

	var subs []Subtitle
	if err := json.Unmarshal(body, &subs); err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}
	return subs, nil
}

func formatSubtitleEntry(s Subtitle) string {
	lang := boldText(s.Display)
	if s.IsHearingImpaired {
		lang += yellowText(" [SDH]")
	}
	format := cyanText(strings.ToUpper(s.Format))
	source := grayText(s.Source)
	release := s.Release
	if release == "" {
		release = grayText("unknown release")
	}
	dl := ""
	if s.DownloadCount != nil {
		dl = grayText(fmt.Sprintf("  ↓%d", *s.DownloadCount))
	}
	origin := ""
	if s.Origin != "" {
		origin = magentaText(" [" + s.Origin + "]")
	}
	return fmt.Sprintf("%s  %s  %s  %s%s%s", lang, format, source, release, origin, dl)
}

func getSources(wyzieKey string) ([]string, error) {
	u := fmt.Sprintf("%s/sources?key=%s", wyzieBase, url.QueryEscape(wyzieKey))
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	// Response shape: { "opensubtitles": true, "subdl": false, ... }
	var raw map[string]bool
	if err := json.Unmarshal(body, &raw); err != nil {
		return []string{"opensubtitles", "subdl", "podnapisi"}, nil
	}
	var sources []string
	for name, enabled := range raw {
		if enabled {
			sources = append(sources, name)
		}
	}
	return sources, nil
}
