package endpoints

import (
	"net/http"
	"sort"
	"time"
)

var Client = http.Client{
	Timeout: 2 * time.Second,
}

// SortEndpoints sort them by height and latency
func SortEndpoints(array []Endpoint) []Endpoint {
	sort.Slice(array, func(i, j int) bool {
		if array[i].Height == -1 {
			return false
		}

		if array[i].Height != array[j].Height {
			return array[i].Height > array[j].Height
		}

		if array[i].Latency == -1 {
			return false
		}
		return array[i].Latency < array[j].Latency
	})

	return array
}
