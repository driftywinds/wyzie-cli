package main

import (
	"fmt"
	"os"
	"runtime"
)

// ANSI color codes
const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
	gray    = "\033[90m"
)

var colorEnabled bool

func initColor() {
	if runtime.GOOS == "windows" {
		colorEnabled = enableWindowsVT()
	} else {
		// Check if stdout is a terminal
		fi, err := os.Stdout.Stat()
		colorEnabled = err == nil && (fi.Mode()&os.ModeCharDevice) != 0
	}
}

func colorize(code, s string) string {
	if !colorEnabled {
		return s
	}
	return code + s + reset
}

func boldText(s string) string    { return colorize(bold, s) }
func dimText(s string) string     { return colorize(dim, s) }
func redText(s string) string     { return colorize(red, s) }
func greenText(s string) string   { return colorize(green, s) }
func yellowText(s string) string  { return colorize(yellow, s) }
func blueText(s string) string    { return colorize(blue, s) }
func magentaText(s string) string { return colorize(magenta, s) }
func cyanText(s string) string    { return colorize(cyan, s) }
func grayText(s string) string    { return colorize(gray, s) }

// Printing helpers
func printSuccess(msg string) { fmt.Println(greenText("✓ ") + msg) }
func printError(msg string)   { fmt.Fprintln(os.Stderr, redText("✗ ")+msg) }
func printInfo(msg string)    { fmt.Println(cyanText("→ ") + msg) }
func printWarn(msg string)    { fmt.Println(yellowText("⚠ ") + msg) }

func printBanner() {
	fmt.Println()
	fmt.Println(boldText(cyanText("  ██╗    ██╗██╗   ██╗███████╗██╗███████╗    ███████╗██╗   ██╗██████╗ ███████╗")))
	fmt.Println(boldText(cyanText("  ██║    ██║╚██╗ ██╔╝╚══███╔╝██║██╔════╝    ██╔════╝██║   ██║██╔══██╗██╔════╝")))
	fmt.Println(boldText(cyanText("  ██║ █╗ ██║ ╚████╔╝   ███╔╝ ██║█████╗      ███████╗██║   ██║██████╔╝███████╗")))
	fmt.Println(boldText(cyanText("  ██║███╗██║  ╚██╔╝   ███╔╝  ██║██╔══╝      ╚════██║██║   ██║██╔══██╗╚════██║")))
	fmt.Println(boldText(cyanText("  ╚███╔███╔╝   ██║   ███████╗██║███████╗    ███████║╚██████╔╝██████╔╝███████║")))
	fmt.Println(boldText(cyanText("   ╚══╝╚══╝    ╚═╝   ╚══════╝╚═╝╚══════╝    ╚══════╝ ╚═════╝ ╚═════╝ ╚══════╝")))
	fmt.Println(grayText("                           Subtitle Downloader powered by Wyzie"))
	fmt.Println()
}
