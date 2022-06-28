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
	producer       kafka.IProducer
	identityClient *identity.Client
	datasets       datasetAPIClient
	ctblr          cantabularClient
	cfg            *config.Config
}

// New creates and initialises a new API
func New(ctx context.Context, cfg *config.Config, r chi.Router, idc *identity.Client, rsp responder, g generator, d datastore, ds datasetAPIClient, c cantabularClient, p kafka.IProducer) *API {
	api := &API{
		Router:         r,
		respond:        rsp,
		generate:       g,
		store:          d,
		cfg:            cfg,
		identityClient: idc,
		datasets:       ds,
		ctblr:          c,
		producer:       p,
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
	api.Router.Get("/filters/{id}/dimensions/{dimension}", api.getFilterDimension)
	api.Router.Put("/filters/{id}/dimensions/{name}", api.updateFilterDimension)
	api.Router.Get("/filters/{id}/dimensions/{dimension}", api.getFilterDimension)

	api.Router.Get("/filters/{id}/dimensions/{name}/options", api.getFilterDimensionOptions)
	api.Router.Post("/filters/{id}/dimensions/{dimension}/options/{option}", api.addFilterDimensionOption)
	api.Router.Delete("/filters/{id}/dimensions/{dimension}/options/{option}", api.deleteFilterDimensionOption)

	api.Router.Get("/filter-outputs/{filter-output-id}", api.getFilterOutput)

	api.Router.Get("/datasets/{dataset_id}/editions/{edition}/versions/{version}/json", api.getDatasetJSON)
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
	r.Get("/filters/{id}/dimensions/{dimension}", api.getFilterDimension)
	r.Post("/filters/{id}/dimensions", api.addFilterDimension)
	r.Put("/filters/{id}/dimensions/{name}", api.updateFilterDimension)

	r.Delete("/filters/{id}/dimensions/{dimension}/options/{option}", api.deleteFilterDimensionOption)
	r.Post("/filters/{id}/dimensions/{dimension}/options/{option}", api.addFilterDimensionOption)
	r.Get("/filters/{id}/dimensions/{name}/options", api.getFilterDimensionOptions)

	r.Get("/filter-outputs/{filter-output-id}", api.getFilterOutput)
	r.Put("/filter-outputs/{filter_output_id}", api.putFilterOutput)
	r.Post("/filter-outputs/{filter_output_id}/events", api.addFilterOutputEvent)

	api.Router.Mount("/", r)
}
