package coin

import (
	"errors"
	"math/big"
)

// EthereumAPI implements CurrencyAPI for Ethereum
type EthereumAPI struct {
}

func (a *EthereumAPI) GetBalance(addressDesc string) (*big.Int, error) {
	return nil, nil
}

func (a *EthereumAPI) GetTransactions(addressDesc string, since int) ([]string, error) {
	return nil, nil
}

func (a *EthereumAPI) CreateAddress(pubKey string) (string, error) {
	return "", nil
}

func (a *EthereumAPI) DeriveAddressFromXPubKey(xPubKey string) ([]string, error) {
	return nil, nil
}

func (a *EthereumAPI) ValidateAddress(addressDesc string) error {
	return nil
}

func (a *EthereumAPI) ValidatePubKey(pubKey string) error {
	return errors.New("not implemented")
}
