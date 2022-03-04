package api

import (
	"testing"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateFiltersRequestValid(t *testing.T) {
	Convey("Given a valid createFilterRequest request object", t, func() {
		req := createFilterRequest{
			Dataset: &model.Dataset{
				ID:      "test-dataset-id",
				Edition: "test-edition",
				Version: 1,
			},
			PopulationType: "test-blob",
			Dimensions: []model.Dimension{
				{
					Name:         "test-dimension-1",
					Options:      []string{"a", "b", "c"},
					DimensionURL: "http://dim-1.com",
					IsAreaType:   true,
				},
				{
					Name:         "test-dimension-2",
					Options:      []string{"1", "2", "3"},
					DimensionURL: "http://dim-2.com",
					IsAreaType:   false,
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

		Convey("If there are less than 2 dimensions", func() {
			r := req
			r.Dimensions = r.Dimensions[:1]
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

func TestCreateFilterOutputsRequestValid(t *testing.T) {

	Convey("Given a valid createFilterOutputsRequest request object", t, func() {
		blankInfo := model.FileInfo{
			HREF:    " ",
			Size:    " ",
			Public:  " ",
			Private: " ",
			Skipped: true,
		}

		partialblankInfo := model.FileInfo{
			HREF:    " ",
			Size:    " ",
			Public:  "test1 test ",
			Private: " ",
			Skipped: true,
		}

		validInfo := model.FileInfo{
			HREF:    "test ",
			Size:    "  tets tts t",
			Public:  "test1 test ",
			Private: " t e s t",
			Skipped: true,
		}

		req := createFilterOutputsRequest{
			State: "published",
			Downloads: model.Downloads{
				CSV:  &blankInfo,
				CSVW: new(model.FileInfo),
				TXT:  new(model.FileInfo),
				XLS:  &partialblankInfo,
			},
		}

		Convey("When Valid() is called with invalid input", func() {
			err := req.Valid()
			So(err, ShouldNotBeNil)
		})

		Convey("When Valid() is called with valid input", func() {
			r := req
			r.Downloads.CSV = &validInfo
			r.Downloads.CSVW = &validInfo
			r.Downloads.TXT = &validInfo
			r.Downloads.XLS = &validInfo
			err := r.Valid()
			So(err, ShouldBeNil)
		})
	})

}
