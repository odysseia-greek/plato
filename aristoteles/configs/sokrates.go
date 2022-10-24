package configs

import (
	"github.com/odysseia-greek/plato/elastic"
	"github.com/odysseia/aristoteles"
)

type SokratesConfig struct {
	Elastic    elastic.Client
	Randomizer aristoteles.Random
	SearchWord string
	Index      string
}
