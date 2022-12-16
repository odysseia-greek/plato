package service

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	version        string = "v1"
	healthEndPoint string = "health"
)

func Health(path url.URL, client HttpClient) (*http.Response, error) {
	response, err := client.Get(&path)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling token endpoint", http.StatusOK, response.StatusCode)
	}

	return response, nil
}
