package bing

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/bevly/bevly/httpagent"
	"github.com/bevly/bevly/websearch"
)

const SearchBaseURL = "https://api.datamarket.azure.com/Bing/Search/Web"

var ErrBadJson = errors.New("unexpected json structure")

type BingSearch struct {
	BaseURL     string
	ApiKey      string
	ResultCount int
	agent       *httpagent.Agent
	auth        string
}

func DefaultSearch() websearch.Search {
	bing := &BingSearch{
		BaseURL:     SearchBaseURL,
		ApiKey:      DefaultApiKey(),
		ResultCount: 10,
		agent:       httpagent.New(),
	}
	bing.init()
	return bing
}

func SearchWithURLKey(url, apiKey string) websearch.Search {
	bing := &BingSearch{
		BaseURL:     url,
		ApiKey:      apiKey,
		ResultCount: 10,
		agent:       httpagent.New(),
	}
	bing.init()
	return bing
}

func (b *BingSearch) init() {
	b.auth = "Basic " +
		base64.StdEncoding.EncodeToString([]byte(b.ApiKey+":"+b.ApiKey))
}

func (b *BingSearch) Search(terms string) ([]websearch.Result, error) {
	searchURL := b.SearchURL(terms)

	log.Printf("BingSearch(%s): GET %s\n", terms, searchURL)
	parsedURL, err := url.Parse(searchURL)
	if err != nil {
		return nil, err
	}

	request := &http.Request{
		URL: parsedURL,
		Header: http.Header{
			"User-Agent":    {b.agent.UserAgent},
			"Authorization": {b.auth},
		},
	}
	res, err := b.agent.Client.Do(request)
	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	return b.parseResult(res)
}

type bingResponse struct {
	Body bingContainer `json:"d"`
}

type bingContainer struct {
	Results []bingResult `json:"results"`
}

type bingResult struct {
	URL  string `json:"Url"`
	Text string `json:"Title"`
}

func (b *BingSearch) parseResult(res *http.Response) ([]websearch.Result, error) {
	if res.StatusCode >= 300 || res.StatusCode < 200 {
		return nil, fmt.Errorf("bing search http err:%d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var bingRes bingResponse
	if err = json.Unmarshal(body, &bingRes); err != nil {
		return nil, fmt.Errorf("malformed bing response: %s (%s)",
			string(body), err)
	}

	results := bingRes.Body.Results
	webresults := make([]websearch.Result, 0, len(results))
	for _, res := range results {
		webresults = append(webresults, websearch.Result{
			URL:  res.URL,
			Text: res.Text,
		})
	}
	return webresults, nil
}

func (b *BingSearch) SearchURL(terms string) string {
	return b.BaseURL + "?" +
		url.Values{
			"$format": {"json"},
			"$top":    {strconv.Itoa(b.ResultCount)},
			"Query":   {"'{" + terms + "}'"},
		}.Encode()
}

func DefaultApiKey() string {
	file := os.Getenv("BING_API_KEY")
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(fmt.Errorf("could not read Bing API key from %s: %s", file, err))
	}
	return strings.TrimSpace(string(content))
}
