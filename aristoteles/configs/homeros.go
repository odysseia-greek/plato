package configs

import (
	"github.com/odysseia-greek/plato/cache"
	"github.com/odysseia-greek/plato/service"
)

type HomerosConfig struct {
	HttpClients service.OdysseiaClient
	Cache       cache.Client
}
