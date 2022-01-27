package api

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/middleware"
	"github.com/ONSdigital/dp-api-clients-go/v2/identity"
	"github.com/ONSdigital/dp-authorisation/auth"
	dphandlers "github.com/ONSdigital/dp-net/v2/handlers"

	"github.com/gorilla/mux"
)

// API provides a struct to wrap the api around
type API struct {
	Router         *mux.Router
	store          datastore
	respond        responder
	generate       generator
	identityClient *identity.Client
	cfg            *config.Config
}

// New creates and initialises a new API
func New(ctx context.Context, cfg *config.Config, r *mux.Router, idc *identity.Client, rsp responder, g generator, d datastore) *API {
	api := &API{
		Router:          r,
		respond:         rsp,
		generate:        g,
		store:           d,
		cfg:             cfg,
		identityClient:  idc,
	}

	if cfg.EnablePrivateEndpoints{
		api.enablePrivateEndpoints()
	} else {
		api.enablePublicEndpoints()
	}

	return api
}

func (api *API) enablePublicEndpoints(){
	api.Router.HandleFunc("/filters", api.createFilter).Methods("POST")
}
func (api *API) enablePrivateEndpoints(){
	checkIdentity := dphandlers.IdentityWithHTTPClient(api.identityClient)
	permissions := middleware.NewPermissions(api.cfg.ZebedeeURL, api.cfg.EnablePermissionsAuth)

	api.Router.Use(checkIdentity)
	api.Router.Use(middleware.LogIdentity())
	api.Router.Use(permissions.Require(auth.Permissions{Read: true}))

	api.Router.HandleFunc("/filters", api.createFilter).Methods("POST")
}
