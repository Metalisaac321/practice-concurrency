package bootstrap

import (
	"net/http"
	"time"

	"github.com/Metalisaac321/practice-concurrency/internal/platform/server"
	"github.com/Metalisaac321/practice-concurrency/internal/pokemon/application"
	"github.com/Metalisaac321/practice-concurrency/internal/pokemon/platform"
)

const (
	host = "localhost"
	port = 8080
)

func Run() error {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	pokemonApiClient := platform.NewPokemonApiClient(httpClient)
	searchPokemons := application.NewSearchPokemons(pokemonApiClient)
	srv := server.New(host, port, searchPokemons)
	return srv.Run()
}
