package payment

import apipayment "github.com/nuxtblog/nuxtblog/api/payment"

func New() apipayment.IPaymentV1 { return &ControllerV1{} }

type ControllerV1 struct{}
