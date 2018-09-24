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

func TestGetPayments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	Convey("Given no payment", t, func() {
		Convey("Given a HTTP request for /payments", func() {
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

	Convey("Given payments", t, func() {
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).
			AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", "100.21").
			AddRow("216d4da9-e59a-4cc6-8df3-3da6e7580b77", "100.21").
			AddRow("7eb8277a-6c91-45e9-8a03-a27f82aca350", "100.21").
			AddRow("97fe60ba-1334-439f-91db-32cc3cde036a", "100.21").
			AddRow("ab4bbd28-33c6-4231-9b64-0e96190f59ef", "100.21").
			AddRow("7f172f5c-f810-4ebe-b015-cb1fc24c6b66", "100.21").
			AddRow("502758ff-505f-4d81-b9d2-83aa9c01ebe2", "100.21").
			AddRow("09fe827a-b3c2-4437-b999-6c0e780c0983", "100.21").
			AddRow("de1f6882-4dba-485a-a632-a80f59fbe4a6", "100.21").
			AddRow("b71afd98-4fba-40a4-b8f3-087d005187e3", "100.21").
			AddRow("dbb89036-4007-47ff-8fab-00bdd5cc4021", "100.21").
			AddRow("52611302-0758-4f69-aa15-c5f55ab7c3eb", "100.21").
			AddRow("6cd862ab-6d40-4a86-8037-77d446b3f6fc", "100.21").
			AddRow("09a8fe0d-e239-4aff-8098-7923eadd0b98", "100.21"))

		Convey("Given a HTTP request for /payments", func() {
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
}
