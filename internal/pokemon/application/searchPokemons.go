package application

import (
	"context"

	"github.com/Metalisaac321/practice-concurrency/internal/pokemon"
)

type SearchPokemons struct {
	pokemonApiClient pokemon.PokemonApiClient
}

func NewSearchPokemons(pokemonApiClient pokemon.PokemonApiClient) SearchPokemons {
	return SearchPokemons{
		pokemonApiClient: pokemonApiClient,
	}
}

func (useCase SearchPokemons) Execute(ctx context.Context) ([]pokemon.Pokemon, error) {
	pokemons, err := useCase.pokemonApiClient.GetPokemons(ctx)
	if err != nil {
		return []pokemon.Pokemon{}, err
	}

	return pokemons, nil
}
