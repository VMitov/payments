package main

import (
	"net/http"

	"github.com/VMitov/payments/pkg/payment"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (api *api) listPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := payment.Select(api.db)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	if err := render.Render(w, r, newPaymentListResponse(payments)); err != nil {
		render.Render(w, r, errRender(err))
		return
	}
}

func (api *api) getPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "paymentID")
	if paymentID == "" {
		render.Render(w, r, errNotFound)
		return
	}

	payment, err := payment.Get(api.db, paymentID)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	if payment.ID == "" {
		render.Render(w, r, errNotFound)
		return
	}

	if err := render.Render(w, r, newPaymentResponse(payment)); err != nil {
		render.Render(w, r, errRender(err))
		return
	}
}

func newPaymentListResponse(paymentList []payment.Payment) *payment.ListResource {
	listResource := &payment.ListResource{
		Data: []*payment.Resource{},
	}
	for i := range paymentList {
		listResource.Data = append(listResource.Data, newPaymentResponse(&paymentList[i]))
	}

	return listResource
}

func newPaymentResponse(p *payment.Payment) *payment.Resource {
	return &payment.Resource{Payment: p, Type: payment.Type}
}
