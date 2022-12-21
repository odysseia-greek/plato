package aristoteles

import (
	uuid2 "github.com/google/uuid"
	"github.com/odysseia-greek/plato/cache"
	"path/filepath"
)

func (c *Config) getBadgerClient() (cache.Client, error) {
	uuid := uuid2.New().String()
	badgerPath := filepath.Join("/tmp", "badger", uuid)
	badger, err := cache.NewBadgerClient(badgerPath)
	if err != nil {
		return nil, err
	}

	return badger, nil
}
