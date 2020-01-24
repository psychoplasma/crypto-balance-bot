package cryptobot

// Account represents Account to be subscribed to bot
type Account struct {
	Name        string
	CurrencyID  string
	IsHD        bool
	HDParams    *HD
	AddressList []string
}

// HD represents HD wallet parameters to drive accounts for the given master public key
type HD struct {
	MasterPubKey string
	HdPath       string
	AccountIndex string
}

func (a *Account) AddByAddress() error {
	return nil
}

func (a *Account) AddByPubKey() error {
	return nil
}

func (a *Account) AddByMasterPubKey() error {
	return nil
}

func (a *Account) RemoveByAddress() error {
	return nil
}

func (a *Account) RemoveByAsset() error {
	return nil
}

func (a *Account) RemoveAll() error {
	return nil
}
