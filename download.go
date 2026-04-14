package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func downloadSubtitle(sub Subtitle, destPath string) error {
	resp, err := http.Get(sub.URL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("download error %d", resp.StatusCode)
	}
	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func suggestFilename(media, year, lang, format string, hi bool) string {
	name := safeFilename(media)
	if year != "" && year != "????" {
		name += "." + year
	}
	name += "." + lang
	if hi {
		name += ".sdh"
	}
	name += "." + strings.ToLower(format)
	return name
}
