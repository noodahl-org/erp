package ollama

import (
	"context"

	"resty.dev/v3"
)

type OllamaClient interface {
	Generate(ctx context.Context, prompt string, format map[string]interface{}) ([]byte, error)
}

type ollamaClient struct {
	resty *resty.Client
}

func NewOllamaClient(opts ...func(*ollamaClient)) OllamaClient {
	o := &ollamaClient{
		resty: resty.New(),
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func WithBaseURL(baseURL string) func(*ollamaClient) {
	return func(o *ollamaClient) {
		o.resty.SetBaseURL(baseURL)
	}
}

func (o *ollamaClient) Generate(ctx context.Context, prompt string, format map[string]interface{}) ([]byte, error) {
	payload := map[string]interface{}{
		"model":      "gemma3:latest",
		"prompt":     prompt, // Also fixed this to use the parameter instead of the literal "prompt"
		"stream":     false,
		"max_tokens": "128",
		"format":     format,
	}
	res, err := o.resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("/api/generate")
	if err != nil {
		return nil, err
	}

	return res.Bytes(), res.Err
}
