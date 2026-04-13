package wallet

import apiwallet "github.com/nuxtblog/nuxtblog/api/wallet"

func New() apiwallet.IWalletV1 { return &ControllerV1{} }

type ControllerV1 struct{}
