package main

import (
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetPayments(t *testing.T) {
	Convey("Given a HTTP request for /payments", t, func() {
		req := httptest.NewRequest("GET", "/payments", nil)
		resp := httptest.NewRecorder()
		Convey("When the request is handled by the Router", func() {
			newRouter().ServeHTTP(resp, req)
			Convey("Then the response should be a 200", func() {
				So(resp.Code, ShouldEqual, 200)
				So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				So(resp.Body.String(), ShouldEqual, `{"data":[{"id":"1","type":"Payment"}],"links":{"self":""}}`+"\n")
			})
		})
	})
}
