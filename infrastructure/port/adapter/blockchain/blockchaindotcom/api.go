package blockchaindotcom

import (
	"fmt"
	"math/big"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/net"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain"
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
	BlockHeight uint64   `json:"block_height"`
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

// Block is a data structure returning from Blockchain.com API
type Block struct {
	Hash      string   `json:"hash"`
	Time      int64    `json:"time"`
	Index     uint64   `json:"block_index"`
	Height    uint64   `json:"height"`
	TxIndexes []uint64 `json:"txIndexes"`
}

// API implements CurrencyAPI for Bitcoin
type API struct {
	t blockchain.Translator
}

// NewAPI creates a new instance of API
func NewAPI(t blockchain.Translator) *API {
	return &API{
		t: t,
	}
}

// GetAccountMovements fetches txs of the given address since the given block height.
// Worth to mention that this does not fetches since exact sinceBlockHeight,
// rather guarantees that txs at sinceBlockHeight will be included. There may be
// past transactions as well. Therefore the changes should be applied in
// an idempotent way in the domain.
func (a *API) GetAccountMovements(address string, sinceBlockHeight uint64) (*domain.AccountMovements, error) {
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

	return a.t.ToAccountMovements(address, txs)
}

// GetLatestBlockHeight fetches the latest block number
func (a *API) GetLatestBlockHeight() (uint64, error) {
	b, err := a.fetchLatestBlock()
	if err != nil {
		return 0, err
	}

	return b.Height, nil
}

// API call to https://blockchain.info/rawaddr/$bitcoin_address.
// For further info: https://www.blockchain.com/api/blockchain_api
func (a *API) fetchAddressInfo(address string, pLimit int, pOffset int) (*AddressInfo, error) {
	url := fmt.Sprintf("https://blockchain.info/rawaddr/%s?n=%d&offset=%d", address, pLimit, pOffset)
	ad := &AddressInfo{}
	if err := net.GetJSON(url, ad); err != nil {
		return nil, err
	}

	return ad, nil
}

// API call to https://blockchain.info/latestblock
// For further info: https://www.blockchain.com/api/blockchain_api
func (a *API) fetchLatestBlock() (*Block, error) {
	b := &Block{}
	if err := net.GetJSON("https://blockchain.info/latestblock", b); err != nil {
		return nil, err
	}

	return b, nil
}
