package rapi

import (
	"encoding/json"
	"github.com/almostmoore/kadastr/feature"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	baseUrl string
}

func NewClient() *Client {
	return &Client{
		baseUrl: "https://pkk5.rosreestr.ru",
	}
}

func (c *Client) GetFeature(number string) (feature.Entity, error) {
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}

	req, err := http.NewRequest("GET", c.baseUrl+"/api/features/1/"+number, nil)
	if err != nil {
		return feature.Entity{}, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return feature.Entity{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return feature.Entity{}, err
	}

	var answer featureAnswer
	err = json.Unmarshal(body, &answer)
	if err != nil {
		log.Printf("Тело ответа для номера %s:\n%s\n", number, string(body))
	}

	return answer.Feature.Attributes, err
}
