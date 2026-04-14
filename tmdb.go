package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const tmdbBase = "https://api.themoviedb.org/3"

type TMDBResult struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`         // movies
	Name          string `json:"name"`          // TV
	ReleaseDate   string `json:"release_date"`  // movies
	FirstAirDate  string `json:"first_air_date"` // TV
	MediaType     string `json:"media_type"`
	Overview      string `json:"overview"`
}

type TMDBExternalIDs struct {
	ImdbID string `json:"imdb_id"`
}

type tmdbSearchResponse struct {
	Results []TMDBResult `json:"results"`
}

func (r *TMDBResult) DisplayTitle() string {
	if r.Title != "" {
		return r.Title
	}
	return r.Name
}

func (r *TMDBResult) Year() string {
	date := r.ReleaseDate
	if date == "" {
		date = r.FirstAirDate
	}
	if len(date) >= 4 {
		return date[:4]
	}
	return "????"
}

func (r *TMDBResult) IsTV() bool {
	return r.MediaType == "tv"
}

func tmdbGet(path, apiKey string, params url.Values) ([]byte, error) {
	params.Set("api_key", apiKey)
	u := fmt.Sprintf("%s%s?%s", tmdbBase, path, params.Encode())
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("TMDB request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("invalid TMDB API key — visit https://www.themoviedb.org/settings/api")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("TMDB error %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func searchTMDB(query, apiKey string) ([]TMDBResult, error) {
	params := url.Values{
		"query":         {query},
		"include_adult": {"false"},
		"page":          {"1"},
	}
	body, err := tmdbGet("/search/multi", apiKey, params)
	if err != nil {
		return nil, err
	}
	var resp tmdbSearchResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	// Filter to movies and TV only
	var results []TMDBResult
	for _, r := range resp.Results {
		if r.MediaType == "movie" || r.MediaType == "tv" {
			results = append(results, r)
		}
	}
	return results, nil
}

func getExternalIDs(id int, mediaType, apiKey string) (*TMDBExternalIDs, error) {
	path := fmt.Sprintf("/%s/%d/external_ids", mediaType, id)
	body, err := tmdbGet(path, apiKey, url.Values{})
	if err != nil {
		return nil, err
	}
	var ids TMDBExternalIDs
	if err := json.Unmarshal(body, &ids); err != nil {
		return nil, err
	}
	return &ids, nil
}

// formatMediaEntry produces a display string for a TMDB result.
func formatMediaEntry(r TMDBResult, imdbID string) string {
	kind := "Movie"
	if r.IsTV() {
		kind = "TV Show"
	}
	tmdbStr := fmt.Sprintf("TMDB: %s", cyanText(fmt.Sprintf("%d", r.ID)))
	imdbStr := ""
	if imdbID != "" {
		imdbStr = "  IMDB: " + cyanText(imdbID)
	}
	title := boldText(r.DisplayTitle())
	year := grayText("(" + r.Year() + ")")
	badge := yellowText("[" + kind + "]")
	return fmt.Sprintf("%s %s  %s  %s%s", title, year, badge, tmdbStr, imdbStr)
}

// safeFilename returns a filesystem-safe version of a string.
func safeFilename(s string) string {
	replacer := strings.NewReplacer(
		"/", "-", "\\", "-", ":", "-", "*", "-",
		"?", "", "\"", "", "<", "", ">", "", "|", "-",
	)
	return strings.TrimSpace(replacer.Replace(s))
}
