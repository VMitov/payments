package main

import (
	"net/http"

	"github.com/VMitov/payments/pkg/payment"
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
