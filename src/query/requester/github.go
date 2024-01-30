package requester

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hanchon-live/stake/src/query/constants"
	"github.com/hanchon-live/stake/src/query/types"
)

func QueryGithubWithCache(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("authorization", "token "+constants.GithubKey)

	resp, err := Client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("github response status code different from 200: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil || len(string(body)) == 0 {
		return "", err
	}

	bodyString := string(body)
	return bodyString, nil
}

func GetChain(chain string) (types.Chain, error) {
	res, err := GetJsonsFromFolder(constants.ChainRegistryURL, chain+"/chain.json")
	if err != nil {
		return types.Chain{}, err
	}
	var m types.Chain
	err = json.Unmarshal([]byte(res.Content), &m)
	if err != nil {
		return types.Chain{}, err
	}
	return m, nil
}

func GetAsset(chain string) (types.AssetList, error) {
	res, err := GetJsonsFromFolder(constants.ChainRegistryURL, chain+"/assetlist.json")
	if err != nil {
		return types.AssetList{}, err
	}
	var m types.AssetList
	err = json.Unmarshal([]byte(res.Content), &m)
	if err != nil {
		return types.AssetList{}, err
	}
	return m, nil
}

func GetJsonsFromFolder(url string, folder string) (types.File, error) {
	// TODO: we can use t.sha to compare if the value is up to date
	apiResp, err := QueryGithubWithCache(url)
	if err != nil {
		return types.File{}, err
	}

	var m types.TreeResponse
	err = json.Unmarshal([]byte(apiResp), &m)
	if err != nil {
		return types.File{}, err
	}

	for _, t := range m.Tree {
		if t.Mode == "100644" {
			// Is file
			if strings.Contains(t.Path, folder) {
				fileResponse, err := QueryGithubWithCache(t.URL)
				if err != nil {
					return types.File{}, err
				}

				var m types.Content
				err = json.Unmarshal([]byte(fileResponse), &m)
				if err == nil {
					rawDecodedText, err := base64.StdEncoding.DecodeString(m.Content)
					if err != nil {
						return types.File{}, err
					}
					return types.File{Content: string(rawDecodedText), URL: t.Path}, nil
				}
			}
		}
	}
	return types.File{}, fmt.Errorf("element not found")
}
