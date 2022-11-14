package aristoteles

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	PLATO         = "plato"
	ODYSSEIA_PATH = "odysseia-greek"
)

func (c *Config) GetPlatoPath(path string) string {
	dirs, _ := ioutil.ReadDir(path)
	var currentModTime time.Time
	var latestPath string
	for _, dir := range dirs {

		if dir.Name() == PLATO {
			latestPath = dir.Name()
			break
		}

		if currentModTime.Before(dir.ModTime()) {
			currentModTime = dir.ModTime()
			latestPath = dir.Name()
		}
	}

	return filepath.Join(path, latestPath)
}

func (c *Config) OdysseiaRootPath(path string) string {
	_, callingFile, _, _ := runtime.Caller(0)
	callingDir := filepath.Dir(callingFile)
	dirParts := strings.Split(callingDir, string(os.PathSeparator))
	var odysseiaPath []string
	for i, part := range dirParts {
		if part == path {
			odysseiaPath = dirParts[0 : i+1]
		}
	}
	l := "/"
	for _, path := range odysseiaPath {
		l = filepath.Join(l, path)
	}

	return l
}
