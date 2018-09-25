package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VMitov/payments/pkg/payment"
	"github.com/jmoiron/sqlx"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestPayments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	t.Run("TestPaymentsGETNone", func(t *testing.T) {
		Convey("Given a HTTP request for /payments", t, func() {
			req := httptest.NewRequest("GET", "/payments", nil)
			resp := httptest.NewRecorder()
			Convey("When the request is handled by the Router", func() {
				newRouter(&api{db: sqlxDB}).ServeHTTP(resp, req)
				Convey("Then the response should be a 200", func() {
					So(resp.Code, ShouldEqual, 200)
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
					So(resp.Body.String(), ShouldEqual, `{"data":[],"links":{"self":""}}`+"\n")
				})
			})
		})
	})

	t.Run("TestPaymentsGETList", func(t *testing.T) {
		Convey("Given a HTTP request for /payments", t, func() {
			mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).
				AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", 10021).
				AddRow("216d4da9-e59a-4cc6-8df3-3da6e7580b77", 10021).
				AddRow("7eb8277a-6c91-45e9-8a03-a27f82aca350", 10021).
				AddRow("97fe60ba-1334-439f-91db-32cc3cde036a", 10021).
				AddRow("ab4bbd28-33c6-4231-9b64-0e96190f59ef", 10021).
				AddRow("7f172f5c-f810-4ebe-b015-cb1fc24c6b66", 10021).
				AddRow("502758ff-505f-4d81-b9d2-83aa9c01ebe2", 10021).
				AddRow("09fe827a-b3c2-4437-b999-6c0e780c0983", 10021).
				AddRow("de1f6882-4dba-485a-a632-a80f59fbe4a6", 10021).
				AddRow("b71afd98-4fba-40a4-b8f3-087d005187e3", 10021).
				AddRow("dbb89036-4007-47ff-8fab-00bdd5cc4021", 10021).
				AddRow("52611302-0758-4f69-aa15-c5f55ab7c3eb", 10021).
				AddRow("6cd862ab-6d40-4a86-8037-77d446b3f6fc", 10021).
				AddRow("09a8fe0d-e239-4aff-8098-7923eadd0b98", 10021))

			req := httptest.NewRequest("GET", "/payments", nil)
			resp := httptest.NewRecorder()
			Convey("When the request is handled by the Router", func() {
				newRouter(&api{db: sqlxDB}).ServeHTTP(resp, req)
				Convey("Then the response should be a 200", func() {
					So(resp.Code, ShouldEqual, 200)
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")

					result := &payment.ListResource{}
					if err := json.Unmarshal(resp.Body.Bytes(), result); err != nil {
						t.Fatal(err)
					}
					resource, err := json.MarshalIndent(result, "", "    ")
					if err != nil {
						t.Fatal(err)
					}

					expected, err := ioutil.ReadFile("../../testdata/payments.json")
					if err != nil {
						t.Fatal(err)
					}

					So(string(resource), ShouldEqual, strings.TrimRight(string(expected), "\n"))
				})
			})
		})
	})

	t.Run("TestPaymentsGETOne", func(t *testing.T) {
		Convey("Given a HTTP request for /payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", t, func() {
			mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).
				AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", 10021))
			req := httptest.NewRequest("GET", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)
			resp := httptest.NewRecorder()
			Convey("When the request is handled by the Router", func() {
				newRouter(&api{db: sqlxDB}).ServeHTTP(resp, req)
				Convey("Then the response should be a 200", func() {
					So(resp.Code, ShouldEqual, 200)
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
					So(strings.TrimRight(resp.Body.String(), "\n"), ShouldEqual,
						`{"id":"4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43","amount":"100.21","type":"Payment"}`)
				})
			})
		})
	})

	t.Run("TestPaymentsGETMissing", func(t *testing.T) {
		Convey("Given a HTTP request for /payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", t, func() {
			mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}))
			req := httptest.NewRequest("GET", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)
			resp := httptest.NewRecorder()
			Convey("When the request is handled by the Router", func() {
				newRouter(&api{db: sqlxDB}).ServeHTTP(resp, req)
				Convey("Then the response should be a 200", func() {
					So(resp.Code, ShouldEqual, 404)
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
					So(strings.TrimRight(resp.Body.String(), "\n"), ShouldEqual,
						`{"status":"Resource not found."}`)
				})
			})
		})
	})

	t.Run("TestPaymentsCreate", func(t *testing.T) {
		Convey("Given a HTTP request for POST:/payments", t, func() {
			mock.ExpectQuery("INSERT INTO payments").
				WithArgs(10021).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"))

			mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).
				AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", 10021))

			req := httptest.NewRequest("POST", "/payments", strings.NewReader(`{"amount": "100.21", "type": "Payment"}`))
			resp := httptest.NewRecorder()
			Convey("When the request is handled by the Router", func() {
				newRouter(&api{db: sqlxDB}).ServeHTTP(resp, req)
				Convey("Then the response should be a 201", func() {
					So(strings.TrimRight(resp.Body.String(), "\n"), ShouldEqual, `{"id":"4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43","amount":"100.21","type":"Payment"}`)
					So(resp.Code, ShouldEqual, 201)
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				})
			})
		})
	})

}
