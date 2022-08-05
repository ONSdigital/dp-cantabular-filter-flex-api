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
func New(_ context.Context, cfg *config.Config, r chi.Router, idc *identity.Client, rsp responder, g generator, d datastore, ds datasetAPIClient, c cantabularClient, p kafka.IProducer) *API {
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
	api.Router.Get("/filters/{id}/dimensions/{dimension}/options", api.getFilterDimensionOptions)
	api.Router.Delete("/filters/{id}/dimensions/{dimension}/options", api.deleteFilterDimensionOptions)
	api.Router.Get("/filters/{id}/dimensions/{dimension}", api.getFilterDimension)
	api.Router.Put("/filters/{id}/dimensions/{dimension}", api.updateFilterDimension)
	api.Router.Get("/filters/{id}/dimensions/{dimension}", api.getFilterDimension)

	api.Router.Get("/filters/{id}/dimensions/{dimension}/options", api.getFilterDimensionOptions)
	api.Router.Post("/filters/{id}/dimensions/{dimension}/options/{option}", api.addFilterDimensionOption)
	api.Router.Delete("/filters/{id}/dimensions/{dimension}/options/{option}", api.deleteFilterDimensionOption)

	api.Router.Get("/filter-outputs/{id}", api.getFilterOutput)

	api.Router.Get("/datasets/{dataset_id}/editions/{edition}/versions/{version}/json", api.getDatasetJSON)
}

func (api *API) enablePrivateEndpoints() {
	// This is a hack to work around the issue whereby a chi.Router cannot set middleware after a route has been added
	// Make a copy of the router that has been provided to the constructor of the api, and copy its routes after any
	// necessary middleware has been added
	cp := api.Router
	api.Router = chi.NewRouter()
	api.Router.Use(cp.Middlewares()...)

	api.Router.Use(
		dphandlers.IdentityWithHTTPClient(api.identityClient),
		middleware.LogIdentity(),
		middleware.NewPermissions(api.cfg.ZebedeeURL, api.cfg.EnablePermissionsAuth).Require(auth.Permissions{Read: true}))

	api.Router.Post("/filters", api.createFilter)
	api.Router.Get("/filters/{id}", api.getFilter)
	api.Router.Put("/filters/{id}", api.putFilter)
	api.Router.Post("/filters/{id}/submit", api.submitFilter)

	api.Router.Get("/filters/{id}/dimensions", api.getFilterDimensions)
	api.Router.Post("/filters/{id}/dimensions", api.addFilterDimension)
	api.Router.Get("/filters/{id}/dimensions/{dimension}", api.getFilterDimension)
	api.Router.Put("/filters/{id}/dimensions/{dimension}", api.updateFilterDimension)

	api.Router.Get("/filters/{id}/dimensions/{dimension}/options", api.getFilterDimensionOptions)
	api.Router.Delete("/filters/{id}/dimensions/{dimension}/options", api.deleteFilterDimensionOptions)
	api.Router.Post("/filters/{id}/dimensions/{dimension}/options/{option}", api.addFilterDimensionOption)
	api.Router.Delete("/filters/{id}/dimensions/{dimension}/options/{option}", api.deleteFilterDimensionOption)

	api.Router.Get("/filter-outputs/{id}", api.getFilterOutput)
	api.Router.Put("/filter-outputs/{id}", api.putFilterOutput)
	api.Router.Post("/filter-outputs/{id}/events", api.addFilterOutputEvent)

	for _, p := range cp.Routes() {
		for m, h := range p.Handlers {
			if m == "*" {
				api.Router.Handle(p.Pattern, h)
				continue
			}
			api.Router.Method(m, p.Pattern, h)
		}
	}
}
