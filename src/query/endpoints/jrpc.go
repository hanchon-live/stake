package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hanchon-live/stake/src/query/constants"
)

var jwg sync.WaitGroup

func PingJrpc(endpoint string, c chan Endpoint) {
	defer jwg.Done()

	transactionURL := fmt.Sprintf("%stx?hash=0x0000000000000000000000000000000000000000000000000000000000000000", endpoint)

	// make request
	resp, err := Client.Get(transactionURL)
	if err != nil {
		e := Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil || len(string(body)) == 0 {
		e := Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	var jsonTransactionRes JrpcTransactionErrorResponse
	err = json.Unmarshal(body, &jsonTransactionRes)
	if err != nil {
		e := Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	if strings.Contains(jsonTransactionRes.Error.Data, constants.IndexingDisabledError) {
		e := Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	// record start time to measure latency
	start := time.Now()

	url := fmt.Sprintf("%sstatus", endpoint)

	// make request
	resp, err = Client.Get(url)

	if err != nil {
		e := Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	// compute latency
	duration := time.Since(start).Seconds()

	body, err = io.ReadAll(resp.Body)

	if err != nil || len(string(body)) == 0 {
		e := Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	var jsonRes JrpcStatusResponse
	err = json.Unmarshal(body, &jsonRes)
	if err != nil {
		e := Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	if jsonRes.Result.NodeInfo.Other.TxIndex == "on" {

		height, err := strconv.Atoi(jsonRes.Result.SyncInfo.LatestBlockHeight)
		if err != nil {
			height = -1
		}

		e := Endpoint{
			URL:     endpoint,
			Latency: duration,
			Height:  height,
		}

		c <- e

	} else {
		e := Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
	}
}

func ProcessJrpc(jrpcEndpoints []string) []Endpoint {
	// create a channel to receive results for each jrpc endpoint
	jrpcChannel := make(chan Endpoint, len(jrpcEndpoints))

	for _, v := range jrpcEndpoints {
		// ping jrpc endpoint & get results in a goroutine
		go PingJrpc(v, jrpcChannel)
		// add goroutine to jrpc wait group
		jwg.Add(1)
	}

	jrpcResults := make([]Endpoint, 0)

	done := make(chan struct{})

	// Loop over values sent via channel.
	// This has to be as a separate goroutine in order to keep channel listener open
	// while waiting for all endpoints to be pinged & processed
	go func() {
		for r := range jrpcChannel {
			jrpcResults = append(jrpcResults, r)
		}
		close(done)
	}()
	// wait for all endpoints to be pinged & processed
	jwg.Wait()
	// close channel
	close(jrpcChannel)
	// wait for all values to be read
	<-done

	sorted := SortEndpoints(jrpcResults)

	return sorted
}
