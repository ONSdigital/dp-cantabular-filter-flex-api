package api

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"

	"github.com/gorilla/mux"
)

// API provides a struct to wrap the api around
type API struct {
	Router    *mux.Router
	store     datastore
	respond   responder
	generate  generator
	cfg       *config.Config
}

// New creates and initialises a new API
func New(ctx context.Context, cfg *config.Config, r *mux.Router, rsp responder, g generator, d datastore) *API {
	api := &API{
		Router:    r,
		respond:   rsp,
		generate:  g,
		store:     d,
		cfg:       cfg,
	}

	r.HandleFunc("/filters", api.createFilter).Methods("POST")

	return api
}
