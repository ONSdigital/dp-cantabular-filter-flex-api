package api

import (
	"context"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	validatorv10 "github.com/go-playground/validator/v10"

	"github.com/ONSdigital/dp-api-clients-go/v2/identity"
	"github.com/ONSdigital/dp-authorisation/auth"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/middleware"
	dphandlers "github.com/ONSdigital/dp-net/v2/handlers"

	"github.com/go-chi/chi/v5"
)

// API provides a struct to wrap the api around
type API struct {
	Router         chi.Router
	store          datastore
	respond        responder
	generate       generator
	identityClient *identity.Client
	datasets       datasetAPIClient
	ctblr          cantabularClient
	cfg            *config.Config
	validate       *validatorv10.Validate
}

// New creates and initialises a new API
func New(ctx context.Context, cfg *config.Config, r chi.Router, idc *identity.Client, rsp responder, g generator, d datastore, ds datasetAPIClient, c cantabularClient) *API {
	api := &API{
		Router:         r,
		respond:        rsp,
		generate:       g,
		store:          d,
		cfg:            cfg,
		identityClient: idc,
		datasets:       ds,
		ctblr:          c,
		validate:       validatorv10.New(),
	}

	if cfg.EnablePrivateEndpoints {
		api.enablePrivateEndpoints()
	} else {
		api.enablePublicEndpoints()
	}

	return api
}

func (api *API) enablePublicEndpoints() {
	api.Router.Post("/filters", api.createFilter)
	api.Router.Get("/filters/{id}/dimensions", api.getFilterDimensions)
	api.Router.Post("/filter-outputs", api.CreateFilterOutput)
}

func (api *API) enablePrivateEndpoints() {
	r := chi.NewRouter()

	permissions := middleware.NewPermissions(api.cfg.ZebedeeURL, api.cfg.EnablePermissionsAuth)
	checkIdentity := dphandlers.IdentityWithHTTPClient(api.identityClient)

	r.Use(checkIdentity)
	r.Use(middleware.LogIdentity())
	r.Use(permissions.Require(auth.Permissions{Read: true}))

	r.Post("/filters", api.createFilter)
	r.Get("/filters/{id}/dimensions", api.getFilterDimensions)
	r.Post("/filter-outputs", api.CreateFilterOutput)

	api.Router.Mount("/", r)
}
