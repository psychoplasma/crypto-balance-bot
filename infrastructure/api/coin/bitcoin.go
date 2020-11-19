package coin

import (
	"errors"
	"math/big"
)

// BitcoinAPI implements CurrencyAPI for Bitcoin
type BitcoinAPI struct {
}

func (a *BitcoinAPI) GetBalance(addressDesc string) (*big.Int, error) {
	return nil, nil
}

func (a *BitcoinAPI) GetTransactions(addressDesc string, since int) ([]string, error) {
	return nil, nil
}

func (a *BitcoinAPI) CreateAddress(pubKey string) (string, error) {
	return "", nil
}

func (a *BitcoinAPI) DeriveAddressFromXPubKey(xPubKey string) ([]string, error) {
	return nil, nil
}

func (a *BitcoinAPI) ValidateAddress(addressDesc string) error {
	return nil
}

func (a *BitcoinAPI) ValidatePubKey(pubKey string) error {
	return errors.New("not implemented")
}
