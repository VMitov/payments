package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VMitov/payments/pkg/payment"
	"github.com/jmoiron/sqlx"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestPayments(t *testing.T) {
	testCases := map[string]struct {
		given  string
		givenF func(mock sqlmock.Sqlmock)
		when   string
		getReq func() *http.Request
		thens  map[string]func(resp *httptest.ResponseRecorder)
	}{
		"GETNone": {
			given: "Given a HTTP request for /payments",
			getReq: func() *http.Request {
				return httptest.NewRequest("GET", "/payments", nil)
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 200": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 200)
				},
				"And the payload should have no payments": func(resp *httptest.ResponseRecorder) {
					So(resp.Body.String(), ShouldEqual, `{"data":[],"links":{"self":"/payments"}}`+"\n")
				},
				"The Content-Type should be json": func(resp *httptest.ResponseRecorder) {
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				},
			},
		},
		"GETList": {
			given: "Given a HTTP request for /payments",
			givenF: func(mock sqlmock.Sqlmock) {
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
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("GET", "/payments", nil)
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 200": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 200)
				},
				"Then the payload should contain all payments": func(resp *httptest.ResponseRecorder) {
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
				},
				"The Content-Type should be json": func(resp *httptest.ResponseRecorder) {
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				},
			},
		},
		"GETOne": {
			given: "Given a HTTP request for /payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", 10021))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("GET", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 200": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 200)
				},
				"And the payload should be the payment with the id from the request": func(resp *httptest.ResponseRecorder) {
					So(strings.TrimRight(resp.Body.String(), "\n"), ShouldEqual,
						`{"data":{"id":"4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43","amount":"100.21","type":"Payment","links":{"self":"/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"}}}`)
				},
				"The Content-Type should be json": func(resp *httptest.ResponseRecorder) {
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				},
			},
		},
		"GETMissing": {
			given: "Given a HTTP request for /payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("GET", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 404": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 404)
				},
				"The Content-Type should be json": func(resp *httptest.ResponseRecorder) {
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				},
			},
		},
		"GETBadUUID": {
			given: "Given a HTTP request for /payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("GET", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 404": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 404)
				},
				"The Content-Type should be json": func(resp *httptest.ResponseRecorder) {
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				},
			},
		},
		"Create": {
			given: "Given a HTTP request for POST:/payments",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO payments").
					WithArgs(10021).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"))

				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", 10021))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("POST", "/payments", strings.NewReader(`{"data": {"amount": "100.21", "type": "Payment"}}`))
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 201": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 201)
				},
				"And the payload should be as the one of the request": func(resp *httptest.ResponseRecorder) {
					So(strings.TrimRight(resp.Body.String(), "\n"), ShouldEqual, `{"data":{"id":"4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43","amount":"100.21","type":"Payment","links":{"self":"/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"}}}`)
				},
				"The Content-Type should be json": func(resp *httptest.ResponseRecorder) {
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				},
			},
		},
		"CreateNoData": {
			given: "Given a HTTP request for POST:/payments",
			getReq: func() *http.Request {
				return httptest.NewRequest("POST", "/payments", strings.NewReader(`{"data": {}}`))
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 400": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 400)
				},
				"The Content-Type should be json": func(resp *httptest.ResponseRecorder) {
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				},
			},
		},
		"CreateNoPayload": {
			given: "Given a HTTP request for POST:/payments",
			getReq: func() *http.Request {
				return httptest.NewRequest("POST", "/payments", strings.NewReader(`{}`))
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 400": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 400)
				},
			},
		},
		"Update": {
			given: "Given a HTTP request for PUT:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", 10022))

				mock.ExpectExec("UPDATE payments").
					WithArgs(10022, "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", 10022))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("PUT", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", strings.NewReader(`{"data":{"amount": "100.22", "type": "Payment"}}`))
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 200": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 200)
				},
				"And the payload should be the updated payment": func(resp *httptest.ResponseRecorder) {
					So(strings.TrimRight(resp.Body.String(), "\n"), ShouldEqual, `{"data":{"id":"4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43","amount":"100.22","type":"Payment","links":{"self":"/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"}}}`)
				},
				"The Content-Type should be json": func(resp *httptest.ResponseRecorder) {
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				},
			},
		},
		"UpdateNoData": {
			given: "Given a HTTP request for PUT:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			getReq: func() *http.Request {
				return httptest.NewRequest("PUT", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", strings.NewReader(`{"data":{}}`))
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 400": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 400)
				},
			},
		},
		"UpdateNoPayload": {
			given: "Given a HTTP request for PUT:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			getReq: func() *http.Request {
				return httptest.NewRequest("PUT", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", strings.NewReader(`{}`))
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 400": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 400)
				},
			},
		},
		"UpdateMissing": {
			given: "Given a HTTP request for PUT:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("PUT", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", strings.NewReader(`{"data":{"amount": "100.22", "type": "Payment"}}`))
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 404": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 404)
				},
			},
		},
		"UpdateBadUUID": {
			given: "Given a HTTP request for PUT:/payments/bad-uuid",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("PUT", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", strings.NewReader(`{"data":{"amount": "100.22", "type": "Payment"}}`))
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 404": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 404)
				},
			},
		},
		"Delete": {
			given: "Given a HTTP request to DELETE:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", 10022))

				mock.ExpectExec("DELETE FROM payments").
					WithArgs("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("DELETE", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 200": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 200)
				},
			},
		},
		"DeleteNonExisting": {
			given: "Given a HTTP request to DELETE:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43 which is not existing",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("DELETE", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 404": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 404)
				},
			},
		},
		"DeleteInvalidUUID": {
			given: "Given a HTTP request to DELETE:/payments/bad-uuid with wrong uuid",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("DELETE", "/payments/bad-uuid", nil)
			},
			thens: map[string]func(resp *httptest.ResponseRecorder){
				"Then the response should be a 404": func(resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 404)
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			sqlxDB := sqlx.NewDb(db, "sqlmock")

			Convey(tc.given, t, func() {
				if tc.givenF != nil {
					tc.givenF(mock)
				}
				Convey("When the request is handled by the Router", func() {
					req := tc.getReq()
					resp := httptest.NewRecorder()
					newRouter(&api{db: sqlxDB}).ServeHTTP(resp, req)

					for then, f := range tc.thens {
						Convey(then, func() {
							f(resp)
						})
					}
				})
			})

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
