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
	baseRouter     router
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

type router func() chi.Router

// New creates and initialises a new API
func New(_ context.Context, cfg *config.Config, baseRouter router, idc *identity.Client, rsp responder, g generator, d datastore, ds datasetAPIClient, c cantabularClient, p kafka.IProducer) *API {
	api := &API{
		baseRouter:     baseRouter,
		respond:        rsp,
		generate:       g,
		store:          d,
		cfg:            cfg,
		identityClient: idc,
		datasets:       ds,
		ctblr:          c,
		producer:       p,
	}

	api.initRouter()

	return api
}

func (api *API) initRouter() {
	if api.cfg.EnablePrivateEndpoints {
		api.enablePrivateEndpoints()
	} else {
		api.enablePublicEndpoints()
	}
}

func (api *API) enablePublicEndpoints() {
	api.Router = api.baseRouter()

	api.Router.Post("/filters", api.createFilter)
	api.Router.Get("/filters/{id}", api.getFilter)
	api.Router.Put("/filters/{id}", api.putFilter)
	api.Router.Post("/filters/{id}/submit", api.submitFilter)

	api.Router.Get("/filters/{id}/dimensions", api.getFilterDimensions)
	api.Router.Get("/filters/{id}/dimensions/{dimension}", api.getFilterDimension)
	api.Router.Get("/filters/{id}/dimensions/{name}/options", api.getFilterDimensionOptions)
	api.Router.Post("/filters/{id}/dimensions", api.addFilterDimension)
	api.Router.Put("/filters/{id}/dimensions/{name}", api.updateFilterDimension)

	api.Router.Post("/filters/{id}/dimensions/{dimension}/options/{option}", api.addFilterDimensionOption)

	api.Router.Get("/filter-outputs/{filter-output-id}", api.getFilterOutput)

	api.Router.Get("/datasets/{dataset_id}/editions/{edition}/versions/{version}/json", api.getDatasetJSON)
}

func (api *API) enablePrivateEndpoints() {
	api.Router = chi.NewRouter()
	api.Router.Use(api.baseRouter().Middlewares()...)

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
	api.Router.Put("/filters/{id}/dimensions/{name}", api.updateFilterDimension)
	api.Router.Get("/filters/{id}/dimensions/{dimension}", api.getFilterDimension)
	api.Router.Get("/filters/{id}/dimensions/{name}/options", api.getFilterDimensionOptions)
	api.Router.Post("/filters/{id}/dimensions/{dimension}/options/{option}", api.addFilterDimensionOption)

	api.Router.Get("/filter-outputs/{filter-output-id}", api.getFilterOutput)
	api.Router.Put("/filter-outputs/{filter_output_id}", api.putFilterOutput)
	api.Router.Post("/filter-outputs/{filter_output_id}/events", api.addFilterOutputEvent)

	for _, p := range api.baseRouter().Routes() {
		for m, h := range p.Handlers {
			if m == "*" {
				api.Router.Handle(p.Pattern, h)
				continue
			}
			api.Router.Method(m, p.Pattern, h)
		}
	}
}

// Reset is intended for testing purposes only, and should not be regarded as part of the standard public api of the package
func (api *API) Reset() {
	var err error

	api.cfg, err = config.Get()
	if err != nil {
		panic("error on api reset: " + err.Error())
	}

	api.initRouter()
}
