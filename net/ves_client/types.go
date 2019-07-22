package vesclient

// ECCKey is the private key object in memory
type ECCKey struct {
	PrivateKey []byte `json:"private_key"`
	ChainID    uint64 `json:"chain_id"`
}

// ECCKeyAlias is the private key object in json
type ECCKeyAlias struct {
	PrivateKey string `json:"private_key"`
	ChainID    uint64 `json:"chain_id"`
	Alias      string `json:"alias"`
}

// EthAccount is the account object in memory
type EthAccount struct {
	Address    string `json:"address"`
	ChainID    uint64 `json:"chain_id"`
	PassPhrase string `json:"pass_phrase"`
}

// EthAccountAlias is the account object in json
type EthAccountAlias struct {
	EthAccount
	Alias string `json:"alias"`
}

// ECCKeys is the object saved in files
type ECCKeys struct {
	Keys  []*ECCKey
	Alias map[string]ECCKey
}

// EthAccounts is the object saved in files
type EthAccounts struct {
	Accs  []*EthAccount
	Alias map[string]EthAccount
}
