package types

// Utils
type Pagination struct {
	NextKey interface{} `json:"next_key"`
	Total   string      `json:"total"`
}

// Validators
type ValidatorAPIResponse struct {
	Validators []Validator `json:"validators"`
	Pagination Pagination  `json:"pagination"`
}

type ConsensusKey struct {
	TypeURL string `json:"type_url"`
	Value   string `json:"value"`
}

type Description struct {
	Moniker         string `json:"moniker"`
	Identity        string `json:"identity"`
	Website         string `json:"website"`
	SecurityContact string `json:"security_contact"`
	Details         string `json:"details"`
}

type CommissionRate struct {
	Rate          string `json:"rate"`
	MaxRate       string `json:"max_rate"`
	MaxChangeRate string `json:"max_change_rate"`
}

type Commission struct {
	CommissionRate CommissionRate `json:"commission_rates"`
	UpdateTime     string         `json:"update_time"`
}

type Validator struct {
	OperatorAddress   string       `json:"operator_address"`
	ConsensusKey      ConsensusKey `json:"consensus_pubkey"`
	Jailed            bool         `json:"jailed"`
	Status            string       `json:"status"`
	Tokens            string       `json:"tokens"`
	DelegatorShares   string       `json:"delegator_shares"`
	Description       Description  `json:"description"`
	UnbondingHeight   string       `json:"unbonding_height"`
	UnbondingTime     string       `json:"unbonding_time"`
	Commission        Commission   `json:"commission"`
	MinSelfDelegation string       `json:"min_self_delegation"`
	Rank              int          `json:"rank"`
}

// Accounts
type PubKeyAccount struct {
	Type string `json:"@type"`
	Key  string `json:"string"`
}

type BaseAccount struct {
	Address       string        `json:"address"`
	PubKey        PubKeyAccount `json:"pub_key"`
	AccountNumber string        `json:"account_number"`
	Sequence      string        `json:"sequence"`
}

type BaseAccountDetails struct {
	Type        string      `json:"@type"`
	BaseAccount BaseAccount `json:"base_account"`
	CodeHash    string      `json:"code_hash"`
}

type BaseAccountResponse struct {
	Account BaseAccountDetails `json:"account"`
}

type AccountDetails struct {
	Type          string        `json:"@type"`
	Address       string        `json:"address"`
	PubKey        PubKeyAccount `json:"pub_key"`
	AccountNumber string        `json:"account_number"`
	Sequence      string        `json:"sequence"`
}

type AccountResponse struct {
	Account AccountDetails `json:"account"`
}

// Balances
type Balance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type BalancesResponse struct {
	Balances   []Balance  `json:"balances"`
	Pagination Pagination `json:"pagination"`
}

// Delegations
type DelegationResponses struct {
	DelegationResponses []DelegationResponse `json:"delegation_responses"`
	Pagination          Pagination           `json:"pagination"`
}

type DelegationResponse struct {
	Delegation struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
	} `json:"delegation"`
	Balance struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balance"`
}

// Rewards
type Reward struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Rewards struct {
	ValidatorAddress string   `json:"validator_address"`
	Reward           []Reward `json:"reward"`
}

type RewardsResponse struct {
	Rewards []Rewards `json:"rewards"`
	Total   []Reward  `json:"total"`
}

// Unbonding
type UnbondingEntry struct {
	CreationHeight string `json:"creation_height"`
	CompletionTime string `json:"completion_time"`
	InitialBalance string `json:"initial_balance"`
	Balance        string `json:"balance"`
}

type Unbonding struct {
	DelegatorAddress string           `json:"delegator_addresss"`
	ValidatorAddress string           `json:"validator_address"`
	Entries          []UnbondingEntry `json:"entries"`
}

type UnbondingResponse struct {
	UnbondingResponses []Unbonding `json:"unbonding_responses"`
	Pagination         Pagination  `json:"pagination"`
}

// Proposals
type TallyResponse struct {
	Tally FinalTallyResult `json:"tally"`
}

type FinalTallyResult struct {
	Yes        string `json:"yes"`
	No         string `json:"no"`
	Abstain    string `json:"abstain"`
	NoWithVeto string `json:"no_with_veto"`
}

type ProposalContent struct {
	Type        string `json:"@type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Recipient   string `json:"recipient"`
	Amount      []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"amount"`
	Changes []struct {
		Subspace string `json:"subspace"`
		Key      string `json:"key"`
		Value    string `json:"value"`
	} `json:"changes"`
}

type Proposal struct {
	ProposalID       string           `json:"proposal_id"`
	Content          ProposalContent  `json:"content"`
	Status           string           `json:"status"`
	FinalTallyResult FinalTallyResult `json:"final_tally_result"`
	SubmitTime       string           `json:"submit_time"`
	DepositEndTime   string           `json:"deposit_end_time"`
	TotalDeposit     []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"total_deposit"`
	VotingStartTime string `json:"voting_start_time"`
	VotingEndTime   string `json:"voting_end_time"`
}

type ProposalsResponse struct {
	Proposals []Proposal `json:"proposals"`
}

// Transactions
type Message struct {
	TypeURL string `json:"@type"`
}

type Body struct {
	Messages []Message `json:"messages"`
}

type Tx struct {
	Body Body `json:"body"`
}

type TxResponse struct {
	Height string `json:"height"`
	TxHash string `json:"txhash"`
	Code   int64  `json:"code"`
	Tx     Tx     `json:"tx"`
}

type TransactionResponse struct {
	Txs         []Tx         `json:"txs"`
	TxResponses []TxResponse `json:"tx_responses"`
	Pagination  Pagination   `json:"pagination"`
}
