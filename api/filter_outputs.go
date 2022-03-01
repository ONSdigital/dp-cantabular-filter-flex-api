package api

import (
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	dprequest "github.com/ONSdigital/dp-net/v2/request"
	"github.com/pkg/errors"
)

func (api *API) createFilterOutput(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if !dprequest.IsCallerPresent(ctx) {
		api.respond.Error(
			ctx,
			w,
			http.StatusNotFound,
			Error{
				err:     errors.New("unauthenticated request on unpublished dataset"),
				message: "caller not found",
			},
		)
		return
	}

	var req createFilterOutputsRequest

	if err := api.ParseRequest(r.Body, &req); err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "failed to parse request"),
		)
		return
	}

	err := validateFilterOutput(req.Downloads)
	if err != nil {
		api.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "failed to validate request output"),
		)
		return
	}

	f := model.FilterOutputResponse{
		Download: req.Downloads,
	}
	if err := api.store.CreateFilterOutput(ctx, &f); err != nil {
		api.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.Wrap(err, "failed to create filter"),
		)
		return
	}

	resp := createFilterOutputsResponse{
		ID:        f.ID,
		Downloads: f.Download,
	}

	api.respond.JSON(ctx, w, http.StatusCreated, resp)
}

func isValid(fi *model.FileInfo) error {

	cutset := " "

	if len(strings.Trim(fi.HREF, cutset)) == 0 {
		//	fi.Skipped = true
		return errors.Errorf("HREF is empty in input")
	}

	if len(strings.Trim(fi.Private, cutset)) == 0 {
		//	fi.Skipped = true
		return errors.Errorf("Private is empty in input")
	}

	if len(strings.Trim(fi.Public, cutset)) == 0 {
		//	fi.Skipped = true
		return errors.Errorf("Public is empty in input")
	}

	if len(strings.Trim(fi.Size, cutset)) == 0 {
		//	fi.Skipped = true
		return errors.Errorf("Size is empty in input")
	}

	return nil
}

func validateFilterOutput(filterOutput model.FilterOutput) error {

	//make sure that there are logical values in the input structure
	if err := isValid(filterOutput.CSV); err != nil {
		return err
	}
	if err := isValid(filterOutput.CSVW); err != nil {
		return err
	}
	if err := isValid(filterOutput.TXT); err != nil {
		return err
	}
	if err := isValid(filterOutput.XLS); err != nil {
		return err
	}

	return nil
}
