package types

// Github repo structs
type Tree struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Sha  string `json:"sha"`
	URL  string `json:"url"`
}

type TreeResponse struct {
	Sha       string `json:"sha"`
	URL       string `json:"url"`
	Tree      []Tree `json:"tree"`
	Truncated bool   `json:"truncated"`
}

type Content struct {
	Content string `json:"content"`
	Sha     string `json:"sha"`
}

type File struct {
	Content string
	URL     string
}

// Chain struct from chain registry
type FeeToken struct {
	Denom            string `json:"denom"`
	FixedMinGasPrice int64  `json:"fixed_min_gas_price"`
	LowGasPrice      int64  `json:"low_gas_price"`
	AverageGasPrice  int64  `json:"average_gas_price"`
	HighGasPrice     int64  `json:"high_gas_price"`
}

type FeeTokens struct {
	FeeTokens []FeeToken `json:"fee_tokens"`
}

type StakingToken struct {
	Denom string `json:"denom"`
}

type StakingTokens struct {
	StakingTokens []StakingToken `json:"staking_tokens"`
}

type Endpoint struct {
	Address  string `json:"address"`
	Provider string `json:"provider"`
}

type Apis struct {
	RPC  []Endpoint `json:"rpc"`
	Rest []Endpoint `json:"rest"`
	Grpc []Endpoint `json:"grpc"`
	Evm  []Endpoint `json:"evm-http-jsonrpc"`
}

type Explorer struct {
	Kind   string `json:"kind"`
	URL    string `json:"url"`
	TxPage string `json:"tx_page"`
}

type Chain struct {
	ChainName    string        `json:"chain_name"`
	NetworkType  string        `json:"network_type"`
	PrettyName   string        `json:"pretty_name"`
	ChainID      string        `json:"chain_id"`
	Bech32Prefix string        `json:"bech32_prefix"`
	Fees         FeeTokens     `json:"fees"`
	Staking      StakingTokens `json:"staking"`
	Apis         Apis          `json:"apis"`
	Explorers    []Explorer    `json:"explorers"`
}

// AssetList struct from chain registry
type DenomUnit struct {
	Denom    string `json:"denom"`
	Exponent int64  `json:"exponent"`
}

type Asset struct {
	Description string      `json:"description"`
	DenomUnits  []DenomUnit `json:"denom_units"`
	Base        string      `json:"base"`
	Name        string      `json:"name"`
	Display     string      `json:"display"`
	Symbol      string      `json:"symbol"`
	LogoUris    struct {
		Svg string `json:"svg"`
		Png string `json:"png"`
	} `json:"logo_URIs"`
	CoinGeckoID string `json:"coingecko_id"`
}

type AssetList struct {
	ChainName string  `json:"chain_name"`
	Assets    []Asset `json:"assets"`
}
