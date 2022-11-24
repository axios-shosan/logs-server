package elastic

import (
	"errors"
	"net/http"
	"os"
	"strings"

	es8 "github.com/elastic/go-elasticsearch/v8"
)

type Transport struct {
	Username string
	Password string
}

// RoundTrip implementation used to log in.
func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return http.DefaultClient.Do(req)
}

// NewClient creates a new client to the variable Client.
func NewClient(url, username, password string) (*es8.Client, error) {
	cfg := es8.Config{
		Addresses: []string{
			url,
		},
		Transport: &Transport{
			Username: username,
			Password: password,
		},
	}
	return es8.NewClient(cfg)
}

func Init() (*es8.Client, error) {

	client, err := NewClient(
		os.Getenv("ELASTICSEARCH_URL"),
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
	)
	if err != nil {
		return &es8.Client{}, err
	}

	err = CreateIndexIfNotExist("requests", client)
	if err != nil {
		return &es8.Client{}, err
	}

	err = CreateIndexIfNotExist("transactions", client)
	if err != nil {
		return &es8.Client{}, err
	}

	return client, nil
}

func CreateIndexIfNotExist(index string, client *es8.Client) error {

	resp, err := client.Indices.Exists([]string{index})
	if err != nil {
		return err
	}

	
	if resp.StatusCode == http.StatusNotFound {
		body := `
		{
			"settings" : {
				"analysis" : {
					"analyzer" : {
						"default" : {
							"tokenizer" : "standard",
								"filter" : ["asciifolding", "lowercase"]
						}
					}
				}
			}
		}`
		resp, err := client.Indices.Create(
			index,
			client.Indices.Create.WithBody(strings.NewReader(body)),
		)
		if err != nil {
			return err
		}
		if resp.IsError() {
			return errors.New(resp.String())
		}
	}
	return nil
}
