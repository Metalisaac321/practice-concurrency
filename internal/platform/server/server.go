package server

import (
	"fmt"
	"log"

	"github.com/Metalisaac321/practice-concurrency/internal/platform/server/handler/pokemons"
	"github.com/Metalisaac321/practice-concurrency/internal/pokemon/application"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpAddr string
	engine   *gin.Engine

	//deps
	searchPokemons application.SearchPokemons
}

func New(host string, port uint, searchPokemons application.SearchPokemons) Server {
	srv := Server{
		engine:         gin.New(),
		httpAddr:       fmt.Sprintf("%s:%d", host, port),
		searchPokemons: searchPokemons,
	}

	srv.registerRoutes()
	return srv
}

func (s *Server) Run() error {
	log.Println("Server running on", s.httpAddr)
	return s.engine.Run(s.httpAddr)
}

func (s *Server) registerRoutes() {
	s.engine.GET("/pokemons", pokemons.SearchPokemonsHandler(s.searchPokemons))
}
