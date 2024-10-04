package utils

import "net/url"

// CreateURL dynamically generates a URL with encoded query parameters
func CreateURL(baseURL string, params map[string]string) (string, error) {
	// Parse the base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// Initialize query parameters
	query := url.Values{}

	// Add each query parameter to the URL
	for key, value := range params {
		query.Add(key, value)
	}

	// Encode the parameters and append them to the URL
	parsedURL.RawQuery = query.Encode()

	// Return the full URL as a string
	return parsedURL.String(), nil
}
