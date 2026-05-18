package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/michaelstoeckel/pokedexcli/internal/pokecache"
)

type Pokemon struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	BaseExperience int     `json:"base_experience"`
	Height         int     `json:"height"`
	Weight         int     `json:"weight"`
	Stats          []Stats `json:"stats"`
	Types          []Types `json:"types"`
}
type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Stats struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	Stat     Stat `json:"stat"`
}
type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Types struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}

func GetPokemon(url string, cache *pokecache.Cache) (Pokemon, error) {
	// get data from cache if available
	data, ok := cache.Get(url)
	// not in cache get data from url
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return Pokemon{}, err // Return the error instead of killing the app
		}
		defer res.Body.Close() // Ensure closure happens

		if res.StatusCode > 299 {
			return Pokemon{}, fmt.Errorf("pokemon not found!\n")
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return Pokemon{}, err
		}
		// add to cache
		cache.Add(url, data)
	}

	pokemon := Pokemon{}
	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}
	return pokemon, nil
}
