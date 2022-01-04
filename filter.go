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
	AddressOff FilterType = "addressOff"
)

// Filter represents a domain entity which decides
// at what condition notifications will be published
type Filter struct {
	c      condition
	isMust bool
	t      FilterType
}

// NewAmountFilter creates a new instance of Amount type of Filter
func NewAmountFilter(amount string, must bool) (*Filter, error) {
	a, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf("amount(%s) is not a valid number representation", amount)
	}
	return NewFilter(Amount, &amountCondition{Amount: a}, must), nil
}

// NewAddressOnFilter creates a new instance of AddressOn type of Filter
func NewAddressOnFilter(address string, must bool) (*Filter, error) {
	if address == "" {
		return nil, fmt.Errorf("empty address")
	}
	return NewFilter(AddressOn, &addressOnCondition{Address: address}, must), nil
}

// NewAddressOffFilter creates a new instance of AddressOff type of Filter
func NewAddressOffFilter(address string, must bool) (*Filter, error) {
	if address == "" {
		return nil, fmt.Errorf("empty address")
	}
	return NewFilter(AddressOff, &addressOffCondition{Address: address}, must), nil
}

// NewFilter creates a new instance of Filter
func NewFilter(t FilterType, c condition, must bool) *Filter {
	return &Filter{
		c:      c,
		isMust: must,
		t:      t,
	}
}

// CheckCondition checks whether or not the given conditions satisfy this filter
func (f *Filter) CheckCondition(t *Transfer) bool {
	return f.c.CheckAgainst(t)
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

// ToString returns human-readable string representation of this filter
func (f *Filter) ToString() string {
	return fmt.Sprintf("filters any transfer whose %s. Is this a must? -> %v", f.c.ToString(), f.isMust)
}

// Condition represents condition parameters and
// its condition check for a specific type of Filter
type condition interface {
	CheckAgainst(t *Transfer) bool
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
	ToString() string
}

type amountCondition struct {
	Amount *big.Int `json:"amount"`
}

func (c *amountCondition) CheckAgainst(t *Transfer) bool {
	return t.Amount.Cmp(c.Amount) > -1
}

func (c *amountCondition) Serialize() ([]byte, error) {
	return json.Marshal(c)
}

func (c *amountCondition) Deserialize(data []byte) error {
	return decodeJSONStrictly(data, c)
}

func (c *amountCondition) ToString() string {
	return fmt.Sprintf("amount is greater than or equal to %s", c.Amount.String())
}

type addressOnCondition struct {
	Address string `json:"address"`
}

func (c *addressOnCondition) CheckAgainst(t *Transfer) bool {
	return c.Address == t.Address
}

func (c *addressOnCondition) Serialize() ([]byte, error) {
	return json.Marshal(c)
}

func (c *addressOnCondition) Deserialize(data []byte) error {
	return decodeJSONStrictly(data, c)
}

func (c *addressOnCondition) ToString() string {
	return fmt.Sprintf("second-party address is equal to \"%s\"", c.Address)
}

type addressOffCondition struct {
	Address string `json:"address"`
}

func (c *addressOffCondition) CheckAgainst(t *Transfer) bool {
	return c.Address != t.Address
}

func (c *addressOffCondition) Serialize() ([]byte, error) {
	return json.Marshal(c)
}

func (c *addressOffCondition) Deserialize(data []byte) error {
	return decodeJSONStrictly(data, c)
}

func (c *addressOffCondition) ToString() string {
	return fmt.Sprintf("second-party address is different than \"%s\"", c.Address)
}

func decodeJSONStrictly(data []byte, i interface{}) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	return d.Decode(i)
}
