package duckduckgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/bevly/bevly/html"
	"github.com/bevly/bevly/httpagent"
	"github.com/bevly/bevly/throttle"
	"github.com/bevly/bevly/websearch"
)

const SearchBaseURL = "https://duckduckgo.com/d.js"

var Throttle = throttle.Default(SearchBaseURL)

type DuckSearch struct {
	BaseURL  string
	Throttle *throttle.Throttle
}

func DefaultSearch() websearch.Search {
	return &DuckSearch{BaseURL: SearchBaseURL, Throttle: Throttle}
}

func SearchWithURL(baseURL string) websearch.Search {
	return &DuckSearch{BaseURL: baseURL}
}

func (s *DuckSearch) Search(terms string) ([]websearch.Result, error) {
	if s.Throttle != nil {
		s.Throttle.DelayInvocation()
	}
	response, err := httpagent.Agent().Get(s.SearchURL(terms))
	if err != nil {
		return nil, err
	}
	return extractJsonResults(response)
}

func (s *DuckSearch) SearchURL(terms string) string {
	p := func(name, value string) string {
		return url.Values{name: {value}}.Encode()
	}
	return s.BaseURL + "?" + p("q", terms) + "&" + p("l", "us-en") + "&" +
		p("p", "1") + "&" + p("s", "0")
}

func extractJsonResults(res *http.Response) ([]websearch.Result, error) {
	defer res.Body.Close()
	utf8Bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	json, err := extractDuckDuckGoJson(utf8Bytes)
	if err != nil {
		return nil, err
	}
	if json == nil {
		return nil, errors.New(fmt.Sprintf("no search payload in: %s", string(utf8Bytes)))
	}

	log.Printf("jsonResults size:%d\n", len(json))
	return jsonResults(json), nil
}

func jsonResults(json []interface{}) []websearch.Result {
	results := []websearch.Result{}
	for _, jsonResult := range json {
		jsonMap, ok := jsonResult.(map[string]interface{})
		if ok {
			if jsonMap["u"] != nil && jsonMap["t"] != nil {
				url := jsonMap["u"].(string)
				text := html.Text(jsonMap["t"].(string))
				if url != "" && text != "" {
					results = append(results, websearch.Result{
						URL:  url,
						Text: text,
					})
				}
			}
		}
	}
	return results
}

func extractDuckDuckGoJson(utf8Bytes []byte) ([]interface{}, error) {
	return duckJson(duckJsonBytes(utf8Bytes))
}

var rDuckJsonPayload = regexp.MustCompile(`nrn.*?(\[.*\])`)

func duckJsonBytes(utf8Bytes []byte) []byte {
	submatch := rDuckJsonPayload.FindSubmatch(utf8Bytes)
	if submatch == nil {
		return nil
	}
	return submatch[1]
}

func duckJson(jsonBytes []byte) ([]interface{}, error) {
	if jsonBytes == nil {
		return nil, nil
	}
	var res []interface{}
	err := json.Unmarshal(jsonBytes, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
