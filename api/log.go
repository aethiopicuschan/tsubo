package api

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func accessLog(r *http.Request) {
	log.Info().Str("method", r.Method).Str("url", r.URL.String()).Msg("Access")
}

func warnLog(message string) {
	log.Warn().Msg(message)
}

func errorLog(err error) {
	log.Error().Msg(err.Error())
}
