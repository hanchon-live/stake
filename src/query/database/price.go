package database

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/hanchon-live/stake/src/query/requester"
)

func (db Database) GetPrices(token string) (string, error) {
	token = strings.ToLower(token)
	key := token + "historic"
	item := db.Cache.Get(key)
	if item != nil {
		return item.Value(), nil
	}

	prices, err := requester.QueryCoinGeckoHistoricPrices(token)
	if err != nil {
		return "", err
	}

	chainInfoAsString, err := json.Marshal(prices)
	if err != nil {
		return "", err
	}
	ret := string(chainInfoAsString)

	db.Cache.Set(key, ret, 15*time.Minute)
	return ret, nil
}
