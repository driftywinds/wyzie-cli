package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	initColor()
	printBanner()

	// ── Load / setup config ──────────────────────────────────────────
	cfg, err := loadConfig()
	if err != nil {
		printError("Config error: " + err.Error())
		os.Exit(1)
	}
	if err := ensureConfig(cfg); err != nil {
		printError(err.Error())
		os.Exit(1)
	}

	// ── Step 1: Search for media ─────────────────────────────────────
	separator()
	fmt.Println(boldText("  Search for a Movie or TV Show"))
	fmt.Println()
	query, err := prompt("  Title", "")
	if err != nil || strings.TrimSpace(query) == "" {
		printError("No title entered.")
		os.Exit(1)
	}

	fmt.Println()
	printInfo("Searching TMDB...")

	results, err := searchTMDB(query, cfg.TMDBKey)
	if err != nil {
		printError("TMDB search failed: " + err.Error())
		os.Exit(1)
	}
	if len(results) == 0 {
		printWarn("No results found for: " + query)
		os.Exit(0)
	}

	// Cap at 8 results and fetch their IMDB IDs
	if len(results) > 8 {
		results = results[:8]
	}
	imdbIDs := make([]string, len(results))
	for i, r := range results {
		mediaType := "movie"
		if r.IsTV() {
			mediaType = "tv"
		}
		if ext, err := getExternalIDs(r.ID, mediaType, cfg.TMDBKey); err == nil {
			imdbIDs[i] = ext.ImdbID
		}
	}

	// ── Step 2: Pick media ───────────────────────────────────────────
	entries := make([]string, len(results))
	for i, r := range results {
		entries[i] = formatMediaEntry(r, imdbIDs[i])
	}

	separator()
	mediaIdx, err := pickOne("Select Media", entries)
	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}
	chosen := results[mediaIdx]
	chosenIMDB := imdbIDs[mediaIdx]

	separator()
	fmt.Printf("  %s  %s\n",
		boldText(chosen.DisplayTitle()),
		grayText(fmt.Sprintf("(%s)  TMDB:%d  IMDB:%s", chosen.Year(), chosen.ID, chosenIMDB)),
	)

	// Warn when no IMDB ID — Wyzie may not recognise this entry
	if chosenIMDB == "" {
		fmt.Println()
		printWarn("This entry has no IMDB ID. It may be too new or incomplete — subtitles might not be found.")
		fmt.Println(grayText("  Tip: try a different result that has an IMDB ID listed."))
	}

	// ── Step 3: Season / Episode for TV shows ────────────────────────
	var season, episode int
	if chosen.IsTV() {
		fmt.Println()
		seasonStr, _ := prompt("  Season number", "1")
		episodeStr, _ := prompt("  Episode number", "1")
		season, _ = strconv.Atoi(strings.TrimSpace(seasonStr))
		episode, _ = strconv.Atoi(strings.TrimSpace(episodeStr))
	}

	// ── Step 4: Language ─────────────────────────────────────────────
	separator()
	fmt.Println(boldText("  Subtitle Options"))
	fmt.Println()
	fmt.Println(grayText("  Language codes: ISO 639-1, comma-separated (e.g. en,es,fr)"))
	langInput, _ := prompt("  Language(s)", "en")
	langs := parseCsvInput(langInput)

	hiInput, _ := prompt("  Include hearing-impaired / SDH subtitles? [y/n]", "y")
	hi := strings.ToLower(strings.TrimSpace(hiInput)) == "y"

	// ── Step 5: Source ───────────────────────────────────────────────
	fmt.Println()
	printInfo("Fetching available sources...")
	sources, _ := getSources(cfg.WyzieKey)

	sourceItems := []string{boldText("all") + grayText("  (every enabled source)")}
	sourceLabels := []string{"all"}
	for _, s := range sources {
		sourceItems = append(sourceItems, s)
		sourceLabels = append(sourceLabels, s)
	}

	sourceIdx, err := pickOne("Subtitle Source", sourceItems)
	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}
	selectedSource := sourceLabels[sourceIdx]

	// ── Step 6: Fetch subtitles ──────────────────────────────────────
	separator()
	printInfo("Searching for subtitles...")

	subs, err := searchSubtitles(SubtitleSearchParams{
		TMDBID:   chosen.ID,
		IMDBID:   chosenIMDB,
		Season:   season,
		Episode:  episode,
		Language: langs,
		Source:   selectedSource,
		Hi:       hi,
		WyzieKey: cfg.WyzieKey,
	})
	if err != nil {
		printError("Subtitle search failed: " + err.Error())
		os.Exit(1)
	}
	if len(subs) == 0 {
		printWarn("No subtitles found with those filters. Try broader language/source settings.")
		os.Exit(0)
	}
	printSuccess(fmt.Sprintf("Found %d subtitle(s)", len(subs)))

	// ── Step 7: Pick subtitle ────────────────────────────────────────
	subEntries := make([]string, len(subs))
	for i, s := range subs {
		subEntries[i] = formatSubtitleEntry(s)
	}

	subIdx, err := pickOne("Select Subtitle", subEntries)
	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}
	picked := subs[subIdx]

	// ── Step 8: Save ─────────────────────────────────────────────────
	separator()
	defaultName := suggestFilename(
		chosen.DisplayTitle(), chosen.Year(),
		picked.Language, picked.Format, picked.IsHearingImpaired,
	)
	outFile, _ := prompt("  Save as", defaultName)
	if strings.TrimSpace(outFile) == "" {
		outFile = defaultName
	}

	fmt.Println()
	printInfo("Downloading...")
	if err := downloadSubtitle(picked, outFile); err != nil {
		printError("Download failed: " + err.Error())
		os.Exit(1)
	}

	separator()
	printSuccess("Saved: " + boldText(outFile))
	fmt.Println()
}

func parseCsvInput(s string) []string {
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
