package service

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/plato/models"
	"net/http"
	"net/url"
	"path"
)

type SokratesImpl struct {
	Scheme  string
	BaseUrl string
	Client  HttpClient
}

func NewSokratesConfig(schema OdysseiaApi, ca []byte) (*SokratesImpl, error) {
	client := NewHttpClient(ca, schema.Cert)
	return &SokratesImpl{Scheme: schema.Scheme, BaseUrl: schema.Url, Client: client}, nil
}

func NewFakeSokratesConfig(scheme, baseUrl string, client HttpClient) (*SokratesImpl, error) {
	return &SokratesImpl{Scheme: scheme, BaseUrl: baseUrl, Client: client}, nil
}
func (s *SokratesImpl) Health() (*http.Response, error) {
	healthPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   path.Join(sokratesService, version, healthEndPoint),
	}

	return Health(healthPath, s.Client)
}

func (s *SokratesImpl) GetMethods() (*http.Response, error) {
	methodPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s", sokratesService, version, methods),
	}

	response, err := s.Client.Get(&methodPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, methodPath)
	}

	return response, nil
}

func (s *SokratesImpl) GetCategories(method string) (*http.Response, error) {
	categoryPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s/%s/%s", sokratesService, version, methods, method, categories),
	}

	response, err := s.Client.Get(&categoryPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, categoryPath)
	}

	return response, nil
}

func (s *SokratesImpl) GetChapters(method, category string) (*http.Response, error) {
	chapterPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s", sokratesService, version, methods, method, categories, category, chapters),
	}

	response, err := s.Client.Get(&chapterPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, chapterPath)
	}

	return response, nil
}

func (s *SokratesImpl) CreateQuestion(method, category, chapter string) (*http.Response, error) {
	query := fmt.Sprintf("method=%s&category=%s&chapter=%s", method, category, chapter)
	questionPath := url.URL{
		Scheme:   s.Scheme,
		Host:     s.BaseUrl,
		Path:     fmt.Sprintf("%s/%s/%s", sokratesService, version, question),
		RawQuery: query,
	}

	response, err := s.Client.Get(&questionPath)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, questionPath)
	}

	return response, nil
}

func (s *SokratesImpl) CheckAnswer(checkAnswerRequest models.CheckAnswerRequest) (*http.Response, error) {
	answerPath := url.URL{
		Scheme: s.Scheme,
		Host:   s.BaseUrl,
		Path:   fmt.Sprintf("%s/%s/%s", sokratesService, version, answer),
	}

	body, err := json.Marshal(checkAnswerRequest)
	if err != nil {
		return nil, err
	}

	response, err := s.Client.Post(&answerPath, body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v but got %v while calling %v endpoint", http.StatusOK, response.StatusCode, answerPath)
	}

	return response, nil
}
