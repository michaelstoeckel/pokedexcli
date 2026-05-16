package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/michaelstoeckel/pokedexcli/internal/pokecache"
)

type location struct {
	// Capitalize Name and Url
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Locations struct {
	// Capitalize all fields
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []location `json:"results"`
}

func GetLocations(url string, cache *pokecache.Cache) (Locations, error) {
	// get data from cache if available
	data, ok := cache.Get(url)
	// not in cache get data from url
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return Locations{}, err // Return the error instead of killing the app
		}
		defer res.Body.Close() // Ensure closure happens

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return Locations{}, err
		}
		// add to cache
		cache.Add(url, data)
	}

	loca := Locations{}
	err := json.Unmarshal(data, &loca)
	if err != nil {
		return Locations{}, err
	}
	return loca, nil

}
