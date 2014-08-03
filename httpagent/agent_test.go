package httpagent

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestHttpGet(t *testing.T) {
	agent := Agent()
	agent.ForceEncoding = "latin1"
	res, err := agent.Get("http://beer.friscogrille.com/")
	if err != nil {
		log.Printf("Error getting data from server: %s\n", err)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading from response: %s\n", err)
	}
	ioutil.WriteFile("beer.html", bytes, 0644)
}
