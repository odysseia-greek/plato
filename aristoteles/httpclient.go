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
	herodotosUrl := c.getStringFromEnv(EnvHerodotosService, defaultServiceAddress)
	alexandrosUrl := c.getStringFromEnv(EnvAlexandrosService, defaultServiceAddress)
	dionysiosUrl := c.getStringFromEnv(EnvDionysiosService, defaultServiceAddress)
	sokratesUrl := c.getStringFromEnv(EnvSokratesService, defaultServiceAddress)
	tlsEnabled := c.getBoolFromEnv(EnvTlSKey)

	solonParsed, err := url.Parse(solonUrl)
	if err != nil {
		return nil, err
	}

	herodotosParsed, err := url.Parse(herodotosUrl)
	if err != nil {
		return nil, err
	}
	alexandrosParsed, err := url.Parse(alexandrosUrl)
	if err != nil {
		return nil, err
	}
	dionysiosParsed, err := url.Parse(dionysiosUrl)
	if err != nil {
		return nil, err
	}
	sokratesParsed, err := url.Parse(sokratesUrl)
	if err != nil {
		return nil, err
	}

	var ca []byte

	glg.Debug("getting odysseia client")

	config := service.ClientConfig{
		Ca: ca,
		Solon: service.OdysseiaApi{
			Url:    solonParsed.Host,
			Scheme: solonParsed.Scheme,
			Cert:   nil,
		},
		Herodotos: service.OdysseiaApi{
			Url:    herodotosParsed.Host,
			Scheme: herodotosParsed.Scheme,
			Cert:   nil,
		},
		Dionysios: service.OdysseiaApi{
			Url:    dionysiosParsed.Host,
			Scheme: dionysiosParsed.Scheme,
			Cert:   nil,
		},
		Alexandros: service.OdysseiaApi{
			Url:    alexandrosParsed.Host,
			Scheme: alexandrosParsed.Scheme,
			Cert:   nil,
		},
		Sokrates: service.OdysseiaApi{
			Url:    sokratesParsed.Host,
			Scheme: sokratesParsed.Scheme,
			Cert:   nil,
		},
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
					config.Solon.Cert = []tls.Certificate{loadedCerts}
				case "dionysios":
					config.Dionysios.Cert = []tls.Certificate{loadedCerts}
				case "herodotos":
					config.Herodotos.Cert = []tls.Certificate{loadedCerts}
				case "alexandros":
					config.Alexandros.Cert = []tls.Certificate{loadedCerts}
				case "sokrates":
					config.Sokrates.Cert = []tls.Certificate{loadedCerts}
				}
			}
		}
	}

	glg.Debug("creating new client")

	return service.NewClient(config)
}
