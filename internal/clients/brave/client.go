package brave

import (
	"context"

	"github.com/noodahl-org/erp/internal/clients/brave/models"
	"resty.dev/v3"
)

type BraveClient interface {
}

type braveClient struct {
	client *resty.Client
}

func NewBraveClient(opts func(*braveClient) *braveClient) BraveClient {
	client := resty.New()
	client.SetBaseURL("https://api.search.brave.com/res/v1/web")
	return &braveClient{
		client: client,
	}
}

func (c *braveClient) Search(ctx context.Context, query string) (models.SearchResponse, error) {
	var result models.SearchResponse
	res, err := c.client.R().
		EnableTrace().
		SetQueryParams(map[string]string{
			"q":             query,
			"result_filter": "web",
		}).
		SetResult(&result).
		Get("/search")
	if err != nil {
		return result, err
	}
	return result, res.Err
}
