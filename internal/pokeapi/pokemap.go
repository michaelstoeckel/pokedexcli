package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/michaelstoeckel/pokedexcli/internal/pokecache"
)

type PokemapResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon        PokemapPokemon  `json:"pokemon"`
	VersionDetails []VersionDetail `json:"version_details"`
}

type PokemapPokemon struct {
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

func GetPokemapResponse(url string, cache *pokecache.Cache) (PokemapResponse, error) {
	// get data from cache if available
	data, ok := cache.Get(url)
	// not in cache get data from url
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return PokemapResponse{}, err // Return the error instead of killing the app
		}
		defer res.Body.Close() // Ensure closure happens

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return PokemapResponse{}, err
		}
		// add to cache
		cache.Add(url, data)
	}

	resp := PokemapResponse{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return PokemapResponse{}, err
	}
	return resp, nil
}
