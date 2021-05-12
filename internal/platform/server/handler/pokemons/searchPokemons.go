package pokemons

import (
	"net/http"

	"github.com/Metalisaac321/practice-concurrency/internal/pokemon"
	"github.com/Metalisaac321/practice-concurrency/internal/pokemon/application"
	"github.com/gin-gonic/gin"
)

type searchPokemonResponse struct {
	Pokemons []pokemon.Pokemon `json:"pokemons"`
}

func SearchPokemonsHandler(searchPokemons application.SearchPokemons) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pokemons, err := searchPokemons.Execute(ctx)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
		}

		ctx.JSON(200, searchPokemonResponse{Pokemons: pokemons})
	}
}
