package database

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/hanchon-live/stake/src/query/requester"
	"github.com/hanchon-live/stake/src/query/types"
)

// Sort validators by total voting power
func powerSort(array []types.Validator) []types.Validator {
	sort.Slice(array, func(i, j int) bool {
		a := new(big.Int)
		a, ok := a.SetString(array[i].Tokens, 10)
		if !ok {
			log.Default().Println("Error converting tokens: ", array[i].Tokens)
			a = big.NewInt(0)
		}

		b := new(big.Int)
		b, ok = b.SetString(array[j].Tokens, 10)
		if !ok {
			log.Default().Println("Error converting tokens: ", array[j].Tokens)
			b = big.NewInt(0)
		}

		return a.Cmp(b) == 1
	})

	return array
}

// All the chain validators
func (db Database) GetValidators(chain string) ([]types.Validator, error) {
	chain = strings.ToLower(chain)

	key := chain + "validators"
	item := db.Cache.Get(key)
	if item != nil {
		var m []types.Validator
		err := json.Unmarshal([]byte(item.Value()), &m)
		if err == nil {
			return m, nil
		}
	}

	// Get best endpoint
	endpoint, err := db.GetRestEndpoint(chain)
	if err != nil {
		return []types.Validator{}, err
	}

	// Validators query
	query := "cosmos/staking/v1beta1/validators?pagination.limit=1000"

	response, err := requester.MakeGetRequest(endpoint, query)
	if err != nil {
		return []types.Validator{}, err
	}

	var validatorsObject types.ValidatorAPIResponse
	err = json.Unmarshal([]byte(response), &validatorsObject)

	if err != nil {
		return []types.Validator{}, err
	}

	sortedValidators := powerSort(validatorsObject.Validators)

	sortedValidatorsAsString, err := json.Marshal(sortedValidators)
	if err != nil {
		return []types.Validator{}, err
	}

	db.Cache.Set(key, string(sortedValidatorsAsString), Timeout15sec)

	return sortedValidators, nil
}

// Account info from the chain
func (db Database) getAccount(address string, chain string) (string, error) {
	key := chain + "account" + address
	item := db.Cache.Get(key)
	if item != nil {
		return item.Value(), nil
	}

	// Get best endpoint
	endpoint, err := db.GetRestEndpoint(chain)
	if err != nil {
		return "", err
	}

	// Validators query
	query := "cosmos/auth/v1beta1/accounts/" + address

	response, err := requester.MakeGetRequest(endpoint, query)
	if err != nil {
		return "", err
	}

	db.Cache.Set(key, response, TimeoutOneBlock)

	return response, nil
}

// Account info to be used
type AccountInfo struct {
	AccountNumber uint64 `json:"account_number"`
	Sequence      uint64 `json:"sequence"`
}

func (db Database) GetAccountInfo(chain string, sender string) (AccountInfo, error) {
	chain = strings.ToLower(chain)
	isEthermintNetwork := false
	switch chain {
	case "evmos":
		isEthermintNetwork = true
	case "planq":
		isEthermintNetwork = true
	}

	val, err := db.getAccount(sender, chain)
	if err != nil {
		return AccountInfo{}, fmt.Errorf("error while getting account details, please try again")
	}

	// TODO: vesting accounts are not being considered here
	if isEthermintNetwork {
		var accountDetails types.BaseAccountResponse
		err = json.Unmarshal([]byte(val), &accountDetails)
		if err != nil {
			return AccountInfo{}, err
		}
		accountNumber, err := strconv.ParseUint(accountDetails.Account.BaseAccount.AccountNumber, 10, 64)
		if err != nil {
			return AccountInfo{}, err
		}
		sequence, err := strconv.ParseUint(accountDetails.Account.BaseAccount.Sequence, 10, 64)
		if err != nil {
			return AccountInfo{}, err
		}
		return AccountInfo{AccountNumber: accountNumber, Sequence: sequence}, nil
	}

	// Classic cosmos network
	var accountDetails types.AccountResponse
	err = json.Unmarshal([]byte(val), &accountDetails)
	if err != nil {
		return AccountInfo{}, err
	}

	accountNumber, err := strconv.ParseUint(accountDetails.Account.AccountNumber, 10, 64)
	if err != nil {
		return AccountInfo{}, err
	}

	sequence, err := strconv.ParseUint(accountDetails.Account.Sequence, 10, 64)
	if err != nil {
		return AccountInfo{}, err
	}
	return AccountInfo{AccountNumber: accountNumber, Sequence: sequence}, nil
}

