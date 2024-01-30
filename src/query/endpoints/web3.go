package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"
)

var wwg sync.WaitGroup

func PingWeb3(endpoint string, c chan Endpoint) {
	defer wwg.Done()

	url := fmt.Sprintf("%scosmos/base/tendermint/v1beta1/blocks/latest", endpoint)

	// record start time to measure latency
	start := time.Now()

	// make request
	payload := bytes.NewBuffer([]byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`))
	resp, err := Client.Post(url, "application/json", payload)
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

	var jsonRes Web3Response
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

	if len(jsonRes.Result) == 0 {
		e := Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	convertedHeight, err := strconv.ParseInt(jsonRes.Result[2:], 16, 64)
	if err != nil {
		e := Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	e := Endpoint{
		URL:     endpoint,
		Latency: duration,
		Height:  int(convertedHeight),
	}

	c <- e
}

func ProcessWeb3(web3Endpoints []string) []Endpoint {
	// create a channel to receive results for each web3 endpoint
	web3Channel := make(chan Endpoint, len(web3Endpoints))

	for _, v := range web3Endpoints {
		// ping web3 endpoint & get results in a goroutine
		go PingWeb3(v, web3Channel)
		// add goroutine to web3 wait group
		wwg.Add(1)
	}

	web3Results := make([]Endpoint, 0)

	done := make(chan struct{})

	// Loop over values sent via channel.
	// This has to be as a separate goroutine in order to keep channel listener open
	// while waiting for all endpoints to be pinged & processed
	go func() {
		for r := range web3Channel {
			web3Results = append(web3Results, r)
		}
		close(done)
	}()
	// wait for all endpoints to be pinged & processed
	wwg.Wait()
	// close channel
	close(web3Channel)
	// wait for all values to be read
	<-done

	sorted := SortEndpoints(web3Results)

	return sorted
}
