package main

import (
	"net/http"

	"github.com/VMitov/payments/pkg/payment"
	"github.com/go-chi/render"
)

func listPayments(w http.ResponseWriter, r *http.Request) {
	payments := []*payment.Payment{{ID: "1"}}
	if err := render.Render(w, r, newPaymentListResponse(payments)); err != nil {
		render.Render(w, r, errRender(err))
		return
	}
}

func newPaymentListResponse(paymentList []*payment.Payment) *payment.ListResource {
	listResource := &payment.ListResource{}
	for _, payment := range paymentList {
		listResource.Data = append(listResource.Data, newPaymentResponse(payment))
	}

	return listResource
}

func newPaymentResponse(p *payment.Payment) *payment.Resource {
	return &payment.Resource{Payment: p, Type: payment.Type}
}
