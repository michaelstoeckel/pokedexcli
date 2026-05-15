package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
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

func GetLocations(url string) (Locations, error) {

	res, err := http.Get(url)
	if err != nil {
		return Locations{}, err // Return the error instead of killing the app
	}
	defer res.Body.Close() // Ensure closure happens

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Locations{}, err
	}

	loca := Locations{}
	err = json.Unmarshal(body, &loca)
	if err != nil {
		return Locations{}, err
	}
	return loca, nil

}
