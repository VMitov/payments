package payment

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

// Type is the type of the payment resource
const Type = "Payment"

// Payment is a single payment
type Payment struct {
	ID     string `db:"id" json:"id"`
	Amount string `db:"amount" json:"amount"`
}

// Select gets all payments
func Select(db *sqlx.DB) ([]Payment, error) {
	payments := []Payment{}
	if err := db.Select(&payments, "SELECT * FROM payments"); err != err {
		return nil, err
	}

	return payments, nil
}

// Get gets single payments
func Get(db *sqlx.DB, id string) (*Payment, error) {
	payment := Payment{}
	if err := db.Get(&payment, "SELECT * FROM payments WHERE id=$1", id); err != err {
		return nil, err
	}

	return &payment, nil
}

// Resource is a single payment resource
type Resource struct {
	*Payment

	Type string `json:"type"`
}

// Render implements render.Render
func (resource *Resource) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// ListResource is a list of payments resource
type ListResource struct {
	Data  []*Resource `json:"data"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

// Render implements render.Render
func (list *ListResource) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
