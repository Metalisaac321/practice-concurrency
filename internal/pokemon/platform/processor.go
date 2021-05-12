package platform

import (
	"context"

	"github.com/Metalisaac321/practice-concurrency/internal/pokemon"
)

// PokemonDetailProcessor
type PokemonDetailProcessor struct {
	pokemonsDetailChan chan pokemon.PokemonDetail
	errorsChan         chan error
	pokemonApiClient   pokemon.PokemonApiClient
}

func (p *PokemonDetailProcessor) launch(ctx context.Context, pokemonsDto []pokemon.PokemonDto) {
	for _, pokemonDto := range pokemonsDto {
		go func(pokemonName string) {
			pokemonDetailDto, err := p.pokemonApiClient.GetPokemonDetail(ctx, pokemonName)
			if err != nil {
				p.errorsChan <- err
				return
			}

			totalAbilities := len(pokemonDetailDto.Abilities)
			pokemonAbilitiesProcessor := PokemonAbilitiesProcessor{
				PokemonDetailProcessor: PokemonDetailProcessor{
					errorsChan:       make(chan error, totalAbilities),
					pokemonApiClient: p.pokemonApiClient,
				},
				pokemonAbilityChan: make(chan pokemon.Ability, totalAbilities),
			}
			pokemonAbilitiesProcessor.launch(ctx, pokemonDetailDto.Abilities)

			abilities, err := pokemonAbilitiesProcessor.waitForAbilities(ctx, totalAbilities)
			if err != nil {
				pokemonAbilitiesProcessor.errorsChan <- err
				return
			}

			pokemonDetail := pokemon.PokemonDetail{
				Id:                     pokemonDetailDto.Id,
				Name:                   pokemonDetailDto.Name,
				Order:                  pokemonDetailDto.Order,
				Height:                 pokemonDetailDto.Height,
				Weight:                 pokemonDetailDto.Weight,
				Abilities:              abilities,
				BaseExperience:         pokemonDetailDto.BaseExperience,
				Forms:                  pokemonDetailDto.Forms,
				GameIndices:            pokemonDetailDto.GameIndices,
				HeldItems:              pokemonDetailDto.HeldItems,
				IsDefault:              pokemonDetailDto.IsDefault,
				LocationAreaEncounters: pokemonDetailDto.LocationAreaEncounters,
				Moves:                  pokemonDetailDto.Moves,
				PastTypes:              pokemonDetailDto.PastTypes,
				Stats:                  pokemonDetailDto.Stats,
				Types:                  pokemonDetailDto.Types,
			}

			p.pokemonsDetailChan <- pokemonDetail
		}(pokemonDto.Name)
	}
}

func (p *PokemonDetailProcessor) waitForPokemonsDetails(ctx context.Context, totalPokemonsFetched int) ([]pokemon.PokemonDetail, error) {
	var pokemonsDetails []pokemon.PokemonDetail

	for totalPokemonsFetched > 0 {
		select {
		case pokemonDetail := <-p.pokemonsDetailChan:
			pokemonsDetails = append(pokemonsDetails, pokemonDetail)
			totalPokemonsFetched--
		case err := <-p.errorsChan:
			return []pokemon.PokemonDetail{}, err
		case <-ctx.Done():
			return []pokemon.PokemonDetail{}, ctx.Err()
		}
	}

	return pokemonsDetails, nil
}

// PokemonAbilitiesProcessor
type PokemonAbilitiesProcessor struct {
	PokemonDetailProcessor
	pokemonAbilityChan chan pokemon.Ability
}

func (p *PokemonAbilitiesProcessor) launch(ctx context.Context, abilities []pokemon.AbilityDto) {
	for _, ability := range abilities {
		go func(url string) {
			pokemonAbility, err := p.pokemonApiClient.GetPokemonAbility(ctx, url)
			if err != nil {
				p.errorsChan <- err
				return
			}
			p.pokemonAbilityChan <- pokemonAbility
		}(ability.Ability.Url)
	}
}

func (p *PokemonAbilitiesProcessor) waitForAbilities(ctx context.Context, totalAbilities int) ([]pokemon.Ability, error) {
	var abilities []pokemon.Ability

	for totalAbilities > 0 {
		select {
		case ability := <-p.pokemonAbilityChan:
			abilities = append(abilities, ability)
			totalAbilities--
		case err := <-p.errorsChan:
			return []pokemon.Ability{}, err
		case <-ctx.Done():
			return []pokemon.Ability{}, ctx.Err()
		}
	}

	return abilities, nil
}
