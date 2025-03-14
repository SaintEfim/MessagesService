package websocket

import (
	"net/http"

	"MessagesService/config"

	"github.com/rs/cors"
)

func CorsSettings(cfg *config.Config) *cors.Cors {
	c := cors.New(cors.Options{
		AllowedMethods: []string{http.MethodOptions, http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowedOrigins: cfg.Cors.AllowedOrigins,
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		Debug: true,
	})

	return c
}
