package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
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

var (
	integration bool
	dbconn      string
)

func init() {
	flag.BoolVar(&integration, "integration", false, "Run integration tests")
	flag.StringVar(&dbconn, "db", "postgres://postgres@localhost:5432/payments?sslmode=disable", "postgres://user:pass@address:port/db")
}

func TestPayments(t *testing.T) {
	testCases := map[string]struct {
		given     string
		givenFInt func(db *sqlx.DB)          // given func integration
		givenF    func(mock sqlmock.Sqlmock) // given func with mock
		when      string
		getReq    func() *http.Request
		thens     map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder)
	}{
		"GETNone": {
			given: "Given a HTTP request for /payments",
			getReq: func() *http.Request {
				return httptest.NewRequest("GET", "/payments", nil)
			},
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 200": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 200)
				},
				"And the payload should have no payments": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Body.String(), ShouldEqual, `{"data":[],"links":{"self":"/payments"}}`+"\n")
				},
				"The Content-Type should be json": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				},
			},
		},
		"GETList": {
			given: "Given a HTTP request for /payments",
			givenFInt: func(db *sqlx.DB) {
				db.MustExec(`INSERT INTO payments (id, amount) VALUES
					('4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43', 10021),
					('216d4da9-e59a-4cc6-8df3-3da6e7580b77', 10021),
					('7eb8277a-6c91-45e9-8a03-a27f82aca350', 10021),
					('97fe60ba-1334-439f-91db-32cc3cde036a', 10021),
					('ab4bbd28-33c6-4231-9b64-0e96190f59ef', 10021),
					('7f172f5c-f810-4ebe-b015-cb1fc24c6b66', 10021),
					('502758ff-505f-4d81-b9d2-83aa9c01ebe2', 10021),
					('09fe827a-b3c2-4437-b999-6c0e780c0983', 10021),
					('de1f6882-4dba-485a-a632-a80f59fbe4a6', 10021),
					('b71afd98-4fba-40a4-b8f3-087d005187e3', 10021),
					('dbb89036-4007-47ff-8fab-00bdd5cc4021', 10021),
					('52611302-0758-4f69-aa15-c5f55ab7c3eb', 10021),
					('6cd862ab-6d40-4a86-8037-77d446b3f6fc', 10021),
					('09a8fe0d-e239-4aff-8098-7923eadd0b98', 10021);
				`)
			},
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
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 200": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 200)
				},
				"Then the payload should contain all payments": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
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
				"The Content-Type should be json": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				},
			},
		},
		"GETOne": {
			given: "Given a HTTP request for /payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			givenFInt: func(db *sqlx.DB) {
				fmt.Println("INSERT")
				db.MustExec(`INSERT INTO payments(id, amount) VALUES ('4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43', 10021)`)
				fmt.Println("INSERT DONE")
			},
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", 10021))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("GET", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)
			},
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 200": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 200)
				},
				"And the payload should be the payment with the id from the request": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(strings.TrimRight(resp.Body.String(), "\n"), ShouldEqual,
						`{"data":{"id":"4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43","amount":"100.21","type":"Payment","links":{"self":"/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"}}}`)
				},
				"The Content-Type should be json": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
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
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 404": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 404)
				},
				"The Content-Type should be json": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
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
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 404": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 404)
				},
				"The Content-Type should be json": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
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
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 201": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 201)
				},
				"And the payload should be as the one of the request": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					pay := &payment.Resource{}
					if err := json.Unmarshal(resp.Body.Bytes(), pay); err != nil {
						t.Fatal(err)
					}

					req := httptest.NewRequest("GET", "/payments/"+pay.Data.ID, nil)
					getResp := httptest.NewRecorder()
					newRouter(&api{db: db}).ServeHTTP(getResp, req)

					So(strings.TrimRight(resp.Body.String(), "\n"), ShouldEqual, `{"data":{"id":"`+pay.Data.ID+`","amount":"100.21","type":"Payment","links":{"self":"/payments/`+pay.Data.ID+`"}}}`)
				},
				"The Content-Type should be json": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				},
			},
		},
		"CreateNegativeAmount": {
			given: "Given a HTTP request for POST:/payments with negative amount",
			getReq: func() *http.Request {
				return httptest.NewRequest("POST", "/payments", strings.NewReader(`{"data": {"amount": "-100.21", "type": "Payment"}}`))
			},
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 400": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 400)
				},
			},
		},
		"CreateNoData": {
			given: "Given a HTTP request for POST:/payments",
			getReq: func() *http.Request {
				return httptest.NewRequest("POST", "/payments", strings.NewReader(`{"data": {}}`))
			},
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 400": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 400)
				},
				"The Content-Type should be json": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				},
			},
		},
		"CreateNoPayload": {
			given: "Given a HTTP request for POST:/payments",
			getReq: func() *http.Request {
				return httptest.NewRequest("POST", "/payments", strings.NewReader(`{}`))
			},
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 400": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 400)
				},
			},
		},
		"Update": {
			given: "Given a HTTP request for PUT:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			givenFInt: func(db *sqlx.DB) {
				db.MustExec("INSERT INTO payments (id, amount) VALUES ('4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43', 10021)")
			},
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", 10021))

				mock.ExpectExec("UPDATE payments").
					WithArgs(10022, "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", 10022))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("PUT", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", strings.NewReader(`{"data":{"amount": "100.22", "type": "Payment"}}`))
			},
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 200": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 200)
				},
				"And the payload should be the updated payment": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(strings.TrimRight(resp.Body.String(), "\n"), ShouldEqual, `{"data":{"id":"4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43","amount":"100.22","type":"Payment","links":{"self":"/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"}}}`)
				},
				"The Content-Type should be json": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
				},
			},
		},
		"UpdateNoData": {
			given: "Given a HTTP request for PUT:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			getReq: func() *http.Request {
				return httptest.NewRequest("PUT", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", strings.NewReader(`{"data":{}}`))
			},
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 400": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 400)
				},
			},
		},
		"UpdateNoPayload": {
			given: "Given a HTTP request for PUT:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			getReq: func() *http.Request {
				return httptest.NewRequest("PUT", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", strings.NewReader(`{}`))
			},
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 400": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
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
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 404": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
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
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 404": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 404)
				},
			},
		},
		"Delete": {
			given: "Given a HTTP request to DELETE:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			givenFInt: func(db *sqlx.DB) {
				db.MustExec("INSERT INTO payments (id, amount) VALUES ('4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43', 10021)")
			},
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", 10021))

				mock.ExpectExec("DELETE FROM payments").
					WithArgs("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("DELETE", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)
			},
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 200": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
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
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 404": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
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
			thens: map[string]func(db *sqlx.DB, resp *httptest.ResponseRecorder){
				"Then the response should be a 404": func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
					So(resp.Code, ShouldEqual, 404)
				},
			},
		},
	}

	for name, tc := range testCases {
		if integration {
			name = name + "Integration"
		}

		t.Run(name, func(t *testing.T) {
			var (
				db   *sqlx.DB
				mock sqlmock.Sqlmock
			)

			if integration {
				var err error
				db, err = sqlx.Connect("postgres", dbconn)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				var (
					mockDB *sql.DB
					err    error
				)
				mockDB, mock, err = sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				defer mockDB.Close()
				db = sqlx.NewDb(mockDB, "sqlmock")
			}

			if integration && tc.givenFInt != nil {
				tc.givenFInt(db)
			}

			Convey(tc.given, t, func() {
				if !integration && tc.givenF != nil {
					tc.givenF(mock)
				}

				Convey("When the request is handled by the Router", func() {
					req := tc.getReq()
					resp := httptest.NewRecorder()
					newRouter(&api{db: db}).ServeHTTP(resp, req)

					for then, f := range tc.thens {
						Convey(then, func() {
							f(db, resp)
						})
					}
				})
			})

			if !integration {
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}

			if integration {
				// Clean DB
				db.MustExec("DELETE FROM payments")
			}
		})
	}
}
