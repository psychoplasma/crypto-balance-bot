package cryptobot_test

import (
	"math/big"
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

func TestCheckAgainst_WithAmountType(t *testing.T) {
	if _, err := domain.NewAmountFilter(nil, true); err == nil {
		t.Fatalf("expected an error but got nothing")
	}

	amountLess := big.NewInt(4)
	amountMore := big.NewInt(6)
	filterAmount := big.NewInt(5)
	f, err := domain.NewAmountFilter(filterAmount, true)
	if err != nil {
		t.Fatal(err)
	}

	if !f.CheckCondition(amountMore) {
		t.Fatalf("expected to pass the condition check but failed")
	}

	if f.CheckCondition(amountLess) {
		t.Fatalf("expected to fail the condition check but passed")
	}
}

func TestCheckAgainst_WithAddressOnType(t *testing.T) {
	if _, err := domain.NewAddressOnFilter("", true); err == nil {
		t.Fatalf("expected an error but got nothing")
	}

	addressSame := "address-1"
	addressDiff := "address-2"
	filterAddress := "address-1"
	f, err := domain.NewAddressOnFilter(filterAddress, true)
	if err != nil {
		t.Fatal(err)
	}

	if !f.CheckCondition(addressSame) {
		t.Fatalf("expected to pass the condition check but failed")
	}

	if f.CheckCondition(addressDiff) {
		t.Fatalf("expected to fail the condition check but passed")
	}
}

func TestCheckAgainst_WithAddressOffType(t *testing.T) {
	if _, err := domain.NewAddressOffFilter("", true); err == nil {
		t.Fatalf("expected an error but got nothing")
	}

	addressSame := "address-1"
	addressDiff := "address-2"
	filterAddress := "address-1"
	f, err := domain.NewAddressOffFilter(filterAddress, true)
	if err != nil {
		t.Fatal(err)
	}

	if !f.CheckCondition(addressDiff) {
		t.Fatalf("expected to pass the condition check but failed")
	}

	if f.CheckCondition(addressSame) {
		t.Fatalf("expected to fail the condition check but passed")
	}
}

func TestSerializeCondition_WithAmountType(t *testing.T) {
	expected := "{\"amount\":5}"
	f, err := domain.NewAmountFilter(big.NewInt(5), true)
	if err != nil {
		t.Fatal(err)
	}

	d, err := f.SerializeCondition()
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != expected {
		t.Fatalf("expected \"%s\" but got \"%s\"", expected, string(d))
	}
}

func TestSerializeCondition_WithAddressOnType(t *testing.T) {
	expected := "{\"address\":\"address-1\"}"
	f, err := domain.NewAddressOnFilter("address-1", true)
	if err != nil {
		t.Fatal(err)
	}

	d, err := f.SerializeCondition()
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != expected {
		t.Fatalf("expected \"%s\" but got \"%s\"", expected, string(d))
	}
}

func TestSerializeCondition_WithAddressOffType(t *testing.T) {
	expected := "{\"address\":\"address-1\"}"
	f, err := domain.NewAddressOffFilter("address-1", true)
	if err != nil {
		t.Fatal(err)
	}

	d, err := f.SerializeCondition()
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != expected {
		t.Fatalf("expected \"%s\" but got \"%s\"", expected, string(d))
	}
}

func TestDeserializeCondition_WithAmountType(t *testing.T) {
	condition := "{\"amount\":5}"
	amountLess := big.NewInt(4)
	amountMore := big.NewInt(6)
	f, err := domain.NewFilter(domain.Amount, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	if err := f.DeserializeCondition([]byte(condition)); err != nil {
		t.Fatal(err)
	}

	if !f.CheckCondition(amountMore) {
		t.Fatalf("expected to pass the condition check but failed")
	}

	if f.CheckCondition(amountLess) {
		t.Fatalf("expected to fail the condition check but passed")
	}
}

func TestDeserializeCondition_WithAddressOnType(t *testing.T) {
	condition := "{\"address\":\"address-1\"}"
	addressSame := "address-1"
	addressDiff := "address-2"
	f, err := domain.NewFilter(domain.AddressOn, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	if err := f.DeserializeCondition([]byte(condition)); err != nil {
		t.Fatal(err)
	}

	if !f.CheckCondition(addressSame) {
		t.Fatalf("expected to pass the condition check but failed")
	}

	if f.CheckCondition(addressDiff) {
		t.Fatalf("expected to fail the condition check but passed")
	}
}

func TestDeserializeCondition_WithAddressOffType(t *testing.T) {
	condition := "{\"address\":\"address-1\"}"
	addressSame := "address-1"
	addressDiff := "address-2"
	f, err := domain.NewFilter(domain.AddressOff, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	if err := f.DeserializeCondition([]byte(condition)); err != nil {
		t.Fatal(err)
	}

	if !f.CheckCondition(addressDiff) {
		t.Fatalf("expected to pass the condition check but failed")
	}

	if f.CheckCondition(addressSame) {
		t.Fatalf("expected to fail the condition check but passed")
	}
}

func TestDeserializeCondition_WithUnknownType(t *testing.T) {
	condition := "{\"address\":\"address-1\"}"
	f, err := domain.NewFilter(domain.FilterType("asdadas"), nil, true)
	if err != nil {
		t.Fatal(err)
	}

	if err := f.DeserializeCondition([]byte(condition)); err == nil {
		t.Fatalf("expected an error but got nothing")
	}
}

func TestDeserializeCondition_WithErrornousCondition(t *testing.T) {
	condition := "{\"whatdowewant\":\"moretests\", \"amount\": 5}"
	f, err := domain.NewFilter(domain.Amount, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	if err := f.DeserializeCondition([]byte(condition)); err == nil {
		t.Fatalf("expected an error but got nothing")
	}
}
