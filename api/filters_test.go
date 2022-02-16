package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/api/mock"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidateDimensions(t *testing.T) {
	ctx := context.Background()
	var api API

	Convey("Given a set of filter dimensions which match existing dimensions found in Version doc", t, func() {
		existingDims := []dataset.VersionDimension{
			{
				Name: "foo",
				ID:   "foo01",
			},
			{
				Name: "bar",
				ID:   "bar01",
			},
		}

		filterDims := []model.Dimension{
			{
				Name: "foo",
				Options: []string{
					"foo_1",
					"foo_2",
				},
			},
			{
				Name: "bar",
				Options: []string{
					"bar_1",
					"bar_2",
				},
			},
		}

		Convey("When validateDimensions is called", func() {
			expected := map[string]string{
				"foo": "foo01",
				"bar": "bar01",
			}

			got, err := api.validateDimensions(ctx, filterDims, existingDims)
			So(err, ShouldBeNil)
			So(got, ShouldResemble, expected)
		})
	})

	Convey("Given a set of filter dimensions which include dimensions not found in Version doc", t, func() {
		existingDims := []dataset.VersionDimension{
			{
				Name: "foo",
				ID:   "foo01",
			},
			{
				Name: "bar",
				ID:   "bar01",
			},
		}

		filterDims := []model.Dimension{
			{
				Name: "foo",
				Options: []string{
					"foo_1",
					"foo_2",
				},
			},
			{
				Name: "bar",
				Options: []string{
					"bar_1",
					"bar_2",
				},
			},
			{
				Name: "alice",
				Options: []string{
					"alice_1",
					"alice_2",
				},
			},
		}

		Convey("When validateDimensions is called", func() {
			_, err := api.validateDimensions(ctx, filterDims, existingDims)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestValidateDimensionOptions(t *testing.T) {
	ctx := context.Background()
	var api API

	Convey("Given a Cantabular Client which cannot find provided dimension options", t, func() {
		api.ctblr = &mock.CantabularClient{
			ErrStatus:    http.StatusInternalServerError,
			OptionsHappy: false,
		}

		Convey("When validateDimenOptions is called", func() {
			dimIDs := map[string]string{
				"foo": "foo01",
				"bar": "bar01",
			}

			req := createFilterRequest{
				Dimensions: []model.Dimension{
					{
						Name: "foo",
						Options: []string{
							"foo_1",
							"foo_2",
						},
					},
					{
						Name: "bar",
						Options: []string{
							"bar_1",
							"bar_2",
						},
					},
				},
			}

			err := api.validateDimensionOptions(ctx, req, dimIDs)
			So(err, ShouldNotBeNil)
		})

		Convey("When validateDimenOptions is called but there are no options selected", func() {
			dimIDs := map[string]string{
				"foo": "foo01",
				"bar": "bar01",
			}

			req := createFilterRequest{
				Dimensions: []model.Dimension{
					{
						Name: "foo",
					},
					{
						Name: "bar",
					},
				},
			}

			err := api.validateDimensionOptions(ctx, req, dimIDs)
			So(err, ShouldBeNil)
		})
	})

	Convey("Given a Cantabular Client which can find provided dimension options", t, func() {
		api.ctblr = &mock.CantabularClient{
			ErrStatus:    http.StatusInternalServerError,
			OptionsHappy: true,
		}

		Convey("When validateDimenOptions is called", func() {
			dimIDs := map[string]string{
				"foo": "foo01",
				"bar": "bar01",
			}

			req := createFilterRequest{
				Dimensions: []model.Dimension{
					{
						Name: "foo",
						Options: []string{
							"foo_1",
							"foo_2",
						},
					},
					{
						Name: "bar",
						Options: []string{
							"bar_1",
							"bar_2",
						},
					},
				},
			}

			err := api.validateDimensionOptions(ctx, req, dimIDs)
			So(err, ShouldBeNil)
		})
	})
}
