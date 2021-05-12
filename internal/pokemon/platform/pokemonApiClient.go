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

func (pokemonApiClient *PokemonApiClient) GetPokemons(ctx context.Context) ([]pokemon.PokemonDto, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%spokemon", POKEMON_API_BASE_URL), nil)
	if err != nil {
		return []pokemon.PokemonDto{}, fmt.Errorf("Error creating request: %v", err)
	}

	response, err := pokemonApiClient.httpClient.Do(request)
	if err != nil {
		return []pokemon.PokemonDto{}, fmt.Errorf("Error making request: %v", err)
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return []pokemon.PokemonDto{}, fmt.Errorf("Unexpected status: got %v", response.Status)
	}

	var getPokemonsResponseDto GetPokemonsResponseDto
	err = json.NewDecoder(response.Body).Decode(&getPokemonsResponseDto)
	if err != nil {
		return []pokemon.PokemonDto{}, fmt.Errorf("Error decoding: %v", err)
	}

	return getPokemonsResponseDto.Results, nil
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
