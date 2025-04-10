package brave

import (
	"context"
	"github.com/noodahl-org/erp/api/clients/brave/models"
	"resty.dev/v3"
)

type BraveClient interface {
	Search(ctx context.Context, query string) (models.SearchResponse, error)
}

type braveClient struct {
	client *resty.Client
}

func NewBraveClient(opts ...func(*braveClient) *braveClient) BraveClient {
	client := resty.New()
	client.SetBaseURL("https://api.search.brave.com/res/v1/web")
	b := &braveClient{
		client: client,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func WithAPIKey(key string) func(*braveClient) *braveClient {
	return func(b *braveClient) *braveClient {
		b.client.SetHeader("X-Subscription-Token", key)
		return b // Return the client after setting the header
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
