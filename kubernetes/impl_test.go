package kubernetes

import (
	"github.com/kpango/glg"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestClientCreation(t *testing.T) {
	ns := "test"
	defaultKubeConfig := "/.kube/config"

	t.Run("ConfigWhenNonIsSetAlsoWorks", func(t *testing.T) {
		var filePath string

		homeDir, err := os.UserHomeDir()
		if err != nil {
			glg.Error(err)
		}

		filePath = filepath.Join(homeDir, defaultKubeConfig)

		cfg, err := ioutil.ReadFile(filePath)
		assert.Nil(t, err)

		kube, err := NewKubeClient(cfg, ns)
		assert.NotNil(t, kube)
		assert.Nil(t, err)
	})

	t.Run("ConfigWhenNonIsSetAlsoWorks", func(t *testing.T) {
		var filePath string

		homeDir, err := os.UserHomeDir()
		if err != nil {
			glg.Error(err)
		}

		filePath = filepath.Join(homeDir, defaultKubeConfig)

		cfg, err := ioutil.ReadFile(filePath)
		assert.Nil(t, err)

		kube, err := NewKubeClient(cfg, ns)
		assert.NotNil(t, kube)
		assert.Nil(t, err)

		kubeContext, err := kube.Cluster().GetCurrentContext()
		assert.NotEqual(t, "", kubeContext)
		assert.Nil(t, err)
	})
}
