package configs

import (
	"github.com/odysseia-greek/plato/elastic"
	"github.com/odysseia-greek/plato/randomizer"
)

type SokratesConfig struct {
	Elastic    elastic.Client
	Randomizer randomizer.Random
	SearchWord string
	Index      string
}
