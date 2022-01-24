package api

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
)

// responder handles responding to http requests
type responder interface{
	JSON(context.Context, http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, error)
}

type datastore interface{
	CreateFilter(context.Context, *model.Filter) error
}
