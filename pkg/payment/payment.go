package payment

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

// Type is the type of the payment resource
const Type = "Payment"

// Payment is a single payment
type Payment struct {
	ID string `db:"id" json:"id"`
}

// Select gets all payments
func Select(db *sqlx.DB) ([]Payment, error) {
	payments := []Payment{}
	if err := db.Select(&payments, "SELECT * FROM payments"); err != err {
		return nil, err
	}

	return payments, nil
}

// Resource is a single payment resource
type Resource struct {
	*Payment

	Type string `json:"type"`
}

// ListResource is a list of payments resource
type ListResource struct {
	Data  []*Resource `json:"data"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

// Render writes json representation of the ListResource to http.ResponseWriter
func (list *ListResource) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
