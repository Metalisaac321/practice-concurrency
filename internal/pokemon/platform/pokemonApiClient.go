package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Metalisaac321/practice-concurrency/internal/pokemon"
)

const (
	POKEMON_API_BASE_URL = "https://pokeapi.co/api/v2/"
)

type PokemonApiClient struct {
	httpClient *http.Client
}

func NewPokemonApiClient(httpClient *http.Client) *PokemonApiClient {
	return &PokemonApiClient{
		httpClient: httpClient,
	}
}

// GetPokemons
type GetPokemonsResponseDto struct {
	Results []pokemon.PokemonDto `json:"results"`
}

func (pokemonApiClient *PokemonApiClient) GetPokemons(ctx context.Context) ([]pokemon.Pokemon, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%spokemon", POKEMON_API_BASE_URL), nil)
	if err != nil {
		return []pokemon.Pokemon{}, fmt.Errorf("Error creating request: %v", err)
	}

	response, err := pokemonApiClient.httpClient.Do(request)
	if err != nil {
		return []pokemon.Pokemon{}, fmt.Errorf("Error making request: %v", err)
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return []pokemon.Pokemon{}, fmt.Errorf("Unexpected status: got %v", response.Status)
	}

	var getPokemonsResponseDto GetPokemonsResponseDto
	err = json.NewDecoder(response.Body).Decode(&getPokemonsResponseDto)
	if err != nil {
		return []pokemon.Pokemon{}, fmt.Errorf("Error decoding: %v", err)
	}

	pokemonsFromApi := getPokemonsResponseDto.Results
	totalPokemonsFetched := len(pokemonsFromApi)

	pokemonDetailProcessor := PokemonDetailProcessor{
		pokemonsDetailChan: make(chan pokemon.PokemonDetail, totalPokemonsFetched),
		errorsChan:         make(chan error, totalPokemonsFetched),
		pokemonApiClient:   pokemonApiClient,
	}
	pokemonDetailProcessor.launch(ctx, pokemonsFromApi)

	pokemonsDetails, err := pokemonDetailProcessor.waitForPokemonsDetails(ctx, totalPokemonsFetched)
	if err != nil {
		return []pokemon.Pokemon{}, fmt.Errorf("Error getting pokemond details: %v", err)
	}

	var pokemons []pokemon.Pokemon
	for _, pokemonDetail := range pokemonsDetails {
		pokemons = append(pokemons, pokemon.NewPokemon(
			pokemon.Pokemon{
				Id:        pokemonDetail.Id,
				Name:      pokemonDetail.Name,
				Height:    pokemonDetail.Height,
				Weight:    pokemonDetail.Weight,
				Order:     pokemonDetail.Order,
				Abilities: pokemonDetail.Abilities,
			},
		))
	}

	return pokemons, nil
}

// GetPokemonDetail
func (pokemonApiClient *PokemonApiClient) GetPokemonDetail(ctx context.Context, pokemonName string) (pokemon.PokemonDetailDto, error) {
	pokemonDetailDto := pokemon.PokemonDetailDto{}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%spokemon/%s", POKEMON_API_BASE_URL, pokemonName), nil)
	if err != nil {
		return pokemonDetailDto, fmt.Errorf("%v", err)
	}

	response, err := pokemonApiClient.httpClient.Do(request)
	if err != nil {
		return pokemonDetailDto, fmt.Errorf("%v", err)
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return pokemonDetailDto, fmt.Errorf("Unexpected status: got %v", response.Status)
	}

	err = json.NewDecoder(response.Body).Decode(&pokemonDetailDto)
	if err != nil {
		return pokemonDetailDto, fmt.Errorf("Error decoding json: %v", err)
	}

	return pokemonDetailDto, nil
}

// GetPokemonAbilities
func (pokemonApiClient *PokemonApiClient) GetPokemonAbility(ctx context.Context, url string) (pokemon.Ability, error) {
	ability := pokemon.Ability{}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ability, fmt.Errorf("in creating request: %v", err)
	}

	response, err := pokemonApiClient.httpClient.Do(request)
	if err != nil {
		return ability, fmt.Errorf("in making request%v", err)
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return ability, fmt.Errorf("Unexpected status: got %v", response.Status)
	}

	err = json.NewDecoder(response.Body).Decode(&ability)
	if err != nil {
		return ability, fmt.Errorf("Error decoding json: %v", err)
	}

	return ability, nil
}
