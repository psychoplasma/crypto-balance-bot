package cryptobot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
)

// FilterType represents types of filters
type FilterType string

// Defined filter types
const (
	Amount     FilterType = "amount"
	AddressOn  FilterType = "addressOn"
	AddressOff FilterType = "blacklistOff"
)

// Filter represents a domain entity which decides
// at what condition notifications will be published
type Filter struct {
	c      condition
	isMust bool
	t      FilterType
}

// NewAmountFilter creates a new instance of Amount type of Filter
func NewAmountFilter(amount *big.Int, must bool) (*Filter, error) {
	if amount == nil {
		return nil, fmt.Errorf("nil amount")
	}
	return NewFilter(Amount, &amountCondition{Amount: amount}, must)
}

// NewAddressOnFilter creates a new instance of AddressOn type of Filter
func NewAddressOnFilter(address string, must bool) (*Filter, error) {
	if address == "" {
		return nil, fmt.Errorf("empty address")
	}
	return NewFilter(AddressOn, &addressOnCondition{Address: address}, must)
}

// NewAddressOffFilter creates a new instance of AddressOff type of Filter
func NewAddressOffFilter(address string, must bool) (*Filter, error) {
	if address == "" {
		return nil, fmt.Errorf("empty address")
	}
	return NewFilter(AddressOff, &addressOffCondition{Address: address}, must)
}

// NewFilter creates a new instance of Filter
func NewFilter(t FilterType, c condition, must bool) (*Filter, error) {
	return &Filter{
		c:      c,
		isMust: must,
		t:      t,
	}, nil
}

// CheckCondition checks whether or not the given conditions satisfy this filter
func (f *Filter) CheckCondition(condition interface{}) bool {
	return f.c.CheckAgainst(condition)
}

// IsMust returns true if this filter always must be satisfied to publish a notification
// regardless of the other filters, false otherwise
func (f *Filter) IsMust() bool {
	return f.isMust
}

// Type returns type property
func (f *Filter) Type() FilterType {
	return f.t
}

// SerializeCondition serializes the condition data
// which can be different for each type of condition
func (f *Filter) SerializeCondition() ([]byte, error) {
	return f.c.Serialize()
}

// DeserializeCondition deserializes the given data to
// the corresponding condition according to type of the filter
func (f *Filter) DeserializeCondition(data []byte) error {
	switch f.t {
	case Amount:
		f.c = new(amountCondition)
		return f.c.Deserialize(data)
	case AddressOn:
		f.c = new(addressOnCondition)
		return f.c.Deserialize(data)
	case AddressOff:
		f.c = new(addressOffCondition)
		return f.c.Deserialize(data)
	default:
		return fmt.Errorf("unrecognized filter type: %s", f.t)
	}
}

// Condition represents condition parameters and
// its condition check for a specific type of Filter
type condition interface {
	CheckAgainst(i interface{}) bool
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
}

type amountCondition struct {
	Amount *big.Int `json:"amount"`
}

func (c *amountCondition) CheckAgainst(amount interface{}) bool {
	a, _ := amount.(*big.Int)
	return a.Cmp(c.Amount) > -1
}

func (c *amountCondition) Serialize() ([]byte, error) {
	return json.Marshal(c)
}

func (c *amountCondition) Deserialize(data []byte) error {
	return decodeJSONStrictly(data, c)
}

type addressOnCondition struct {
	Address string `json:"address"`
}

func (c *addressOnCondition) CheckAgainst(address interface{}) bool {
	return c.Address == address.(string)
}

func (c *addressOnCondition) Serialize() ([]byte, error) {
	return json.Marshal(c)
}

func (c *addressOnCondition) Deserialize(data []byte) error {
	return decodeJSONStrictly(data, c)
}

type addressOffCondition struct {
	Address string `json:"address"`
}

func (c *addressOffCondition) CheckAgainst(address interface{}) bool {
	return c.Address != address.(string)
}

func (c *addressOffCondition) Serialize() ([]byte, error) {
	return json.Marshal(c)
}

func (c *addressOffCondition) Deserialize(data []byte) error {
	return decodeJSONStrictly(data, c)
}

func decodeJSONStrictly(data []byte, i interface{}) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	return d.Decode(i)
}
