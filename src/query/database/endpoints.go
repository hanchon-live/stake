package database

import (
	"fmt"
	"strings"

	endpointFunctions "github.com/hanchon-live/stake/src/query/endpoints"
)

func (db Database) GetWeb3Endpoint(chain string) (string, error) {
	chain = strings.ToLower(chain)
	key := chain + "web3endpoint"
	item := db.Cache.Get(key)
	if item != nil {
		return item.Value(), nil
	}

	chainInfo, err := db.GetChain(chain)
	if err != nil {
		return "", err
	}

	endpoints := []string{}
	for _, v := range chainInfo.Apis.Evm {
		endpoint := v.Address
		if endpoint[len(endpoint)-1:] != "/" {
			endpoint += "/"
		}
		endpoints = append(endpoints, endpoint)
	}

	endpointsSorted := endpointFunctions.ProcessWeb3(endpoints)
	if len(endpointsSorted) > 0 {
		db.Cache.Set(key, endpointsSorted[0].URL, Timeout30min)
		return endpointsSorted[0].URL, nil
	}
	return "", fmt.Errorf("not found")
}

func (db Database) GetJrpcEndpoint(chain string) (string, error) {
	chain = strings.ToLower(chain)
	key := chain + "jrpcendpoint"
	item := db.Cache.Get(key)
	if item != nil {
		return item.Value(), nil
	}

	chainInfo, err := db.GetChain(chain)
	if err != nil {
		return "", err
	}

	endpoints := []string{}
	for _, v := range chainInfo.Apis.RPC {
		endpoint := v.Address
		if endpoint[len(endpoint)-1:] != "/" {
			endpoint += "/"
		}
		endpoints = append(endpoints, endpoint)
	}
	endpointsSorted := endpointFunctions.ProcessJrpc(endpoints)

	if len(endpointsSorted) > 0 {
		db.Cache.Set(key, endpointsSorted[0].URL, Timeout30min)
		return endpointsSorted[0].URL, nil
	}
	return "", fmt.Errorf("not found")
}

func (db Database) GetRestEndpoint(chain string) (string, error) {
	standartEndpoint := true
	chain = strings.ToLower(chain)
	if chain == "gravity" {
		standartEndpoint = false
	}

	key := chain + "restendpoint"
	item := db.Cache.Get(key)
	if item != nil {
		return item.Value(), nil
	}

	chainInfo, err := db.GetChain(chain)
	if err != nil {
		return "", err
	}

	endpoints := []string{}
	for _, v := range chainInfo.Apis.Rest {
		endpoint := v.Address
		if endpoint[len(endpoint)-1:] != "/" {
			endpoint += "/"
		}
		endpoints = append(endpoints, endpoint)
	}
	endpointsSorted := endpointFunctions.ProcessRest(endpoints, standartEndpoint)
	if len(endpointsSorted) > 0 {
		// TODO: Make this timeout 15 seconds and add a gorutine to the main to keep it up to date
		db.Cache.Set(key, endpointsSorted[0].URL, Timeout30min)
		return endpointsSorted[0].URL, nil
	}
	return "", fmt.Errorf("not found")
}