// Balances
func (db Database) GetAccountBalance(chain string, address string) (types.BalancesResponse, error) {
	chain = strings.ToLower(chain)

	key := chain + "accountbalance" + address
	item := db.Cache.Get(key)
	if item != nil {
		var m types.BalancesResponse
		err := json.Unmarshal([]byte(item.Value()), &m)
		if err == nil {
			return m, nil
		}
	}

	// Get best endpoint
	endpoint, err := db.GetRestEndpoint(chain)
	if err != nil {
		return types.BalancesResponse{}, err
	}

	// Validators query
	query := "cosmos/bank/v1beta1/balances/" + address

	response, err := requester.MakeGetRequest(endpoint, query)
	if err != nil {
		return types.BalancesResponse{}, err
	}

	var m types.BalancesResponse
	err = json.Unmarshal([]byte(response), &m)
	if err != nil {
		return types.BalancesResponse{}, err
	}

	db.Cache.Set(key, response, Timeout15sec)

	return m, nil
}

// Governance
// Delegations
// nolint
func (db Database) GetAccountDelegations(chain string, address string) (types.DelegationResponses, error) {
	chain = strings.ToLower(chain)

	key := chain + "accountdelegations" + address
	item := db.Cache.Get(key)
	if item != nil {
		var m types.DelegationResponses
		err := json.Unmarshal([]byte(item.Value()), &m)
		if err == nil {
			return m, nil
		}
	}

	// Get best endpoint
	endpoint, err := db.GetRestEndpoint(chain)
	if err != nil {
		return types.DelegationResponses{}, err
	}

	// Validators query
	query := "cosmos/staking/v1beta1/delegations/" + address + "?pagination.limit=200"

	response, err := requester.MakeGetRequest(endpoint, query)
	if err != nil {
		return types.DelegationResponses{}, err
	}

	var m types.DelegationResponses
	err = json.Unmarshal([]byte(response), &m)
	if err != nil {
		return types.DelegationResponses{}, err
	}

	db.Cache.Set(key, response, Timeout15sec)

	return m, nil
}

// Rewards
// nolint
func (db Database) GetAccountRewards(chain string, address string) (types.RewardsResponse, error) {
	chain = strings.ToLower(chain)

	key := chain + "accountrewards" + address
	item := db.Cache.Get(key)
	if item != nil {
		var m types.RewardsResponse
		err := json.Unmarshal([]byte(item.Value()), &m)
		if err == nil {
			return m, nil
		}
	}

	// Get best endpoint
	endpoint, err := db.GetRestEndpoint(chain)
	if err != nil {
		return types.RewardsResponse{}, err
	}

	// Validators query
	query := "cosmos/distribution/v1beta1/delegators/" + address + "/rewards"

	response, err := requester.MakeGetRequest(endpoint, query)
	if err != nil {
		return types.RewardsResponse{}, err
	}

	var m types.RewardsResponse
	err = json.Unmarshal([]byte(response), &m)
	if err != nil {
		return types.RewardsResponse{}, err
	}

	db.Cache.Set(key, response, Timeout15sec)

	return m, nil
}

// Unbonding
// nolint
func (db Database) GetAccountUnbonding(chain string, address string) (types.UnbondingResponse, error) {
	chain = strings.ToLower(chain)

	key := chain + "accountunbonding" + address
	item := db.Cache.Get(key)
	if item != nil {
		var m types.UnbondingResponse
		err := json.Unmarshal([]byte(item.Value()), &m)
		if err == nil {
			return m, nil
		}
	}

	// Get best endpoint
	endpoint, err := db.GetRestEndpoint(chain)
	if err != nil {
		return types.UnbondingResponse{}, err
	}

	// Validators query
	query := "cosmos/staking/v1beta1/delegators/" + address + "/unbonding_delegations"

	response, err := requester.MakeGetRequest(endpoint, query)
	if err != nil {
		return types.UnbondingResponse{}, err
	}

	var m types.UnbondingResponse
	err = json.Unmarshal([]byte(response), &m)
	if err != nil {
		return types.UnbondingResponse{}, err
	}

	db.Cache.Set(key, response, Timeout15sec)

	return m, nil
}

// Tally
// nolint
func (db Database) GetProposalsTally(chain string, proposalId string) (types.TallyResponse, error) {
	chain = strings.ToLower(chain)
	key := chain + "proposaltally" + proposalId
	item := db.Cache.Get(key)
	if item != nil {
		var tally types.TallyResponse
		err := json.Unmarshal([]byte(item.Value()), &tally)
		if err == nil {
			return tally, nil
		}
	}

	// Get best endpoint
	endpoint, err := db.GetRestEndpoint(chain)
	if err != nil {
		return types.TallyResponse{}, err
	}

	// Query
	query := "cosmos/gov/v1beta1/proposals/" + proposalId + "/tally"

	response, err := requester.MakeGetRequest(endpoint, query)
	if err != nil {
		return types.TallyResponse{}, err
	}

	var m types.TallyResponse
	err = json.Unmarshal([]byte(response), &m)
	if err != nil {
		return types.TallyResponse{}, err
	}

	db.Cache.Set(key, response, Timeout15sec)

	return m, nil
}

