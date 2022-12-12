package service

import (
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/models"
	"net/http"
	"net/url"
)

type HerodotosImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

const (
	pre      string = "/herodotos/v1"
	authors  string = "authors"
	books    string = "books"
	question string = "createQuestion"
	sentence string = "checkSentence"
)

func NewHerodotosConfig(schema OdysseiaApi, ca []byte) (*HerodotosImpl, error) {
	client := NewHttpClient(ca, schema.Cert)
	return &HerodotosImpl{Scheme: schema.Scheme, BaseUrl: schema.Url, Client: client}, nil
}

func NewFakeHerodotosConfig(scheme, baseUrl string, client HttpClient) (*HerodotosImpl, error) {
	return &HerodotosImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}

func (h *HerodotosImpl) GetAuthors() (*models.Authors, error) {
	path := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   fmt.Sprintf("%s/%s", pre, authors),
	}

	response, err := h.Client.Get(&path)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, path)
	}

	defer response.Body.Close()

	var authors models.Authors
	err = json.NewDecoder(response.Body).Decode(&authors)
	if err != nil {
		return nil, err
	}

	glg.Debug(authors)

	return &authors, nil
}

func (h *HerodotosImpl) GetBooks(authorId string) (*models.Books, error) {
	path := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s/%s", pre, authors, authorId, books),
	}

	response, err := h.Client.Get(&path)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, path)
	}

	defer response.Body.Close()

	var books models.Books
	err = json.NewDecoder(response.Body).Decode(&books)
	if err != nil {
		return nil, err
	}

	glg.Debug(books)

	return &books, nil
}

func (h *HerodotosImpl) CreateQuestion(author, book string) (*models.CreateSentenceResponse, error) {
	query := fmt.Sprintf("author=%s&book=%s", author, book)
	path := url.URL{
		Scheme:   h.Scheme,
		Host:     h.BaseUrl,
		Path:     fmt.Sprintf("%s/%s", pre, question),
		RawQuery: query,
	}

	response, err := h.Client.Get(&path)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, path)
	}

	defer response.Body.Close()

	var sentence models.CreateSentenceResponse
	err = json.NewDecoder(response.Body).Decode(&sentence)
	if err != nil {
		return nil, err
	}

	glg.Debug(sentence)

	return &sentence, nil
}

func (h *HerodotosImpl) CheckSentence(checkSentenceRequest models.CheckSentenceRequest) (*models.CheckSentenceResponse, error) {
	path := url.URL{
		Scheme: h.Scheme,
		Host:   h.BaseUrl,
		Path:   fmt.Sprintf("%s/%s", pre, sentence),
	}

	body, err := json.Marshal(checkSentenceRequest)
	if err != nil {
		return nil, err
	}

	response, err := h.Client.Post(&path, body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, path)
	}

	defer response.Body.Close()

	var sentence models.CheckSentenceResponse
	err = json.NewDecoder(response.Body).Decode(&sentence)
	if err != nil {
		return nil, err
	}

	glg.Debug(sentence)

	return &sentence, nil
}
