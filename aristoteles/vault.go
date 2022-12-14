package aristoteles

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/vault"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	pb "github.com/odysseia-greek/plato/proto"
)

const (
	VAULT = "vault"
)

func (c *Config) getConfigFromVault() (*pb.ElasticConfigVault, error) {
	sidecarService := os.Getenv(EnvPtolemaiosService)
	if sidecarService == "" {
		glg.Infof("defaulting to %s for sidecar", defaultSidecarService)
		sidecarService = defaultSidecarService
	}

	conn, err := grpc.Dial(sidecarService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewPtolemaiosClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.GetSecret(ctx, &pb.VaultRequest{})
	if err != nil {
		return nil, err
	}

	glg.Infof("found secret %s", r)

	return r, nil
}

func (c *Config) getVaultClient() (vault.Client, error) {
	var vaultClient vault.Client

	vaultRootToken := c.getStringFromEnv(EnvRootToken, "")
	vaultAuthMethod := c.getStringFromEnv(EnvAuthMethod, AuthMethodToken)
	vaultService := c.getStringFromEnv(EnvVaultService, c.BaseConfig.VaultService)
	vaultRole := c.getStringFromEnv(EnvVaultRole, defaultRoleName)
	tlsEnabled := c.getBoolFromEnv(EnvTLSEnabled)
	rootPath := c.getStringFromEnv(EnvRootTlSDir, defaultTLSFileLocation)
	secretPath := filepath.Join(rootPath, VAULT)

	glg.Debugf("vaultAuthMethod set to %s", vaultAuthMethod)
	glg.Debugf("secretPath set to %s", secretPath)
	glg.Debugf("tlsEnabled set to %v", tlsEnabled)

	var tlsConfig *api.TLSConfig

	if tlsEnabled {
		insecure := false
		if c.env == "LOCAL" || c.env == "TEST" {
			insecure = !insecure
			secretPath = "/tmp"
		}

		ca := fmt.Sprintf("%s/vault.ca", secretPath)
		cert := fmt.Sprintf("%s/vault.crt", secretPath)
		key := fmt.Sprintf("%s/vault.key", secretPath)

		tlsConfig = vault.CreateTLSConfig(insecure, ca, cert, key, secretPath)
	}

	if vaultAuthMethod == AuthMethodKube {
		jwtToken, err := os.ReadFile(serviceAccountTokenPath)
		if err != nil {
			glg.Error(err)
			return nil, err
		}

		vaultJwtToken := string(jwtToken)

		client, err := vault.CreateVaultClientKubernetes(vaultService, vaultRole, vaultJwtToken, tlsConfig)
		if err != nil {
			glg.Error(err)
			return nil, err
		}

		if c.BaseConfig.HealthCheckOverwrite {
			ticks := 120 * time.Second
			tick := 1 * time.Second
			healthy := vaultClient.CheckHealthyStatus(ticks, tick)
			if !healthy {
				return nil, fmt.Errorf("error getting healthy status from vault")
			}
		}

		vaultClient = client
	} else {
		if c.env == "LOCAL" || c.env == "TEST" {
			glg.Debug("local testing, getting token from file")
			localToken, err := c.getTokenFromFile(defaultNamespace)
			if err != nil {
				glg.Error(err)
				return nil, err
			}
			client, err := vault.NewVaultClient(vaultService, localToken, tlsConfig)
			if err != nil {
				glg.Error(err)
				return nil, err
			}

			vaultClient = client
		} else {
			client, err := vault.NewVaultClient(vaultService, vaultRootToken, tlsConfig)
			if err != nil {
				glg.Error(err)
				return nil, err
			}

			vaultClient = client
		}
	}

	glg.Debug(vaultClient)
	return vaultClient, nil
}

func (c *Config) getTokenFromFile(namespace string) (string, error) {
	rootPath := c.OdysseiaRootPath(ODYSSEIA_PATH)
	platoPath := c.GetPlatoPath(rootPath)
	clusterKeys := filepath.Join(platoPath, "vault", "eratosthenes", fmt.Sprintf("cluster-keys-%s.json", namespace))

	f, err := ioutil.ReadFile(clusterKeys)
	if err != nil {
		glg.Error(fmt.Sprintf("Cannot read fixture file: %s", err))
		return "", err
	}

	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(f, &result)

	return result["root_token"].(string), nil
}
