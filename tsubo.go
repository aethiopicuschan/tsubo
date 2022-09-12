package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/aethiopicuschan/tsubo/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	port := flag.Int("port", 5963, "port")
	flag.Parse()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	http.HandleFunc("/", api.NotFound)
	http.HandleFunc("/bbsmenu", api.BBSMenuAPI)
	http.HandleFunc("/subjects", api.SubjectsAPI)
	log.Info().Int("port", *port).Msg("Starting server")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal().Err(err).Msg("Server could not start")
	}
}
