package responder_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/responder"


	. "github.com/smartystreets/goconvey/convey"
)

type testResponse struct {
	Message string `json:"message"`
}

type testError struct {
	err  error
	msg  string
	code int
}

func (e testError) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e testError) Code() int {
	return e.code
}

func (e testError) Message() string {
	return e.msg
}

func TestJSON(t *testing.T) {
	r := responder.New()

	Convey("Given a valid context and response writer", t, func() {
		ctx := context.Background()
		w := httptest.NewRecorder()

		Convey("Given a valid JSON response object", func() {
			resp := testResponse{
				Message: "Hello, World!",
			}

			Convey("when RespondJSON is called with a given statusCode", func() {
				statusCode := http.StatusCreated

				r.JSON(ctx, w, statusCode, resp)

				Convey("the response writer should record the appropriate status code and response body", func() {
					expectedCode := http.StatusCreated
					expectedBody := `{"message":"Hello, World!"}`

					So(w.Code, ShouldEqual, expectedCode)
					So(w.Body.String(), ShouldResemble, expectedBody)
				})

			})
		})

		Convey("Given an invalid JSON response object", func() {
			resp := make(chan int, 3)

			Convey("when RespondJSON is called with a given statusCode", func() {
				statusCode := http.StatusCreated

				r.JSON(ctx, w, statusCode, resp)

				Convey("the response writer should record an error status code and response body", func() {
					expectedCode := http.StatusInternalServerError
					expectedBody := `{"errors":["Internal Server Error: Badly formed reponse attempt"]}`

					So(w.Code, ShouldEqual, expectedCode)
					So(w.Body.String(), ShouldResemble, expectedBody)
				})
			})
		})
	})
}



func TestError(t *testing.T) {

	r := responder.New()

	Convey("Given a valid context and response writer", t, func() {
		ctx := context.Background()
		w := httptest.NewRecorder()

		Convey("Given a standard Go error", func() {
			err := errors.New("test error")

			Convey("when Error() is called", func() {
				r.Error(ctx, w, err)

				Convey("the response writer should record status code 500 and appropriate error response body", func() {
					expectedCode := http.StatusInternalServerError
					expectedBody := `{"errors":["test error"]}`

					So(w.Code, ShouldEqual, expectedCode)
					So(w.Body.String(), ShouldResemble, expectedBody)
				})

			})
		})

		Convey("Given an error that satisfies interfaces providing Code() and Response() functions", func() {
			err := testError{
				err:  errors.New("test error"),
				msg: "test response",
				code: http.StatusUnauthorized,
			}

			Convey("when Error() is called", func() {
				r.Error(ctx, w, err)

				Convey("the response writer should record the appropriate status code and response message", func() {
					expectedCode := http.StatusUnauthorized
					expectedBody := `{"errors":["test response"]}`

					So(w.Code, ShouldEqual, expectedCode)
					So(w.Body.String(), ShouldResemble, expectedBody)
				})
			})
		})
	})
}

func TestErrors(t *testing.T) {

	r := responder.New()

	Convey("Given a valid context and response writer", t, func() {
		ctx := context.Background()
		w := httptest.NewRecorder()

		Convey("Given a slice of standard Go errors", func() {
			errs := []error{
				errors.New("test error 1"),
				errors.New("test error 2"),
				errors.New("test error 3"),
			}

			Convey("when Errors() is called with given status code and errors", func() {
				r.Errors(ctx, w, http.StatusUnauthorized, errs)

				Convey("the response writer should record appropriate status code and error response body", func() {
					expectedCode := http.StatusUnauthorized
					expectedBody := `{"errors":["test error 1","test error 2","test error 3"]}`

					So(w.Code, ShouldEqual, expectedCode)
					So(w.Body.String(), ShouldResemble, expectedBody)
				})

			})
		})

		Convey("Given an slice errors that satisfy an interface providing a Response() function", func() {
			testErrs := []testError{
				{
					err:  errors.New("test error 1"),
					msg: "test response 1",
				},
				{
					err:  errors.New("test error 2"),
					msg: "test response 2",
				},
				{
					err:  errors.New("test error 3"),
					msg: "test response 3",
				},
			}

			var errs []error
			for _, err := range testErrs{
				errs = append(errs, err)
			}

			Convey("when Errors() is called", func() {
				r.Errors(ctx, w, http.StatusForbidden, errs)

				Convey("the response writer should record the appropriate status code and response message", func() {
					expectedCode := http.StatusForbidden
					expectedBody := `{"errors":["test response 1","test response 2","test response 3"]}`

					So(w.Code, ShouldEqual, expectedCode)
					So(w.Body.String(), ShouldResemble, expectedBody)
				})
			})
		})
	})
}
