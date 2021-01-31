package blockchain

import "strings"

const ethereumAddressPrefix = "0x"

// NormalizeEthereumAddress normalizes the given ethereum address.
// Trims any leading trailing white space, adds a prefix of 0x
// if there is no, and converts to lowercases.
func NormalizeEthereumAddress(address string) string {
	address = strings.TrimSpace(address)

	if !strings.HasPrefix(address, ethereumAddressPrefix) {
		address = ethereumAddressPrefix + address
	}

	return strings.ToLower(address)
}
