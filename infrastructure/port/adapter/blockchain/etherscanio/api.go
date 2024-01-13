package etherscanio

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/net"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain"
)

var (
	pageLimit = 50
)

const transactionStatusSuccess = "1"

// Response is a data structure returning from Etherscan.io API
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

// Transaction is a data structure returning from Etherscan.io API
type Transaction struct {
	BlockHeight string `json:"blockNumber"`
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	Status      string `json:"txreceipt_status"`
	Timestamp   string `json:"timeStamp"`
}

// API implements CurrencyAPI for Bitcoin
type API struct {
	t blockchain.Translator
}

// Delay between consecutive api requests, not to choking api provider
const requestDelay = 200 * time.Millisecond

// NewAPI creates a new instance of API
func NewAPI(t blockchain.Translator) *API {
	return &API{
		t: t,
	}
}

// GetAccountMovements fetches txs of the given address since the given block height
func (a *API) GetAccountMovements(address string, sinceBlockHeight uint64) (*domain.AccountMovements, error) {
	txs, err := a.fetchAddressTxs(address, sinceBlockHeight)
	if err != nil {
		return nil, err
	}

	return a.t.ToAccountMovements(address, txs)
}

// GetLatestBlockHeight fetches the latest block number
func (a *API) GetLatestBlockHeight() (uint64, error) {
	return a.fetchBlockHeightByTimestamp(time.Now().Unix())
}

// API call to https://api.etherscan.io/api?module=account&action=txlist&address=
// For further info: https://etherscan.io/apis#accounts
func (a *API) fetchAddressTxs(address string, startBlock uint64) ([]Transaction, error) {
	defer time.Sleep(requestDelay)
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=%d&sort=desc", address, startBlock)
	r := &Response{}
	if err := net.GetJSON(url, r); err != nil {
		return nil, err
	}

	d, err := json.Marshal(r.Result)
	if err != nil {
		return nil, err
	}

	txs := []Transaction{}
	if err := json.Unmarshal(d, &txs); err != nil {
		return nil, err
	}

	if r.Status != "0" && len(txs) == 0 {
		return nil, errors.New(r.Message)
	}

	return txs, nil
}

// API call to https://api.etherscan.io/api?module=block&action=getblocknobytime
// For further info: https://etherscan.io/apis#blocks
func (a *API) fetchBlockHeightByTimestamp(timestamp int64) (uint64, error) {
	defer time.Sleep(requestDelay)
	url := fmt.Sprintf("https://api.etherscan.io/api?module=block&action=getblocknobytime&timestamp=%d&closest=before", timestamp)
	r := &Response{}
	if err := net.GetJSON(url, r); err != nil {
		return 0, err
	}

	blockHeight, err := strconv.ParseUint(r.Result.(string), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("status: %s, %s", r.Message, err.Error())
	}

	return blockHeight, nil
}
