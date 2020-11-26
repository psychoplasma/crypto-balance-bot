package etherscanio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter"
)

var (
	pageLimit = 50
)

const transactionStatusSuccess = "1"

// Response is a data structure returning from Etherscan.io API
type Response struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []Transaction `json:"result"`
}

// Transaction is a data structure returning from Etherscan.io API
type Transaction struct {
	BlockHeight string `json:"blockNumber"`
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	Status      string `json:"txreceipt_status"`
}

// EthereumAPI implements CurrencyAPI for Bitcoin
type EthereumAPI struct {
	t adapter.Translator
}

// NewEthereumAPI creates a new instance of EthereumAPI
func NewEthereumAPI(t adapter.Translator) *EthereumAPI {
	return &EthereumAPI{
		t: t,
	}
}

// GetTxsOfAddress fetches txs of the given address since the given block height
func (a *EthereumAPI) GetTxsOfAddress(address string, sinceBlockHeight int) ([]*domain.AccountMovement, error) {
	txs, err := a.fetchAddressTxs(address, sinceBlockHeight)
	if err != nil {
		return nil, err
	}

	return a.t.ToAccountMovements(address, txs), nil
}

// API call to https://api.etherscan.io/api?module=account&action=txlist&address=.
// For further info: https://etherscan.io/apis#accounts
func (a *EthereumAPI) fetchAddressTxs(address string, startBlock int) ([]Transaction, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=%d&sort=desc", address, startBlock)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}

	r := &Response{}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		return nil, err
	}

	if r.Status != "0" && len(r.Result) == 0 {
		return nil, errors.New(r.Message)
	}

	return r.Result, nil
}