// Last 10 proposals
func (db Database) GetProposals(chain string) (types.ProposalsResponse, error) {
	chain = strings.ToLower(chain)

	key := chain + "proposals"
	item := db.Cache.Get(key)
	if item != nil {
		var m types.ProposalsResponse
		err := json.Unmarshal([]byte(item.Value()), &m)
		if err == nil {
			return m, nil
		}
	}

	// Get best endpoint
	endpoint, err := db.GetRestEndpoint(chain)
	if err != nil {
		return types.ProposalsResponse{}, err
	}

	// Query
	query := "cosmos/gov/v1beta1/proposals?pagination.limit=10&pagination.reverse=true"
	response, err := requester.MakeGetRequest(endpoint, query)
	if err != nil {
		return types.ProposalsResponse{}, err
	}

	var m types.ProposalsResponse
	err = json.Unmarshal([]byte(response), &m)
	if err != nil {
		return types.ProposalsResponse{}, err
	}

	// Get tally
	const ProposalStatusVotingPeriod = "PROPOSAL_STATUS_VOTING_PERIOD"
	for k := range m.Proposals {
		if m.Proposals[k].Status == ProposalStatusVotingPeriod {
			value, err := db.GetProposalsTally(chain, m.Proposals[k].ProposalID)
			if err != nil {
				return types.ProposalsResponse{}, err
			}
			m.Proposals[k].FinalTallyResult = value.Tally
		}
	}

	responseWithTally, err := json.Marshal(m)
	if err != nil {
		return types.ProposalsResponse{}, err
	}

	db.Cache.Set(key, string(responseWithTally), Timeout30min)

	return m, nil
}

// Last transactions by wallet
// Transactions
type Tx struct {
	Type             string `json:"type"`
	Hash             string `json:"hash"`
	Code             int64  `json:"code"`
	Height           string `json:"height"`
	Message          string `json:"message"`
	AmountOfMessages int    `json:"amount"`
}

type Txns struct {
	Txns []Tx `json:"txns"`
}

func (db Database) GetAccountTransactions(chain string, address string) (Txns, error) {
	res := Txns{
		Txns: []Tx{},
	}
	uniqueHash := make(map[string]bool)
	chain = strings.ToLower(chain)

	key := chain + "accounttransactions" + address
	item := db.Cache.Get(key)
	if item != nil {
		var m Txns
		err := json.Unmarshal([]byte(item.Value()), &m)
		if err == nil {
			return m, nil
		}
	}

	// Get best endpoint
	endpoint, err := db.GetRestEndpoint(chain)
	if err != nil {
		return Txns{}, err
	}

	// Recipient query
	queryRecipient := "cosmos/tx/v1beta1/txs?events=transfer.recipient='" + address + "'&order_by=2"
	responseRecipient, err := requester.MakeGetRequest(endpoint, queryRecipient)
	if err != nil {
		return Txns{}, err
	}

	var m types.TransactionResponse
	err = json.Unmarshal([]byte(responseRecipient), &m)
	if err != nil {
		return Txns{}, err
	}
	if len(m.Txs) != len(m.TxResponses) {
		return Txns{}, fmt.Errorf("tx and responses have different len")
	}
	for k := range m.Txs {
		if _, value := uniqueHash[m.TxResponses[k].TxHash]; !value {
			uniqueHash[m.TxResponses[k].TxHash] = true
			res.Txns = append(res.Txns, Tx{
				Type:             "recipient",
				Hash:             m.TxResponses[k].TxHash,
				Code:             m.TxResponses[k].Code,
				Height:           m.TxResponses[k].Height,
				Message:          m.Txs[k].Body.Messages[0].TypeURL,
				AmountOfMessages: len(m.Txs[k].Body.Messages),
			})
		}
		if k > 30 {
			break
		}
	}

	// Sender query
	querySender := "cosmos/tx/v1beta1/txs?events=message.sender='" + address + "'&order_by=2"
	responseSender, err := requester.MakeGetRequest(endpoint, querySender)
	if err != nil {
		return Txns{}, err
	}

	err = json.Unmarshal([]byte(responseSender), &m)
	if err != nil {
		return Txns{}, err
	}
	if len(m.Txs) != len(m.TxResponses) {
		return Txns{}, fmt.Errorf("tx and responses have different len")
	}
	for k := range m.Txs {
		if _, value := uniqueHash[m.TxResponses[k].TxHash]; !value {
			uniqueHash[m.TxResponses[k].TxHash] = true
			res.Txns = append(res.Txns, Tx{
				Type:             "sender",
				Hash:             m.TxResponses[k].TxHash,
				Code:             m.TxResponses[k].Code,
				Height:           m.TxResponses[k].Height,
				Message:          m.Txs[k].Body.Messages[0].TypeURL,
				AmountOfMessages: len(m.Txs[k].Body.Messages),
			})
		}
		if k > 30 {
			break
		}
	}

	sort.Slice(res.Txns, func(i, j int) bool {
		return res.Txns[i].Height < res.Txns[j].Height
	})

	responseAsString, err := json.Marshal(res)
	if err != nil {
		return Txns{}, err
	}

	db.Cache.Set(key, string(responseAsString), Timeout15sec)

	return res, nil
}
