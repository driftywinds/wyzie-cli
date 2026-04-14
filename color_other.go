//go:build !windows

package main

func enableWindowsVT() bool { return false }
