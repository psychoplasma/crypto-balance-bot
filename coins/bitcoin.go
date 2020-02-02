package coins

import "math/big"

// BitcoinAPI implements CurrencyAPI for Bitcoin
type BitcoinAPI struct {
}

func (a *BitcoinAPI) GetBalance(addressDesc string) (*big.Int, error) {
	return nil, nil
}

func (a *BitcoinAPI) CreateAddress(pubKey string) (string, error) {
	return "", nil
}

func (a *BitcoinAPI) ValidateAddress(addressDesc string) error {
	return nil
}
