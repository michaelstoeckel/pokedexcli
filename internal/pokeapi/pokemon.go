package pokeapi

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
