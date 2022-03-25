package api

import (
	"context"
	// "net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	kafka "github.com/ONSdigital/dp-kafka/v3"

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
	Producer       *kafka.Producer
	identityClient *identity.Client
	datasets       datasetAPIClient
	ctblr          cantabularClient
	cfg            *config.Config
}

// New creates and initialises a new API
func New(ctx context.Context, cfg *config.Config, r chi.Router, idc *identity.Client, rsp responder, g generator, d datastore, ds datasetAPIClient, c cantabularClient, producer *kafka.Producer) *API {
	api := &API{
		Router:         r,
		respond:        rsp,
		generate:       g,
		store:          d,
		cfg:            cfg,
		identityClient: idc,
		datasets:       ds,
		ctblr:          c,

		Producer: producer,
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
	api.Router.Get("/filters/{id}", api.getFilter)
	api.Router.Put("/filters/{id}", api.putFilter)
	api.Router.Post("/filters/{id}/submit", api.submitFilter)
	api.Router.Get("/filters/{id}/dimensions", api.getFilterDimensions)
	api.Router.Post("/filters/{id}/dimensions", api.addFilterDimension)
}

func (api *API) enablePrivateEndpoints() {
	r := chi.NewRouter()

	permissions := middleware.NewPermissions(api.cfg.ZebedeeURL, api.cfg.EnablePermissionsAuth)
	checkIdentity := dphandlers.IdentityWithHTTPClient(api.identityClient)

	r.Use(checkIdentity)
	r.Use(middleware.LogIdentity())
	r.Use(permissions.Require(auth.Permissions{Read: true}))

	r.Post("/filters", api.createFilter)
	r.Get("/filters/{id}", api.getFilter)
	r.Put("/filters/{id}", api.putFilter)
	r.Post("/filters/{id}/submit", api.submitFilter)
	r.Get("/filters/{id}/dimensions", api.getFilterDimensions)
	r.Post("/filters/{id}/dimensions", api.addFilterDimension)
	r.Post("/filter-outputs", api.createFilterOutput)
	api.Router.Mount("/", r)
}
