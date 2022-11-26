package zincsearch

import (
	"context"
	"errors"
	"fmt"
	zinc "github.com/zinclabs/sdk-go-zincsearch"
	"net/http"
	"os"
)

type ZincClient struct {
	Ctx    context.Context
	Client *zinc.APIClient
}

// NewClient creates a new client to the variable Client.
func NewClient(url, username, password string) (ZincClient, error) {
	ctx := context.WithValue(context.Background(), zinc.ContextBasicAuth, zinc.BasicAuth{
		UserName: username,
		Password: password,
	})

	configuration := zinc.NewConfiguration()
	configuration.Servers = zinc.ServerConfigurations{
		zinc.ServerConfiguration{
			URL: url,
		},
	}

	return ZincClient{
		Ctx:    ctx,
		Client: zinc.NewAPIClient(configuration),
	}, nil

}

func Init() (ZincClient, error) {

	zincClient, err := NewClient(
		os.Getenv("ELASTICSEARCH_URL"),
		os.Getenv("ELASTICSEARCH_USERNAME"),
		os.Getenv("ELASTICSEARCH_PASSWORD"),
	)
	if err != nil {
		return ZincClient{}, err
	}

	err = CreateIndexIfNotExist("requests", zincClient)
	if err != nil {
		return ZincClient{}, err
	}

	err = CreateIndexIfNotExist("transactions", zincClient)
	if err != nil {
		return ZincClient{}, err
	}

	return zincClient, nil
}

func CreateIndexIfNotExist(index string, zincClient ZincClient) error {
	_, r, err := zincClient.Client.Index.Exists(zincClient.Ctx, index).Execute()

	if err != nil && r.StatusCode != http.StatusNotFound {
		fmt.Println("err checking Index Exists", r)
		return err
	}

	if r.StatusCode == http.StatusNotFound {
		indexMeta := *zinc.NewMetaIndexSimple() // MetaIndexSimple | Index data
		indexMeta.SetName(index)
		resp, r, err := zincClient.Client.Index.Create(zincClient.Ctx).Data(indexMeta).Execute()

		if err != nil {
			fmt.Println("err Creating new Index", r)
			return err
		}
		if r.StatusCode != 200 {
			return errors.New(*resp.Message)
		}
	}
	return nil
}
