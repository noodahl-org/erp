package workflows

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	bravemodels "github.com/noodahl-org/erp/internal/clients/brave/models"
	"github.com/noodahl-org/erp/internal/models"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type EquipmentActivity struct {
	Equipment               models.Equipment
	UserManualSearchResults bravemodels.SearchResponse
}

func (c *WorkflowClient) NewEquipmentWorkflow(ctx workflow.Context, id uuid.UUID) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
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

	// Now that we have the equipment details, try to fetch a user manual
	userManualQuery := fmt.Sprintf("%s %s operator/owners manual", act.Equipment.Make, act.Equipment.Model)
	err = workflow.ExecuteActivity(actx, c.BraveSearchActivity, userManualQuery).Get(ctx, &act.UserManualSearchResults)
	if err != nil {
		return err
	}

	log.Println("Equipment:", act.Equipment)
	log.Println("User Manual Search Results:", act.UserManualSearchResults)

	return nil
}

func (e *WorkflowClient) FetchEquipmentActivity(ctx context.Context, id uuid.UUID) (*models.Equipment, error) {
	var result models.Equipment
	return &result, e.db.Fetch(ctx, &models.Equipment{ID: id.String()}, &result)
}

func (e *WorkflowClient) BraveSearchActivity(ctx context.Context, query string) (*bravemodels.SearchResponse, error) {
	var result bravemodels.SearchResponse
	result.Query = bravemodels.Query{
		Original: "HELLO THERE",
	}
	return &result, nil
}
