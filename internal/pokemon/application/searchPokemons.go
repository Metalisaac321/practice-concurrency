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

type processor struct {
	pokemonsDetailChan chan pokemon.PokemonDetailDto
	errorsChan         chan error
	pokemonApiClient   pokemon.PokemonApiClient
}

func (p *processor) launch(ctx context.Context, pokemonsDto []pokemon.PokemonDto) {
	for _, pokemonDto := range pokemonsDto {
		go func(pokemonName string) {
			pokemonDetail, err := p.pokemonApiClient.GetPokemonDetail(ctx, pokemonName)
			if err != nil {
				p.errorsChan <- err
				return
			}
			p.pokemonsDetailChan <- pokemonDetail
		}(pokemonDto.Name)
	}
}

func (p *processor) waitForPokemonsDetails(ctx context.Context, totalPokemonsFetched int) ([]pokemon.PokemonDetailDto, error) {
	var pokemonsDetails []pokemon.PokemonDetailDto

	for totalPokemonsFetched > 0 {
		select {
		case pokemonDetail := <-p.pokemonsDetailChan:
			pokemonsDetails = append(pokemonsDetails, pokemonDetail)
			totalPokemonsFetched--
		case err := <-p.errorsChan:
			return []pokemon.PokemonDetailDto{}, err
		case <-ctx.Done():
			return []pokemon.PokemonDetailDto{}, ctx.Err()
		}
	}

	return pokemonsDetails, nil
}

func (useCase SearchPokemons) Execute(ctx context.Context) ([]pokemon.Pokemon, error) {
	pokemonsFromApi, err := useCase.pokemonApiClient.GetPokemons(ctx)
	if err != nil {
		return []pokemon.Pokemon{}, err
	}
	totalPokemonsFetched := len(pokemonsFromApi)

	processor := processor{
		pokemonsDetailChan: make(chan pokemon.PokemonDetailDto, totalPokemonsFetched),
		errorsChan:         make(chan error, totalPokemonsFetched),
		pokemonApiClient:   useCase.pokemonApiClient,
	}
	processor.launch(ctx, pokemonsFromApi)
	pokemonsDetails, err := processor.waitForPokemonsDetails(ctx, totalPokemonsFetched)

	if err != nil {
		return []pokemon.Pokemon{}, err
	}

	var pokemons []pokemon.Pokemon
	for _, pokemonDetail := range pokemonsDetails {
		var abilities []pokemon.Ability
		for _, abilityDto := range pokemonDetail.Abilities {
			abilities = append(abilities, pokemon.Ability{
				Id:   0,
				Name: abilityDto.Ability.Name,
			})
		}

		pokemons = append(pokemons, pokemon.NewPokemon(
			pokemon.Pokemon{
				Id:        pokemonDetail.Id,
				Name:      pokemonDetail.Name,
				Height:    pokemonDetail.Height,
				Weight:    pokemonDetail.Weight,
				Abilities: abilities,
			},
		))
	}
	return pokemons, nil
}
