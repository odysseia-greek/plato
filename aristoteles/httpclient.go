package aristoteles

import (
	"crypto/tls"
	"errors"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/service"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

func (c *Config) getOdysseiaClient() (service.OdysseiaClient, error) {
	solonUrl := c.getStringFromEnv(EnvSolonService, defaultServiceAddress)
	PtolemaiosUrl := c.getStringFromEnv(EnvPtolemaiosService, defaultServiceAddress)
	HerodotosUrl := c.getStringFromEnv(EnvHerodotosService, defaultServiceAddress)
	AlexandrosUrl := c.getStringFromEnv(EnvAlexandrosService, defaultServiceAddress)
	DionysiosUrl := c.getStringFromEnv(EnvDionysiosService, defaultServiceAddress)
	SokratesUrl := c.getStringFromEnv(EnvSokratesService, defaultServiceAddress)
	tlsEnabled := c.getBoolFromEnv(EnvTlSKey)

	var ca []byte

	glg.Debug("getting odysseia client")

	config := service.ClientConfig{
		Ca:         ca,
		Solon:      service.OdysseiaApi{},
		Ptolemaios: service.OdysseiaApi{},
		Herodotos:  service.OdysseiaApi{},
		Dionysios:  service.OdysseiaApi{},
		Alexandros: service.OdysseiaApi{},
		Sokrates:   service.OdysseiaApi{},
	}

	if tlsEnabled {
		glg.Debug("setting up certs because TLS is enabled")
		rootPath := os.Getenv("CERT_ROOT")
		glg.Debugf("rootPath: %s", rootPath)
		dirs, err := ioutil.ReadDir(rootPath)
		if err != nil {
			glg.Error(err)
			return nil, err
		}

		for _, dir := range dirs {
			if dir.IsDir() {
				dirPath := filepath.Join(rootPath, dir.Name())
				glg.Debugf("found directory: %s", dirPath)

				certPath := filepath.Join(dirPath, "tls.crt")
				keyPath := filepath.Join(dirPath, "tls.key")

				if _, err := os.Stat(certPath); errors.Is(err, os.ErrNotExist) {
					glg.Debugf("cannot get file because it does not exist: %s", certPath)
					continue
				}

				if _, err := os.Stat(keyPath); errors.Is(err, os.ErrNotExist) {
					glg.Debugf("cannot get file because it does not exist: %s", keyPath)
					continue
				}

				loadedCerts, err := tls.LoadX509KeyPair(certPath, keyPath)
				if err != nil {
					glg.Error(err)
					return nil, err
				}

				if ca == nil {
					caPath := filepath.Join(rootPath, dir.Name(), "tls.pem")
					if _, err := os.Stat(caPath); errors.Is(err, os.ErrNotExist) {
						glg.Debugf("cannot get file because it does not exist: %s", caPath)
						continue
					}
					ca, _ = ioutil.ReadFile(caPath)
					glg.Debugf("writing CA for path %s", caPath)
				}

				switch dir.Name() {
				case "solon":
					parsed, err := url.Parse(solonUrl)
					if err != nil {
						return nil, err
					}
					config.Solon.Scheme = parsed.Scheme
					config.Solon.Url = parsed.Host
					config.Solon.Cert = []tls.Certificate{loadedCerts}
				case "ptolemaios":
					parsed, err := url.Parse(PtolemaiosUrl)
					if err != nil {
						return nil, err
					}
					config.Ptolemaios.Scheme = parsed.Scheme
					config.Ptolemaios.Url = parsed.Host
					config.Ptolemaios.Cert = []tls.Certificate{loadedCerts}
				case "dionysios":
					parsed, err := url.Parse(DionysiosUrl)
					if err != nil {
						return nil, err
					}
					config.Dionysios.Scheme = parsed.Scheme
					config.Dionysios.Url = parsed.Host
					config.Dionysios.Cert = []tls.Certificate{loadedCerts}
				case "herodotos":
					parsed, err := url.Parse(HerodotosUrl)
					if err != nil {
						return nil, err
					}
					config.Herodotos.Scheme = parsed.Scheme
					config.Herodotos.Url = parsed.Host
					config.Herodotos.Cert = []tls.Certificate{loadedCerts}
				case "alexandros":
					parsed, err := url.Parse(AlexandrosUrl)
					if err != nil {
						return nil, err
					}
					config.Alexandros.Scheme = parsed.Scheme
					config.Alexandros.Url = parsed.Host
					config.Alexandros.Cert = []tls.Certificate{loadedCerts}
				case "sokrates":
					parsed, err := url.Parse(SokratesUrl)
					if err != nil {
						return nil, err
					}
					config.Sokrates.Scheme = parsed.Scheme
					config.Sokrates.Url = parsed.Host
					config.Sokrates.Cert = []tls.Certificate{loadedCerts}
				}
			}
		}
	}

	glg.Debug("creating new client")

	return service.NewClient(config)
}
