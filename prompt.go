package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

// prompt shows a label with an optional default and reads input from stdin.
func prompt(label, defaultVal string) (string, error) {
	if defaultVal != "" {
		fmt.Printf("%s %s(%s)%s: ", boldText(label), gray, defaultVal, reset)
	} else {
		fmt.Printf("%s: ", boldText(label))
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal, nil
	}
	return input, nil
}

// pickOne presents a numbered list and returns the 0-based index of the selection.
func pickOne(label string, items []string) (int, error) {
	fmt.Println()
	fmt.Println(boldText("  " + label))
	fmt.Println()
	for i, item := range items {
		num := fmt.Sprintf("%3d", i+1)
		fmt.Printf("  %s  %s\n", grayText(num+"."), item)
	}
	fmt.Println()

	for {
		raw, err := prompt("  Enter number", "1")
		if err != nil {
			return 0, err
		}
		n, err := strconv.Atoi(strings.TrimSpace(raw))
		if err != nil || n < 1 || n > len(items) {
			printWarn(fmt.Sprintf("Please enter a number between 1 and %d", len(items)))
			continue
		}
		return n - 1, nil
	}
}

func separator() {
	fmt.Println()
	fmt.Println(grayText(strings.Repeat("─", 80)))
	fmt.Println()
}
