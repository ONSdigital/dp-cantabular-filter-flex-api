package api_test

import (
	"testing"
	"net/http"
	"bytes"
	"errors"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/api"

	. "github.com/smartystreets/goconvey/convey"
)

func statusCode(err error) int {
	var cerr interface{Code() int}
	if errors.As(err, &cerr) {
		return cerr.Code()
	}

	return 0
}

type fooRequest struct{
	Foo     string
	Bar     int
	IsValid bool
}

func (r *fooRequest) Valid() error{
	if !r.IsValid{
		return errors.New("test invalid request message")
	}
	return nil
}

type barRequest struct{
	Alice string
	Bob   int
}

func TestParseRequest(t *testing.T){
	api := &api.API{}

	Convey("Given a valid request body matching request object", t, func() {
		b := []byte(`{"foo":"I am foo", "bar": 2, "isValid": true}`)
		var req fooRequest

		Convey("When ParseRequest(body, request) is called", func() {
			err := api.ParseRequest(bytes.NewReader(b), &req)
			So(err, ShouldBeNil)

			Convey("The request object should be populated with the expected values ", func() {
				expected := fooRequest{
					Foo:   "I am foo",
					Bar:   2,
					IsValid: true,
				}
				So(req, ShouldResemble, expected)
			})
		})
	})

	Convey("Given a request body that fails validation check with matching request object", t, func() {
		b := []byte(`{"foo":"I am foo", "bar": 2, "isValid": false}`)
		var req fooRequest

		Convey("When ParseRequest(body, request) is called", func() {
			err := api.ParseRequest(bytes.NewReader(b), &req)
			So(err, ShouldNotBeNil)
			So(statusCode(err), ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("Given an invalid JSON request body with matching request object", t, func() {
		b := []byte(`//4tjwopjofo}`)
		var req fooRequest

		Convey("When ParseRequest(body, request) is called", func() {
			err := api.ParseRequest(bytes.NewReader(b), &req)
			So(err, ShouldNotBeNil)
			So(statusCode(err), ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("Given a valid request body matching request object that does not have a Valid function", t, func() {
		b := []byte(`{"Alice":"I am Alice", "Bob": 23}`)
		var req barRequest

		Convey("When ParseRequest(body, request) is called", func() {
			err := api.ParseRequest(bytes.NewReader(b), &req)
			So(err, ShouldBeNil)

			Convey("The request object should be populated with the expected values ", func() {
				expected := barRequest{
					Alice:   "I am Alice",
					Bob:   23,
				}
				So(req, ShouldResemble, expected)
			})
		})
	})
}