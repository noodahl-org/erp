package workflows

import (
	"github.com/go-playground/validator"
	"github.com/gocolly/colly"
	"github.com/noodahl-org/erp/api/clients/brave"
	"github.com/noodahl-org/erp/api/clients/ollama"
	"github.com/noodahl-org/erp/api/clients/postgres"
	"resty.dev/v3"
)

type WorkflowClient struct {
	resty      resty.Client
	webscraper *colly.Collector
	brave      brave.BraveClient
	db         postgres.DBClient
	valid      *validator.Validate
	ollama     ollama.OllamaClient
}

func NewWorkflowClient(opts ...func(*WorkflowClient)) *WorkflowClient {
	cl := &WorkflowClient{
		resty: *resty.New(),
		webscraper: colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:136.0) Gecko/20100101 Firefox/136.0"),
		),
	}
	for _, opt := range opts {
		opt(cl)
	}
	return cl
}

func WithOllamaClient(ollama ollama.OllamaClient) func(*WorkflowClient) {
	return func(w *WorkflowClient) {
		w.ollama = ollama
	}
}

func WithBraveClient(brave brave.BraveClient) func(*WorkflowClient) {
	return func(w *WorkflowClient) {
		w.brave = brave
	}
}

func WithValidator(v *validator.Validate) func(*WorkflowClient) {
	return func(w *WorkflowClient) {
		w.valid = v
	}
}

func WithDB(db postgres.DBClient) func(*WorkflowClient) {
	return func(wc *WorkflowClient) {
		wc.db = db
	}
}
