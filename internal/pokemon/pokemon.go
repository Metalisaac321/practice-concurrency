package pokemon

import (
	"context"
)

// Pokemon
type Ability struct {
	Id           uint   `json:"id"`
	Name         string `json:"name"`
	IsMainSeries bool   `json:"isMainSeries"`
}

type Pokemon struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Order     uint      `json:"order"`
	Height    uint      `json:"height"`
	Weight    uint      `json:"weight"`
	Abilities []Ability `json:"abilities"`
}

func NewPokemon(pokemonDto Pokemon) Pokemon {
	return Pokemon{
		Id:        pokemonDto.Id,
		Name:      pokemonDto.Name,
		Height:    pokemonDto.Height,
		Weight:    pokemonDto.Weight,
		Abilities: pokemonDto.Abilities,
	}
}

// PokemonApiClient
type PokemonDto struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type AbilityDto struct {
	Ability struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"ability"`
}

type PokemonDetailDto struct {
	Id                     uint          `json:"id"`
	Name                   string        `json:"name"`
	Order                  uint          `json:"order"`
	Height                 uint          `json:"height"`
	Weight                 uint          `json:"weight"`
	Abilities              []AbilityDto  `json:"abilities"`
	BaseExperience         uint          `json:"base_experience"`
	Forms                  []interface{} `json:"forms"`
	GameIndices            []interface{} `json:"game_indices"`
	HeldItems              []interface{} `json:"held_items"`
	IsDefault              bool          `json:"is_default"`
	LocationAreaEncounters string        `json:"location_area_encounters"`
	Moves                  []interface{} `json:"moves"`
	PastTypes              []interface{} `json:"past_types"`
	Stats                  []interface{} `json:"stats"`
	Types                  []interface{} `json:"types"`
}

type PokemonDetail struct {
	Id                     uint          `json:"id"`
	Name                   string        `json:"name"`
	Order                  uint          `json:"order"`
	Height                 uint          `json:"height"`
	Weight                 uint          `json:"weight"`
	Abilities              []Ability     `json:"abilities"`
	BaseExperience         uint          `json:"base_experience"`
	Forms                  []interface{} `json:"forms"`
	GameIndices            []interface{} `json:"game_indices"`
	HeldItems              []interface{} `json:"held_items"`
	IsDefault              bool          `json:"is_default"`
	LocationAreaEncounters string        `json:"location_area_encounters"`
	Moves                  []interface{} `json:"moves"`
	PastTypes              []interface{} `json:"past_types"`
	Stats                  []interface{} `json:"stats"`
	Types                  []interface{} `json:"types"`
}

// PokemonApiClient
type PokemonApiClient interface {
	GetPokemons(ctx context.Context) ([]Pokemon, error)
	GetPokemonDetail(ctx context.Context, pokemonName string) (PokemonDetailDto, error)
	GetPokemonAbility(ctx context.Context, pokemonName string) (Ability, error)
}
