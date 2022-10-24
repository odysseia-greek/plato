package aristoteles

import "github.com/odysseia-greek/plato/randomizer"

func (c *Config) getRandomizer() (randomizer.Random, error) {
	return randomizer.NewRandomizerClient()
}
