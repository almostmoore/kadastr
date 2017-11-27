package api_server

import (
	"net/http"
	"github.com/iamsalnikov/kadastr/api_server/messages"
	"encoding/json"
	"bytes"
	"errors"
	"net/url"
)

type Client struct {
	ApiServerAddress string
	client *http.Client
}

func NewClient(apiAddr string) *Client {
	return &Client{
		ApiServerAddress: apiAddr,
		client: &http.Client{},
	}
}

func (c *Client) GetParsingTasksList() ([]messages.ParsingTask, error) {
	tasks := make([]messages.ParsingTask, 0)

	req, err := http.NewRequest(http.MethodGet, c.ApiServerAddress + "/list-parsing", nil)
	if err != nil {
		return tasks, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return tasks, err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tasks)

	return tasks, err
}

func (c *Client) AddParsingTask(quarters []string) (messages.AddParsingTaskAnswer, error) {
	answer := messages.AddParsingTaskAnswer{}

	requestBody := bytes.NewBufferString("")
	encoder := json.NewEncoder(requestBody)
	err := encoder.Encode(&quarters)
	if err != nil {
		return answer, err
	}

	req, err := http.NewRequest(http.MethodPost, c.ApiServerAddress + "/add-parsing", requestBody)
	if err != nil {
		return answer, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return answer, err
	}

	if resp.StatusCode != http.StatusOK {
		return answer, errors.New("Server answer status is not OK")
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&answer)

	return answer, err
}

func (c *Client) FindFeature(quarter, square string) ([]messages.FindFeature, messages.Error, error) {
	features := make([]messages.FindFeature, 0)
	searchError := messages.Error{}

	urlValues := url.Values{
		"quarter": []string{quarter},
		"square": []string{square},
	}

	req, err := http.NewRequest(http.MethodGet, c.ApiServerAddress + "/search?" + urlValues.Encode(), nil)
	if err != nil {
		return features, searchError, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return features, searchError, err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err = decoder.Decode(&searchError)

		return features, searchError, err
	}

	err = decoder.Decode(&features)
	return features, searchError, err
}