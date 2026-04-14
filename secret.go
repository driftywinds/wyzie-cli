package main

// readSecret reads a line of input (without echo suppression - not needed for API keys).
func readSecret() (string, error) {
	return reader.ReadString('\n')
}
