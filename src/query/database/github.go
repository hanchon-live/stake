package database

import (
	"encoding/json"
	"strings"

	"github.com/hanchon-live/stake/src/query/requester"
	"github.com/hanchon-live/stake/src/query/types"
)

func (db Database) GetChain(chain string) (types.Chain, error) {
	chain = strings.ToLower(chain)
	key := chain + "chaininfo"
	item := db.Cache.Get(key)
	if item != nil {
		var m types.Chain
		err := json.Unmarshal([]byte(item.Value()), &m)
		if err == nil {
			return m, nil
		}
	}

	chainInfo, err := requester.GetChain(chain)
	if err != nil {
		return types.Chain{}, err
	}

	chainInfoAsString, err := json.Marshal(chainInfo)
	if err != nil {
		return types.Chain{}, err
	}

	db.Cache.Set(key, string(chainInfoAsString), TimeoutLong)
	return chainInfo, nil
}

func (db Database) GetAsset(chain string) (types.AssetList, error) {
	chain = strings.ToLower(chain)
	key := chain + "chaininfo"
	item := db.Cache.Get(key)
	if item != nil {
		var m types.AssetList
		err := json.Unmarshal([]byte(item.Value()), &m)
		if err == nil {
			return m, nil
		}
	}

	assetInfo, err := requester.GetAsset(chain)
	if err != nil {
		return types.AssetList{}, err
	}

	assetInfoAsString, err := json.Marshal(assetInfo)
	if err != nil {
		return types.AssetList{}, err
	}

	db.Cache.Set(key, string(assetInfoAsString), TimeoutLong)
	return assetInfo, nil
}
