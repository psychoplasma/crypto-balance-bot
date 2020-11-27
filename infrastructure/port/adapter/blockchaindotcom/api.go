package blockchaindotcom

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter"
)

var (
	pageLimit = 50
)

// Input is a data structure returning from Blockchain.com API
type Input struct {
	PrevOutput Output `json:"prev_out"`
}

// Output is a data structure returning from Blockchain.com API
type Output struct {
	Value   *big.Int `json:"value"`
	Address string   `json:"addr"`
}

// Transaction is a data structure returning from Blockchain.com API
type Transaction struct {
	BlockHeight int      `json:"block_height"`
	Hash        string   `json:"hash"`
	Inputs      []Input  `json:"inputs"`
	Outputs     []Output `json:"out"`
}

// AddressInfo is a data structure returning from Blockchain.com API
type AddressInfo struct {
	Address string        `json:"address"`
	TxCount int           `json:"n_tx"`
	Balance *big.Int      `json:"final_balance"`
	Txs     []Transaction `json:"txs"`
}

// BitcoinAPI implements CurrencyAPI for Bitcoin
type BitcoinAPI struct {
	t adapter.Translator
}

// NewBitcoinAPI creates a new instance of BitcoinAPI
func NewBitcoinAPI(t adapter.Translator) *BitcoinAPI {
	return &BitcoinAPI{
		t: t,
	}
}

// GetTxsOfAddress fetches txs of the given address since the given block height.
// Worth to mention that this does not fetches since exact sinceBlockHeight,
// rather guarantees that txs at sinceBlockHeight will be included. There may be
// past transactions as well. Therefore the changes should be applied in
// an idempotent way in the domain.
func (a *BitcoinAPI) GetTxsOfAddress(address string, sinceBlockHeight int) (*domain.AccountMovement, error) {
	txs := []Transaction{}
	ai, err := a.fetchAddressInfo(address, pageLimit, 0)
	if err != nil {
		return nil, err
	}

	txs = append(ai.Txs, txs...)
	txCount := len(txs)
	currPage := 0

	for txCount < ai.TxCount && ai.Txs[pageLimit-1].BlockHeight > sinceBlockHeight {
		currPage++
		ai, err = a.fetchAddressInfo(address, pageLimit, pageLimit*currPage)
		if err != nil {
			return nil, err
		}
		txs = append(txs, ai.Txs...)
		txCount += len(ai.Txs)
	}

	return a.t.ToAccountMovements(address, txs), nil
}

// API call to https://blockchain.info/rawaddr/$bitcoin_address.
// For further info: https://www.blockchain.com/api/blockchain_api
func (a *BitcoinAPI) fetchAddressInfo(address string, pLimit int, pOffset int) (*AddressInfo, error) {
	url := fmt.Sprintf("https://blockchain.info/rawaddr/%s?n=%d&offset=%d", address, pLimit, pOffset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}

	ad := &AddressInfo{}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(ad); err != nil {
		return nil, err
	}

	return ad, nil
}
