package blockchaindotcom

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

type input struct {
	PrevOutput output `json:"prev_out"`
}

type output struct {
	Index   int      `json:"n"`
	Value   *big.Int `json:"value"`
	Address string   `json:"addr"`
	Spent   bool     `json:"spent"`
}

// Transaction is a data structure returning from Blockchain.com API
type Transaction struct {
	BlockHeight int      `json:"block_height"`
	Hash        string   `json:"hash"`
	Inputs      []input  `json:"inputs"`
	Outputs     []output `json:"out"`
}

type addressInfo struct {
	Address string        `json:"address"`
	TxCount int           `json:"n_tx"`
	Balance *big.Int      `json:"final_balance"`
	Txs     []Transaction `json:"txs"`
}

// BitcoinAPI implements CurrencyAPI for Bitcoin
type BitcoinAPI struct {
	T BitcoinTranslator
}

// GetAddressTxs fetches txs of the given address since the given block height(exclusive)
func (a *BitcoinAPI) GetAddressTxs(address string, sinceBlockHeight int) ([]*domain.AccountMovement, error) {
	d, err := a.fetchAddressInfo(address)
	if err != nil {
		return nil, err
	}

	return a.T.ToAccountMovements(address, d), nil
}

func (a *BitcoinAPI) fetchAddressInfo(address string) ([]Transaction, error) {
	endpoint := fmt.Sprintf("https://blockchain.info/rawaddr/%s", address)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}

	ad := &addressInfo{}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(ad); err != nil {
		return nil, err
	}

	return ad.Txs, nil
}
