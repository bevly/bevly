package menu

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/bevly/bevly/fetch/menu/alepdf"
	"github.com/bevly/bevly/httpagent"
	"github.com/bevly/bevly/model"
)

const AlehouseDescription = "ale_houseDescription"

func init() {
	menuFetcherRegistry["ale_house"] = alehouseMenu
}

func alehouseMenu(provider model.MenuProvider) ([]model.Beverage, error) {
	agent := httpagent.New()
	doc, err := agent.GetDoc(provider.URL())
	if err != nil {
		log.Printf("alehouseMenu: Get(%s) failed: %s\n",
			provider.URL(), err)
		return nil, err
	}

	beverages, err := alehouseDrafts(agent, doc)
	if err != nil {
		log.Printf("alehouseMenu: Failed to parse menu: %s\n", err)
		return nil, err
	}

	log.Printf("alehouseMenu: parsed %d beverages from %s\n",
		len(beverages), provider.URL())
	return beverages, nil
}

func alehouseDrafts(agent *httpagent.Agent, doc *goquery.Document) ([]model.Beverage, error) {
	pdfHref, ok := doc.Find("#ales a").Attr("href")
	if !ok {
		return nil, fmt.Errorf("no PDF link in alehouse HTML")
	}

	pdfFile, err := fetchPDF(agent, pdfHref)
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch PDF for Ale House: %s", err)
	}
	defer os.Remove(pdfFile)

	return alepdf.Parse(pdfFile)
}

func fetchPDF(agent *httpagent.Agent, href string) (string, error) {
	pdfTempF, err := ioutil.TempFile("", "alepdf")
	if err != nil {
		return "", err
	}
	pdfPath := pdfTempF.Name()
	pdfTempF.Close()

	_, err = agent.GetFile(href, pdfPath)
	if err != nil {
		os.Remove(pdfPath)
		return "", err
	}
	return pdfPath, nil
}
