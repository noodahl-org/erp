package workflows

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	bravemodels "github.com/noodahl-org/erp/api/clients/brave/models"
	ollamamodels "github.com/noodahl-org/erp/api/clients/ollama/models"
	"github.com/noodahl-org/erp/api/models"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type EquipmentActivity struct {
	Equipment               models.Equipment
	Components              []models.EquipmentComponent
	UserManualSearchResults bravemodels.SearchResponse
}

func (c *WorkflowClient) NewEquipmentWorkflow(ctx workflow.Context, id uuid.UUID) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 3 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 2,
		},
	}

	actx := workflow.WithActivityOptions(ctx, options)
	act := &EquipmentActivity{
		Equipment: models.Equipment{
			ID: id.String(),
		},
		UserManualSearchResults: bravemodels.SearchResponse{},
	}

	err := workflow.ExecuteActivity(actx, c.FetchEquipmentActivity, id).Get(ctx, &act.Equipment)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(actx, c.FetchEquipmentComponentsActivity, id).Get(ctx, &act.Components)
	if err != nil {
		return err
	}

	if len(act.Components) == 0 {
		err = workflow.ExecuteActivity(actx, c.GenerateEquipmentComponentsActivity, act.Equipment).Get(ctx, &act.Components)
		if err != nil {
			return err
		}
	}

	log.Println("Equipment:", act.Equipment)
	log.Println("Components", act.Components)

	return nil
}

func (e *WorkflowClient) FetchEquipmentActivity(ctx context.Context, id uuid.UUID) (*models.Equipment, error) {
	var result models.Equipment
	return &result, e.db.Fetch(ctx, &models.Equipment{ID: id.String()}, &result)
}

func (e *WorkflowClient) FetchEquipmentComponentsActivity(ctx context.Context, id uuid.UUID) ([]models.EquipmentComponent, error) {
	results := []models.EquipmentComponent{}
	return results, e.db.FetchMany(ctx, models.EquipmentComponent{EquipmentID: id.String()}, &results)
}

func (e *WorkflowClient) GenerateEquipmentComponentsActivity(ctx context.Context, equipment models.Equipment) ([]models.EquipmentComponent, error) {
	results := []models.EquipmentComponent{}
	query := fmt.Sprintf(`describe the the major components of %s %s. Respond using JSON`,
		equipment.Make, equipment.Model)
	format := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"components": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"required": []string{"components"},
	}
	llmData, err := e.ollama.Generate(ctx, query, format)
	if err != nil {
		return nil, err
	}

	llmResponse := ollamamodels.OllamaResponse{}
	if err := json.Unmarshal(llmData, &llmResponse); err != nil {
		return nil, err
	}
	compactBuf := bytes.NewBuffer([]byte{})
	if err := json.Compact(compactBuf, []byte(llmResponse.Response)); err != nil {
		return nil, err
	}

	var receiver struct {
		Components []string `json:"components"`
	}

	if err = json.Unmarshal(compactBuf.Bytes(), &receiver); err != nil {
		return nil, err
	}

	if len(receiver.Components) == 0 {
		return nil, errors.New("unable to generate components.")
	}

	for _, component := range receiver.Components {
		results = append(results, models.EquipmentComponent{
			EquipmentID: equipment.ID,
			Name:        component,
		})
	}

	return results, nil
}

func (e *WorkflowClient) BraveSearchActivity(ctx context.Context, query string) (*bravemodels.SearchResponse, error) {
	result, err := e.brave.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
