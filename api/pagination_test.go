package api

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetPaginationParams(t *testing.T) {
	Convey("Given a valid offset and limit, it parses and returns the values", t, func() {
		parsedUrl, err := url.Parse("http://test.test?limit=10&offset=0")
		So(err, ShouldBeNil)

		limit, offset, err := getPaginationParams(parsedUrl, 100)

		Convey("It should return the parsed values", func() {
			So(err, ShouldBeNil)
			So(limit, ShouldEqual, 10)
			So(offset, ShouldEqual, 0)
		})
	})

	Convey("Given an invalid offset or limit", t, func() {
		tests := map[string]struct {
			url          string
			maximumLimit int
		}{
			"Limit cannot be parsed": {
				url:          "http://test.test?limit=dog",
				maximumLimit: 20,
			},
			"Offset cannot be parsed": {
				url:          "http://test.test?offset=dog",
				maximumLimit: 20,
			},
			"Limit exceeds maximum": {
				url:          "http://test.test?limit=21",
				maximumLimit: 20,
			},
			"Limit is negative": {
				url:          "http://test.test?limit=-1",
				maximumLimit: 20,
			},
			"Offset is negative": {
				url:          "http://test.test?offset=-1",
				maximumLimit: 20,
			},
		}

		for desc, test := range tests {
			Convey(desc, func() {
				parsedUrl, err := url.Parse(test.url)
				So(err, ShouldBeNil)

				_, _, err = getPaginationParams(parsedUrl, test.maximumLimit)

				Convey("There should be an error", func() {
					So(err, ShouldNotBeNil)
				})
			})
		}
	})

	Convey("Given no pagination params, it falls back to default values", t, func() {
		parsedUrl, err := url.Parse("http://test.test")
		So(err, ShouldBeNil)

		limit, offset, err := getPaginationParams(parsedUrl, 20)

		So(err, ShouldBeNil)
		So(limit, ShouldEqual, 20)
		So(offset, ShouldEqual, 0)
	})
}
