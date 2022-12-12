package aristoteles

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestFilePath(t *testing.T) {
	config := Config{}
	homeDir, err := os.UserHomeDir()
	assert.Nil(t, err)
	defaultPath := "odysseia-greek"

	t.Run("LocalFilePath", func(t *testing.T) {
		filePath := config.OdysseiaRootPath(defaultPath)
		sut := filepath.Join(homeDir, "/go/src/github.com/", defaultPath)
		assert.Equal(t, sut, filePath)
	})

	t.Run("EmptyPath", func(t *testing.T) {
		filePath := config.OdysseiaRootPath("")
		sut := "/"
		assert.Equal(t, sut, filePath)
	})

	t.Run("PlatoPathFlat", func(t *testing.T) {
		filePath := config.OdysseiaRootPath(defaultPath)
		platoPath := config.GetPlatoPath(filePath)
		sut := filepath.Join(homeDir, "/go/src/github.com/", defaultPath, PLATO)
		assert.Equal(t, sut, platoPath)
	})

	t.Run("PlatoPathSetDifferently", func(t *testing.T) {
		sutDir := "/go/pkg/mod/github.com/odysseia-greek"
		filePath := filepath.Join(homeDir, sutDir)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Skip()
		}
		platoPath := config.GetPlatoPath(filePath)
		assert.Contains(t, platoPath, PLATO)
	})
}
