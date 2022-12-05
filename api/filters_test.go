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

func Test_getFilterVariable(t *testing.T) {
	Convey("Given an array of dimension IDs exists", t, func() {
		Convey("When I have a dimension with no FilterByParent", func() {
			d := model.Dimension{
				Name:           "city",
				FilterByParent: "",
			}

			Convey("Then the 'city' should be returned as the filter variable", func() {
				got := getFilterVariable(d)
				So(got, ShouldResemble, "city")
			})
		})

		Convey("When I have a dimension with a FilterByParent", func() {
			d := model.Dimension{
				Name:           "city",
				FilterByParent: "region",
			}

			Convey("Then 'region' should be returned as the filter variable", func() {
				got := getFilterVariable(d)
				So(got, ShouldResemble, "region")
			})
		})
	})
}

func TestValidateDimensions(t *testing.T) {
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
			err := api.validateDimensionsFromVersion(filterDims, existingDims)
			So(err, ShouldBeNil)
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
			err := api.validateDimensionsFromVersion(filterDims, existingDims)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a set of filter dimensions are duplicate", t, func() {
		existingDims := []dataset.VersionDimension{
			{
				Name: "foo",
				ID:   "foo01",
			},
			{
				Name: "foo",
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
				Name: "foo",
				Options: []string{
					"bar_1",
					"bar_2",
				},
			},
		}

		Convey("When validateDimensions is called", func() {
			err := api.validateDimensionsFromVersion(filterDims, existingDims)
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
				PopulationType: "Example",
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

			err := api.validateDimensionOptionsNew(ctx, req.Dimensions, dimIDs, req.PopulationType)
			So(err, ShouldNotBeNil)
		})

		Convey("When validateDimenOptions is called but there are no options selected", func() {
			dimIDs := map[string]string{
				"foo": "foo01",
				"bar": "bar01",
			}
			req := createFilterRequest{
				PopulationType: "Example",
				Dimensions: []model.Dimension{
					{
						Name: "foo",
					},
					{
						Name: "bar",
					},
				},
			}

			err := api.validateDimensionOptionsNew(ctx, req.Dimensions, dimIDs, req.PopulationType)
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
				PopulationType: "Example",
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

			err := api.validateDimensionOptionsNew(ctx, req.Dimensions, dimIDs, req.PopulationType)
			So(err, ShouldBeNil)
		})
	})

	Convey("Given a Cantabular Client which is expecting the value present in FilterByParent variable", t, func() {
		parentDim := "ParentFoo"

		api.ctblr = &mock.CantabularClient{
			ErrStatus:               http.StatusInternalServerError,
			OptionsHappy:            true,
			ExpectedFilterDimension: parentDim,
		}

		Convey("When validateDimenOptions is called with value as the parent of a variable", func() {
			dimIDs := map[string]string{
				"foo": "foo01",
				"bar": "bar01",
			}
			req := createFilterRequest{
				PopulationType: "Example",
				Dimensions: []model.Dimension{
					{
						Name: "foo",
						Options: []string{
							"foo_1",
							"foo_2",
						},
						FilterByParent: parentDim,
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

			err := api.validateDimensionOptionsNew(ctx, req.Dimensions, dimIDs, req.PopulationType)
			So(err, ShouldBeNil)
		})

		Convey("When validateDimenOptions is called with the correct value but not as the parent", func() {
			dimIDs := map[string]string{
				"foo": "foo01",
				"bar": "bar01",
			}
			req := createFilterRequest{
				PopulationType: "Example",
				Dimensions: []model.Dimension{
					{
						Name: parentDim,
						Options: []string{
							"foo_1",
							"foo_2",
						},
						FilterByParent: "incorrect",
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

			err := api.validateDimensionOptionsNew(ctx, req.Dimensions, dimIDs, req.PopulationType)
			So(err, ShouldNotBeNil)
		})

		Convey("When validateDimenOptions is called without the correct value", func() {
			dimIDs := map[string]string{
				"foo": "foo01",
				"bar": "bar01",
			}
			req := createFilterRequest{
				PopulationType: "Example",
				Dimensions: []model.Dimension{
					{
						Name: "foo",
						Options: []string{
							"foo_1",
							"foo_2",
						},
						FilterByParent: "incorrect",
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

			err := api.validateDimensionOptionsNew(ctx, req.Dimensions, dimIDs, req.PopulationType)
			So(err, ShouldNotBeNil)
		})
	})
}
