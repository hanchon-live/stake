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

var rwg sync.WaitGroup

func PingNonTendermintRest(endpoint string, c chan Endpoint) {
	defer rwg.Done()

	url := fmt.Sprintf("%scosmos/auth/v1beta1/params", endpoint)

	// record start time to measure latency
	start := time.Now()

	// make request
	_, err := Client.Get(url)
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

	e := Endpoint{
		URL:     endpoint,
		Latency: duration,
		Height:  0,
	}

	c <- e
}

func PingRest(endpoint string, c chan Endpoint) {
	defer rwg.Done()

	url := fmt.Sprintf("%scosmos/base/tendermint/v1beta1/blocks/latest", endpoint)

	// record start time to measure latency
	start := time.Now()

	// make request
	resp, err := Client.Get(url)
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

	// get block height from response
	height := -1
	var jsonRes RestResponse
	err = json.Unmarshal(body, &jsonRes)
	if err == nil {
		height, err = strconv.Atoi(jsonRes.Block.Header.Height)
		if err != nil {
			height = -1
		}
	}

	// Check if the tx index is enabled
	transactionURL := fmt.Sprintf("%scosmos/tx/v1beta1/txs?events=transfer.recipient='osmo1pmk2r32ssqwps42y3c9d4clqlca403yd05x9ye'&pagination.offset=0&pagination.limit=30", endpoint)
	// make request
	resp, err = Client.Get(transactionURL)
	if err != nil {
		e := Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

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

	if strings.Contains(string(body), constants.IndexingDisabledError) {
		e := Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}
	// End tx index validation

	e := Endpoint{
		URL:     endpoint,
		Latency: duration,
		Height:  height,
	}

	c <- e
}

func ProcessRest(restEndpoints []string, standartEndpoint bool) []Endpoint {
	// Standart endpoint is only false for gravity chain

	// create a channel to receive results for each rest endpoint
	restChannel := make(chan Endpoint, len(restEndpoints))

	for _, v := range restEndpoints {
		// ping REST endpoint & get results in a goroutine
		if standartEndpoint {
			go PingRest(v, restChannel)
		} else {
			go PingNonTendermintRest(v, restChannel)
		}
		// add goroutine to rest wait group
		rwg.Add(1)
	}

	restResults := make([]Endpoint, 0)

	done := make(chan struct{})

	// Loop over values sent via channel.
	// This has to be as a separate goroutine in order to keep channel listener open
	// while waiting for all endpoints to be pinged & processed
	go func() {
		for r := range restChannel {
			restResults = append(restResults, r)
		}
		close(done)
	}()
	// wait for all endpoints to be pinged & processed
	rwg.Wait()
	// close channel
	close(restChannel)
	// wait for all values to be read
	<-done

	sorted := SortEndpoints(restResults)
	return sorted
}
