package service

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type DionysiosImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

func NewDionysiosConfig(schema OdysseiaApi, ca []byte) (*DionysiosImpl, error) {
	client := NewHttpClient(ca, schema.Cert)
	return &DionysiosImpl{Scheme: schema.Scheme, BaseUrl: schema.Url, Client: client}, nil
}

func NewFakeDionysiosConfig(scheme, baseUrl string, client HttpClient) (*DionysiosImpl, error) {
	return &DionysiosImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (d *DionysiosImpl) Health() (*http.Response, error) {
	healthPath := url.URL{
		Scheme: d.Scheme,
		Host:   d.BaseUrl,
		Path:   path.Join(dionysiosService, version, healthEndPoint),
	}

	return Health(healthPath, d.Client)
}

func (d *DionysiosImpl) Grammar(word string) (*http.Response, error) {
	urlPath := url.URL{
		Scheme: d.Scheme,
		Host:   d.BaseUrl,
		Path:   path.Join(dionysiosService, version, grammar),
	}

	urlPath.Query().Set(searchWord, word)

	response, err := d.Client.Get(&urlPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, urlPath)
	}

	return response, nil
}
