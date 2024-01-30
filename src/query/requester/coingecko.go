package requester

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MarketData struct {
	CurrentPrice struct {
		Usd float64 `json:"usd"`
	} `json:"current_price"`
}

type HistoryResponse struct {
	MarketData MarketData `json:"market_data"`
}

type Price struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

func QueryCoinGeckoHistoricPrices(token string) ([]Price, error) {
	dates := GenerateLast8Days()
	ret := []Price{}
	for _, v := range dates {
		url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/history?date=%s&localization=false", token, v)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return ret, err
		}

		resp, err := Client.Do(req)
		if err != nil {
			return ret, err
		}

		if resp.StatusCode != 200 {
			return ret, fmt.Errorf("coingecko response status code different from 200: %d", resp.StatusCode)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		if err != nil || len(string(body)) == 0 {
			return ret, err
		}

		var m HistoryResponse
		err = json.Unmarshal(body, &m)
		if err != nil {
			return ret, err
		}
		ret = append(ret, Price{Date: v, Price: m.MarketData.CurrentPrice.Usd})
	}

	if ret[len(ret)-1].Price == 0 {
		return ret[0 : len(ret)-1], nil
	}

	return ret[1:], nil
}
