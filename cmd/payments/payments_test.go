package main

import (
	"database/sql"
	"encoding/json"
	"flag"
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
		then      string
		thenF     func(db *sqlx.DB, resp *httptest.ResponseRecorder)
	}{
		"GETNone": {
			given: "Given a HTTP request for /payments",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("GET", "/payments", nil)
			},
			then: "Then the response should be a 200 and no payments in the payload",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 200)
				So(resp.Body.String(), ShouldEqual, `{"data":[],"links":{"self":"/payments"}}`+"\n")
				So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
			},
		},
		"GETList": {
			given: "Given a HTTP request for /payments",
			givenFInt: func(db *sqlx.DB) {
				db.MustExec(`INSERT INTO payments (id, attributes) VALUES
					('4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43', '{"amount": "100.21"}'),
					('216d4da9-e59a-4cc6-8df3-3da6e7580b77', '{"amount": "100.21"}'),
					('7eb8277a-6c91-45e9-8a03-a27f82aca350', '{"amount": "100.21"}'),
					('97fe60ba-1334-439f-91db-32cc3cde036a', '{"amount": "100.21"}'),
					('ab4bbd28-33c6-4231-9b64-0e96190f59ef', '{"amount": "100.21"}'),
					('7f172f5c-f810-4ebe-b015-cb1fc24c6b66', '{"amount": "100.21"}'),
					('502758ff-505f-4d81-b9d2-83aa9c01ebe2', '{"amount": "100.21"}'),
					('09fe827a-b3c2-4437-b999-6c0e780c0983', '{"amount": "100.21"}'),
					('de1f6882-4dba-485a-a632-a80f59fbe4a6', '{"amount": "100.21"}'),
					('b71afd98-4fba-40a4-b8f3-087d005187e3', '{"amount": "100.21"}'),
					('dbb89036-4007-47ff-8fab-00bdd5cc4021', '{"amount": "100.21"}'),
					('52611302-0758-4f69-aa15-c5f55ab7c3eb', '{"amount": "100.21"}'),
					('6cd862ab-6d40-4a86-8037-77d446b3f6fc', '{"amount": "100.21"}'),
					('09a8fe0d-e239-4aff-8098-7923eadd0b98', '{"amount": "100.21"}');
				`)
			},
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "attributes"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", []byte(`{"amount": "100.21"}`)).
					AddRow("216d4da9-e59a-4cc6-8df3-3da6e7580b77", []byte(`{"amount": "100.21"}`)).
					AddRow("7eb8277a-6c91-45e9-8a03-a27f82aca350", []byte(`{"amount": "100.21"}`)).
					AddRow("97fe60ba-1334-439f-91db-32cc3cde036a", []byte(`{"amount": "100.21"}`)).
					AddRow("ab4bbd28-33c6-4231-9b64-0e96190f59ef", []byte(`{"amount": "100.21"}`)).
					AddRow("7f172f5c-f810-4ebe-b015-cb1fc24c6b66", []byte(`{"amount": "100.21"}`)).
					AddRow("502758ff-505f-4d81-b9d2-83aa9c01ebe2", []byte(`{"amount": "100.21"}`)).
					AddRow("09fe827a-b3c2-4437-b999-6c0e780c0983", []byte(`{"amount": "100.21"}`)).
					AddRow("de1f6882-4dba-485a-a632-a80f59fbe4a6", []byte(`{"amount": "100.21"}`)).
					AddRow("b71afd98-4fba-40a4-b8f3-087d005187e3", []byte(`{"amount": "100.21"}`)).
					AddRow("dbb89036-4007-47ff-8fab-00bdd5cc4021", []byte(`{"amount": "100.21"}`)).
					AddRow("52611302-0758-4f69-aa15-c5f55ab7c3eb", []byte(`{"amount": "100.21"}`)).
					AddRow("6cd862ab-6d40-4a86-8037-77d446b3f6fc", []byte(`{"amount": "100.21"}`)).
					AddRow("09a8fe0d-e239-4aff-8098-7923eadd0b98", []byte(`{"amount": "100.21"}`)))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("GET", "/payments", nil)
			},
			then: "Then the response should be a 200 and the payload contains all payments",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 200)

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
				So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
			},
		},
		"GETOne": {
			given: "Given a HTTP request for /payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			givenFInt: func(db *sqlx.DB) {
				db.MustExec(`INSERT INTO payments(id, attributes) VALUES ('4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43', '{"amount": "100.21"}')`)
			},
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "attributes"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", []byte(`{"amount": "100.21"}`)))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("GET", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)
			},
			then: "Then the response should be a 200 and the payload should be the payment with the given id",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 200)
				So(strings.TrimRight(resp.Body.String(), "\n"), ShouldEqual,
					`{"data":{"id":"4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43","attributes":{"amount":"100.21"},"type":"Payment","links":{"self":"/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"}}}`)
				So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
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
			then: "Then the response should be a 404",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 404)
				So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
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
			then: "Then the response should be a 404",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 404)
				So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
			},
		},
		"Create": {
			given: "Given a HTTP request for POST:/payments",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO payments").
					WithArgs(`{"amount": "100.21"}`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"))

				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "attributes"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", []byte(`{"amount": "100.21"}`)))
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "attributes"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", []byte(`{"amount": "100.21"}`)))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("POST", "/payments", strings.NewReader(`{"data": {"type": "Payment", "attributes": {"amount": "100.21"}}}`))
			},
			then: "Then the response should be a 201 and the payload should be as the one from the request",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 201)

				pay := &payment.Resource{}
				if err := json.Unmarshal(resp.Body.Bytes(), pay); err != nil {
					t.Fatal(err)
				}

				if pay.Data == nil {
					So(pay.Data, ShouldNotBeNil)
				}
				req := httptest.NewRequest("GET", "/payments/"+pay.Data.ID, nil)
				getResp := httptest.NewRecorder()
				newRouter(&api{db: db}).ServeHTTP(getResp, req)

				So(strings.TrimRight(getResp.Body.String(), "\n"), ShouldEqual, `{"data":{"id":"`+pay.Data.ID+`","attributes":{"amount":"100.21"},"type":"Payment","links":{"self":"/payments/`+pay.Data.ID+`"}}}`)

				So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
			},
		},
		"CreateNoAttributes": {
			given: "Given a HTTP request for POST:/payments",
			getReq: func() *http.Request {
				return httptest.NewRequest("POST", "/payments", strings.NewReader(`{"data": {"type": "Payment"}}`))
			},
			then: "Then the response should be a 400",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 400)
			},
		},
		"CreateNoData": {
			given: "Given a HTTP request for POST:/payments",
			getReq: func() *http.Request {
				return httptest.NewRequest("POST", "/payments", strings.NewReader(`{"data": {}}`))
			},
			then: "Then the response should be a 400",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 400)
				So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
			},
		},
		"CreateNoPayload": {
			given: "Given a HTTP request for POST:/payments",
			getReq: func() *http.Request {
				return httptest.NewRequest("POST", "/payments", strings.NewReader(`{}`))
			},
			then: "Then the response should be a 400",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 400)
			},
		},
		"CreateNilPayload": {
			given: "Given a HTTP request for POST:/payments",
			getReq: func() *http.Request {
				return httptest.NewRequest("POST", "/payments", nil)
			},
			then: "Then the response should be a 400",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 400)
			},
		},
		"Update": {
			given: "Given a HTTP request for PUT:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			givenFInt: func(db *sqlx.DB) {
				db.MustExec(`INSERT INTO payments (id, attributes) VALUES ('4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43', '{"amount":"100.21"}')`)
			},
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "attributes"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", []byte(`{"amount": "100.21"}`)))

				mock.ExpectExec("UPDATE payments").
					WithArgs(`{"amount": "100.22"}`, "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "attributes"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", []byte(`{"amount": "100.22"}`)))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("PUT", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", strings.NewReader(`{"data":{"type": "Payment", "attributes": {"amount": "100.22"}}}`))
			},
			then: "Then the response should be a 200 and the payload should be updated",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 200)
				So(strings.TrimRight(resp.Body.String(), "\n"), ShouldEqual, `{"data":{"id":"4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43","attributes":{"amount":"100.22"},"type":"Payment","links":{"self":"/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"}}}`)
				So(resp.HeaderMap["Content-Type"], ShouldContain, "application/json; charset=utf-8")
			},
		},
		"UpdateNoData": {
			given: "Given a HTTP request for PUT:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			getReq: func() *http.Request {
				return httptest.NewRequest("PUT", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", strings.NewReader(`{"data":{}}`))
			},
			then: "Then the response should be a 400",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 400)
			},
		},
		"UpdateNoPayload": {
			given: "Given a HTTP request for PUT:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			getReq: func() *http.Request {
				return httptest.NewRequest("PUT", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", strings.NewReader(`{}`))
			},
			then: "Then the response should be a 400",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 400)
			},
		},
		"UpdateMissing": {
			given: "Given a HTTP request for PUT:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("PUT", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", strings.NewReader(`{"data":{"type": "Payment", "attributes": {"amount": "100.21"}}}`))
			},
			then: "Then the response should be a 404",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 404)
			},
		},
		"UpdateBadUUID": {
			given: "Given a HTTP request for PUT:/payments/bad-uuid",
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("PUT", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", strings.NewReader(`{"data":{"type": "Payment", "attributes": {"amount": "100.22"}}}`))
			},
			then: "Then the response should be a 404",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 404)
			},
		},
		"Delete": {
			given: "Given a HTTP request to DELETE:/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			givenFInt: func(db *sqlx.DB) {
				db.MustExec(`INSERT INTO payments (id, attributes) VALUES ('4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43', '{"amount":"100.21"}')`)
			},
			givenF: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "attributes"}).
					AddRow("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", []byte(`{"amount":"100.21"}`)))

				mock.ExpectExec("DELETE FROM payments").
					WithArgs("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			getReq: func() *http.Request {
				return httptest.NewRequest("DELETE", "/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", nil)
			},
			then: "Then the response should be a 200",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 200)
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
			then: "Then the response should be a 404",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 404)
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
			then: "Then the response should be a 404",
			thenF: func(db *sqlx.DB, resp *httptest.ResponseRecorder) {
				So(resp.Code, ShouldEqual, 404)
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

			Convey(tc.given, t, func() {

				Convey("When the request is handled by the Router", func() {
					if integration && tc.givenFInt != nil {
						tc.givenFInt(db)
					}
					if !integration && tc.givenF != nil {
						tc.givenF(mock)
					}

					req := tc.getReq()
					resp := httptest.NewRecorder()
					newRouter(&api{db: db}).ServeHTTP(resp, req)
					Convey(tc.then, func() { tc.thenF(db, resp) })

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
			})

		})
	}
}
