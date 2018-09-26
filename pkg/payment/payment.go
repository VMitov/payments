package payment

import (
	"fmt"
	"net/http"

	"github.com/VMitov/payments/pkg/links"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

// Type is the type of the payment resource
const Type = "Payment"

// Payment is a single payment
type Payment struct {
	ID     string `db:"id" json:"id"`
	Amount int64  `db:"amount" json:"-"`
}

// NewFromResource returns Payment from Resource
func NewFromResource(res *Resource) (*Payment, error) {
	d, err := decimal.NewFromString(res.Data.Amount)
	if err != nil {
		return nil, err
	}

	if res.Data.Payment == nil {
		res.Data.Payment = &Payment{}
	}
	res.Data.Payment.Amount = d.Mul(decimal.New(100, 0)).IntPart()
	return res.Data.Payment, nil
}

// Create persist a payment
func Create(db *sqlx.DB, pay *Payment) (id string, err error) {
	rows, err := db.Queryx(
		db.Rebind(`INSERT INTO payments (amount) VALUES (?) RETURNING id`),
		10021,
	)
	if err != nil {
		return "", err
	}

	if !rows.Next() {
		return "", fmt.Errorf("error getting id")
	}

	rows.Scan(&id)
	return id, nil
}

// Update updates payment
func Update(db *sqlx.DB, id string, pay *Payment) error {
	_, err := db.Exec(
		db.Rebind(`UPDATE payments SET amount=? WHERE id=?`),
		pay.Amount, id,
	)

	return err
}

// Delete deleted payment
func Delete(db *sqlx.DB, id string) error {
	_, err := db.Exec(db.Rebind(`DELETE FROM payments WHERE id=?`), id)
	return err
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

// ResourceData is the data of the payment resource
type ResourceData struct {
	*Payment

	Amount string `json:"amount"`
	Type   string `json:"type"`

	links.Resource
}

func newResourceData(p *Payment, self string) *ResourceData {
	return &ResourceData{
		Payment:  p,
		Type:     Type,
		Amount:   fmt.Sprintf("%.2f", float64(p.Amount)/100),
		Resource: links.Resource{Links: links.Links{Self: self}},
	}
}

// Resource is a single payment resource
type Resource struct {
	Data *ResourceData `json:"data"`
}

// NewResource create new resource from Payment
func NewResource(p *Payment, self string) *Resource {
	return &Resource{
		Data: newResourceData(p, self),
	}
}

// Bind implements render.Binder
func (resource *Resource) Bind(r *http.Request) error {
	if resource.Data == nil || resource.Data.Type != Type {
		return fmt.Errorf("no data")
	}

	_, err := decimal.NewFromString(resource.Data.Amount)
	if err != nil {
		return err
	}

	return nil
}

// Render implements render.Render
func (resource *Resource) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// ListResource is a list of payments resource
type ListResource struct {
	Data []*ResourceData `json:"data"`
	links.Resource
}

// NewListResource returns new payments list resource
func NewListResource(payments []Payment, self string) *ListResource {
	listResource := &ListResource{
		Data:     []*ResourceData{},
		Resource: links.Resource{Links: links.Links{Self: self}},
	}
	for i := range payments {
		listResource.Data = append(
			listResource.Data, newResourceData(&payments[i], self+"/"+payments[i].ID),
		)
	}

	return listResource
}

// Render implements render.Render
func (list *ListResource) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
