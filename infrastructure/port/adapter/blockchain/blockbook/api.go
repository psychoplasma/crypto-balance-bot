package blockbook

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain"
)

const (
	transactionStatusSuccess = 1
	defaultPagingLimit       = 100
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
	VOut      uint     `json:"vout"`
	Sequence  uint64   `json:"sequence"`
	Index     uint     `json:"n"`
	Addresses []string `json:"addresses"`
	IsAddress bool     `json:"isAddress"`
	Value     string   `json:"value"`
	Hex       string   `json:"hex"`
}

// Output is a data structure returning from Blockbook's API
type Output struct {
	Index     uint     `json:"n"`
	Value     string   `json:"value"`
	Spent     bool     `json:"spent"`
	Addresses []string `json:"addresses"`
	IsAddress bool     `json:"isAddress"`
	Hex       string   `json:"hex"`
}

// Transaction is a data structure returning from Blockbook's API
type Transaction struct {
	BlockHeight      uint64           `json:"blockHeight"`
	BlockHash        string           `json:"blockHash"`
	BlockTime        uint64           `json:"blockTime"`
	Confirmations    uint64           `json:"confirmations"`
	EthereumSpecific EthereumSpecific `json:"ethereumSpecific,omitempty"`
	Value            string           `json:"value"`
	ValueIn          string           `json:"valueIn"`
	Fees             string           `json:"fees"`
	Hex              string           `json:"hex"`
	TxID             string           `json:"txid"`
	Version          int              `json:"version"`
	Inputs           []Input          `json:"vin"`
	Outputs          []Output         `json:"vout"`
	TokenTransfers   []TokenTransfer  `json:"tokenTransfers,omitempty"`
}

// TokenTransfer contains info about a token transfer done in a transaction
type TokenTransfer struct {
	Type     string `json:"type"`
	From     string `json:"from"`
	To       string `json:"to"`
	Token    string `json:"token"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals uint   `json:"decimals"`
	Value    string `json:"value"`
}

// EthereumSpecific contains ethereum specific transaction data
type EthereumSpecific struct {
	Status   int    `json:"status"` // 1 OK, 0 Fail, -1 pending
	Nonce    uint64 `json:"nonce"`
	GasLimit uint   `json:"gasLimit"`
	GasUsed  uint   `json:"gasUsed"`
	GasPrice string `json:"gasPrice"`
	Data     string `json:"data,omitempty"`
}

// AddressTxs is a data structure returning from Blockbook's API
type AddressTxs struct {
	Paging
	Address            string        `json:"address"`
	Balance            string        `json:"balance"`
	UnconfirmedBalance string        `json:"unconfirmedBalance"`
	UnconfirmedTxs     int           `json:"unconfirmedTxs"`
	TotalReceived      string        `json:"totalReceived"`
	TotalSent          string        `json:"totalSent"`
	TxCount            uint64        `json:"txs"`
	Transactions       []Transaction `json:"transactions"`
}

// API implements CurrencyAPI for Blockbook
type API struct {
	hostURL     string
	pagingLimit int
	t           blockchain.Translator
}

// NewAPI creates a new instance of BitcoinAPI
func NewAPI(hostURL string, t blockchain.Translator, pagingLimit ...*int) *API {
	api := &API{
		hostURL:     hostURL,
		pagingLimit: defaultPagingLimit,
		t:           t,
	}

	for _, pl := range pagingLimit {
		if pl != nil {
			api.pagingLimit = *pl
		}
	}

	return api
}

// GetAccountMovements fetches txs of the given address since the given block height
func (a *API) GetAccountMovements(address string, sinceBlockHeight uint64) (*domain.AccountMovements, error) {
	currPage := 1
	at, err := a.fetchAddressTxs(address, sinceBlockHeight, currPage)
	if err != nil {
		return nil, err
	}

	txs := at.Transactions
	totalPages := at.TotalPages

	for currPage < totalPages {
		currPage++
		at, err = a.fetchAddressTxs(address, sinceBlockHeight, currPage)
		if err != nil {
			return nil, err
		}
		txs = append(txs, at.Transactions...)
	}

	return a.t.ToAccountMovements(address, txs)
}

// API call to blockbook's api/v2/address endpoint
// For further info: https://github.com/trezor/blockbook/blob/master/docs/api.md#get-address
func (a *API) fetchAddressTxs(address string, since uint64, page int) (*AddressTxs, error) {
	url := fmt.Sprintf("%s/api/v2/address/%s?details=txs&page=%d&pageSize=%d&from=%d", a.hostURL, address, page, a.pagingLimit, since)
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
