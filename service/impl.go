package service

import (
	"crypto/tls"
	"github.com/odysseia-greek/plato/models"
)

type OdysseiaClient interface {
	Solon() Solon
	Ptolemaios() Ptolemaios
	Herodotos() Herodotos
}

type Odysseia struct {
	solon      *SolonImpl
	ptolemaios *PtolemaiosImpl
	herodotos  *HerodotosImpl
}

type Solon interface {
	Health() (*models.Health, error)
	OneTimeToken() (*models.TokenResponse, error)
	Register(requestBody models.SolonCreationRequest) (*models.SolonResponse, error)
}

type Ptolemaios interface {
	GetSecret() (*models.ElasticConfigVault, error)
}

type Herodotos interface {
	GetAuthors() (*models.Authors, error)
	GetBooks(authorId string) (*models.Books, error)
	CreateQuestion(author, book string) (*models.CreateSentenceResponse, error)
	CheckSentence(checkSentenceRequest models.CheckSentenceRequest) (*models.CheckSentenceResponse, error)
}

type ClientConfig struct {
	Ca         []byte
	Solon      OdysseiaApi
	Ptolemaios OdysseiaApi
	Herodotos  OdysseiaApi
	Dionysios  OdysseiaApi
	Alexandros OdysseiaApi
	Sokrates   OdysseiaApi
}

type OdysseiaApi struct {
	Url    string
	Scheme string
	Cert   []tls.Certificate
}

func NewClient(config ClientConfig) (OdysseiaClient, error) {
	solonImpl, err := NewSolonImpl(config.Solon, config.Ca)
	if err != nil {
		return nil, err
	}

	ptolemaiosImpl, err := NewPtolemaiosConfig(config.Ptolemaios, config.Ca)
	if err != nil {
		return nil, err
	}

	herodotosImpl, err := NewHerodotosConfig(config.Herodotos, config.Ca)
	if err != nil {
		return nil, err
	}

	return &Odysseia{
		solon:      solonImpl,
		ptolemaios: ptolemaiosImpl,
		herodotos:  herodotosImpl,
	}, nil
}

func NewFakeClient(config ClientConfig, codes []int, responses []string) (OdysseiaClient, error) {
	client := NewFakeHttpClient(responses, codes)

	solonImpl, err := NewFakeSolonImpl(config.Solon.Scheme, config.Solon.Url, client)
	if err != nil {
		return nil, err
	}

	ptolemaiosImpl, err := NewFakePtolemaiosConfig(config.Ptolemaios.Scheme, config.Ptolemaios.Url, client)
	if err != nil {
		return nil, err
	}

	herodotosImpl, err := NewFakeHerodotosConfig(config.Herodotos.Scheme, config.Herodotos.Url, client)
	if err != nil {
		return nil, err
	}

	return &Odysseia{
		solon:      solonImpl,
		ptolemaios: ptolemaiosImpl,
		herodotos:  herodotosImpl,
	}, nil
}

func (o *Odysseia) Solon() Solon {
	if o == nil {
		return nil
	}
	return o.solon
}

func (o *Odysseia) Ptolemaios() Ptolemaios {
	if o == nil {
		return nil
	}
	return o.ptolemaios
}

func (o *Odysseia) Herodotos() Herodotos {
	if o == nil {
		return nil
	}
	return o.herodotos
}
