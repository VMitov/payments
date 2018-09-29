package main

import (
	"net/http"

	"github.com/VMitov/payments/pkg/payment"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func newPayment(p *payment.Payment) *payment.Resource {
	return payment.NewResource(p, "/payments/"+p.ID)
}

func newPaymentList(ps []payment.Payment) *payment.ListResource {
	return payment.NewListResource(ps, "/payments")
}

func (api *api) createPayment(w http.ResponseWriter, r *http.Request) {
	data := &payment.Resource{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, errInvalidRequest(err))
		return
	}

	pay, err := payment.NewFromResource(data)
	if err != nil {
		render.Render(w, r, errInvalidRequest(err))
		return
	}

	id, err := payment.Create(api.db, pay)
	if err != nil {
		render.Render(w, r, errSystem(err))
		return
	}

	newPay, err := payment.Get(api.db, id)
	if err != nil {
		render.Render(w, r, errSystem(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, newPayment(newPay))
}

func (api *api) updatePayment(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "paymentID")
	if paymentID == "" {
		render.Render(w, r, errNotFound)
		return
	}

	data := &payment.Resource{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, errInvalidRequest(err))
		return
	}

	pay, err := payment.Get(api.db, paymentID)
	if err != nil {
		render.Render(w, r, errInvalidRequest(err))
		return
	}
	if pay.ID == "" {
		render.Render(w, r, errNotFound)
		return
	}

	newPay, err := payment.NewFromResource(data)
	if err != nil {
		render.Render(w, r, errSystem(err))
		return
	}

	if err := payment.Update(api.db, paymentID, newPay); err != nil {
		render.Render(w, r, errSystem(err))
		return
	}

	newPay, err = payment.Get(api.db, paymentID)
	if err != nil {
		render.Render(w, r, errSystem(err))
		return
	}

	render.Render(w, r, newPayment(newPay))
}

func (api *api) listPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := payment.Select(api.db)
	if err != nil {
		render.Render(w, r, errSystem(err))
		return
	}

	if err := render.Render(w, r, newPaymentList(payments)); err != nil {
		render.Render(w, r, errSystem(err))
		return
	}
}

func (api *api) getPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "paymentID")
	if paymentID == "" {
		render.Render(w, r, errNotFound)
		return
	}

	pay, err := payment.Get(api.db, paymentID)
	if err != nil {
		render.Render(w, r, errSystem(err))
		return
	}

	if pay.ID == "" {
		render.Render(w, r, errNotFound)
		return
	}

	if err := render.Render(w, r, newPayment(pay)); err != nil {
		render.Render(w, r, errSystem(err))
		return
	}
}

func (api *api) deletePayment(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "paymentID")
	if paymentID == "" {
		render.Render(w, r, errNotFound)
		return
	}

	pay, err := payment.Get(api.db, paymentID)
	if err != nil {
		render.Render(w, r, errSystem(err))
		return
	}

	if pay.ID == "" {
		render.Render(w, r, errNotFound)
		return
	}

	if err := payment.Delete(api.db, paymentID); err != nil {
		render.Render(w, r, errSystem(err))
		return
	}
}
