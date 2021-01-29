package blockbook

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain"
)

var (
	pageLimit = 50
)

// Paging is a data structure returning from Blockbook's API
type Paging struct {
	Page        int `json:"page"`
	TotalPages  int `json:"totalPages"`
	ItemsOnPage int `json:"itemsOnPage"`
}

// Input is a data structure returning from Blockbook's API
type Input struct {
	TxID      string   `json:"txid"`
	VOut      int      `json:"vout"`
	Sequence  int      `json:"sequence"`
	Index     int      `json:"n"`
	Addresses []string `json:"addresses"`
	IsAddress bool     `json:"isAddress"`
	Value     string   `json:"value"`
	Hex       string   `json:"hex"`
}

// Output is a data structure returning from Blockbook's API
type Output struct {
	Index     int      `json:"n"`
	Value     string   `json:"value"`
	Spent     bool     `json:"spent"`
	Addresses []string `json:"addresses"`
	IsAddress bool     `json:"isAddress"`
	Hex       string   `json:"hex"`
}

// Transaction is a data structure returning from Blockbook's API
type Transaction struct {
	BlockHeight   int      `json:"blockHeight"`
	BlockHash     string   `json:"blockHash"`
	BlockTime     int      `json:"blockTime"`
	Confirmations int      `json:"confirmations"`
	Value         string   `json:"value"`
	ValueIn       string   `json:"valueIn"`
	Fees          string   `json:"fees"`
	Hex           string   `json:"hex"`
	TxID          string   `json:"txid"`
	Version       int      `json:"version"`
	Inputs        []Input  `json:"vin"`
	Outputs       []Output `json:"vout"`
}

// AddressTxs is a data structure returning from Blockchain.com API
type AddressTxs struct {
	Paging
	Address            string        `json:"address"`
	Balance            string        `json:"balance"`
	UnconfirmedBalance string        `json:"unconfirmedBalance"`
	UnconfirmedTxs     int           `json:"unconfirmedTxs"`
	TotalReceived      string        `json:"totalReceived"`
	TotalSent          string        `json:"totalSent"`
	TxCount            int           `json:"txs"`
	Txs                []Transaction `json:"transactions"`
}

// BitcoinAPI implements CurrencyAPI for Bitcoin
type BitcoinAPI struct {
	hostURL string
	t       blockchain.Translator
}

// NewBitcoinAPI creates a new instance of BitcoinAPI
func NewBitcoinAPI(hostURL string, t blockchain.Translator) *BitcoinAPI {
	return &BitcoinAPI{
		hostURL: hostURL,
		t:       t,
	}
}

// GetTxsOfAddress fetches txs of the given address since the given block height
func (a *BitcoinAPI) GetTxsOfAddress(address string, sinceBlockHeight int) (*domain.AccountMovements, error) {
	txs := []Transaction{}
	currPage := 1
	at, err := a.fetchAddressTxs(address, sinceBlockHeight, currPage)
	if err != nil {
		return nil, err
	}

	txs = append(at.Txs, txs...)
	totalPages := at.TotalPages

	for currPage < totalPages {
		currPage++
		at, err = a.fetchAddressTxs(address, sinceBlockHeight, currPage)
		if err != nil {
			return nil, err
		}
		txs = append(txs, at.Txs...)
	}

	return a.t.ToAccountMovements(address, txs)
}

// API call to blockbook's api/v2/address endpoint
// For further info: https://github.com/trezor/blockbook/blob/master/docs/api.md#get-address
func (a *BitcoinAPI) fetchAddressTxs(address string, since int, page int) (*AddressTxs, error) {
	url := fmt.Sprintf("%s/api/v2/address/%s?details=txs&page=%d&pageSize=%d&from=%d", a.hostURL, address, page, pageLimit, since)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}

	ad := &AddressTxs{}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(ad); err != nil {
		return nil, err
	}

	return ad, nil
}
