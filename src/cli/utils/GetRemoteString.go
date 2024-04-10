package utils

import (
	"io"
	"net/http"
	"fmt"
)

func GetRemoteString(url string) (string, error) {
	// Send an HTTP GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", err // Handle errors connecting to the server
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status) // Handle non-OK HTTP responses
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err // Handle errors reading the response body
	}

	// Convert the body bytes to a string
	bodyString := string(bodyBytes)
	return bodyString, nil
}
