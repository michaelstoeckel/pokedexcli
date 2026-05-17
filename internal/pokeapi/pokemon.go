package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/michaelstoeckel/pokedexcli/internal/pokecache"
)

type PokemonResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon        Pokemon         `json:"pokemon"`
	VersionDetails []VersionDetail `json:"version_details"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type VersionDetail struct {
	Version          Version           `json:"version"`
	MaxChance        int               `json:"max_chance"`
	EncounterDetails []EncounterDetail `json:"encounter_details"`
}

type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EncounterDetail struct {
	MinLevel        int             `json:"min_level"`
	MaxLevel        int             `json:"max_level"`
	ConditionValues []interface{}   `json:"condition_values"` // Falls hier Strings oder Objekte reinkommen, ggf. anpassen
	Chance          int             `json:"chance"`
	Method          EncounterMethod `json:"method"`
}

type EncounterMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetPokemonResponse(url string, cache *pokecache.Cache) (PokemonResponse, error) {
	// get data from cache if available
	data, ok := cache.Get(url)
	// not in cache get data from url
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return PokemonResponse{}, err // Return the error instead of killing the app
		}
		defer res.Body.Close() // Ensure closure happens

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return PokemonResponse{}, err
		}
		// add to cache
		cache.Add(url, data)
	}

	resp := PokemonResponse{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return PokemonResponse{}, err
	}
	return resp, nil
}
