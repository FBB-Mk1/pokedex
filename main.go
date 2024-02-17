package main

import (
	"time"

	"github.com/fbb-mk1/pokedex/internal/pkcache"
)

func main() {
	pokeCache := pkcache.NewCache(time.Minute * 5)
	StartRepl(pokeCache)
}
