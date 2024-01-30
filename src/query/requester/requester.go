package requester

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var Client = http.Client{
	Timeout: 4 * time.Second,
}

// Right now is not being used because python is saving the prices
func GetRequestPrice(asset string, vsCurrency string) (string, error) {
	var sb strings.Builder
	sb.WriteString("https://api.coingecko.com/api/v3/simple/price?ids=")
	sb.WriteString(asset)
	sb.WriteString("&vs_currencies=")
	sb.WriteString(vsCurrency)

	resp, err := Client.Get(sb.String())
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil || len(string(body)) == 0 {
		return "", err
	}

	return string(body), nil
}

func MakeGetRequest(endpoint string, query string) (string, error) {
	url := endpoint + query
	resp, err := Client.Get(url)
	if err != nil {
		return "", fmt.Errorf("endpoint is down, %s", err)
	}

	// Handle 404 responses from cosmos api, it's actually element not found
	if resp.StatusCode == 404 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		// endpoint error
		if strings.Contains(string(body), "Cannot GET") {
			return "", fmt.Errorf("endpoint is down")
		}
		// node element not found
		// return `{"error": "Element not found"}`, nil
		return "", fmt.Errorf("element not found")
	}

	if resp.StatusCode == 400 {
		return "", fmt.Errorf("bad request")
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("endpoint is down")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil || len(string(body)) == 0 {
		return "", fmt.Errorf("bad encoding")
	}

	return string(body), nil
}
