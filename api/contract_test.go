package api

import (
	"testing"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	trueVal  = true
	falseVal = false
)

func TestCreateFiltersRequestValid(t *testing.T) {
	Convey("Given a valid createFilterRequest request object", t, func() {
		req := createFilterRequest{
			Dataset: &model.Dataset{
				ID:              "test-dataset-id",
				Edition:         "test-edition",
				Version:         1,
				LowestGeography: "lowest-geography",
			},
			PopulationType: "test-blob",
			Dimensions: []model.Dimension{
				{
					Name:       "test-dimension-1",
					Options:    []string{"a", "b", "c"},
					IsAreaType: &trueVal,
				},
				{
					Name:       "test-dimension-2",
					Options:    []string{"1", "2", "3"},
					IsAreaType: &falseVal,
				},
			},
		}

		Convey("When Valid() is called", func() {
			err := req.Valid()
			So(err, ShouldBeNil)
		})

		Convey("Given datasetID is missing", func() {
			r := req
			r.Dataset.ID = ""
			Convey("When Valid() is called", func() {
				err := r.Valid()
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Given version is 0/missing", func() {
			r := req
			r.Dataset.Version = 0
			Convey("When Valid() is called", func() {
				err := r.Valid()
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Given edition is missing", func() {
			r := req
			r.Dataset.Edition = ""
			Convey("When Valid() is called", func() {
				err := r.Valid()
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Given a dimension is invalid", func() {
			r := req
			r.Dimensions[0] = model.Dimension{}
			Convey("When Valid() is called", func() {
				err := r.Valid()
				So(err, ShouldNotBeNil)
			})
		})
	})
}
