package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/noodahl-org/erp/api/clients/brave"
	"github.com/noodahl-org/erp/api/clients/ollama"
	"github.com/noodahl-org/erp/api/clients/postgres"
	"github.com/noodahl-org/erp/temporal/worker/conf"
	"github.com/noodahl-org/erp/temporal/worker/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var wfclient *workflows.WorkflowClient

func init() {
	os.Setenv("TEMPORAL_DEBUG", "true")
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	conf := &conf.TemporalConf{
		DbPort: dbPort,
		DbUser: os.Getenv("DB_USER"),
		DbPass: os.Getenv("DB_PASS"),
		DbHost: os.Getenv("DB_HOST"),
		DbName: os.Getenv("DB_NAME"),
	}

	valid := validator.New()
	//todo validate conf

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%v database=%s",
		conf.DbHost, conf.DbUser, conf.DbPass, conf.DbPort, conf.DbName)
	db, err := gorm.Open(driver.Open(dsn), nil)
	if err != nil {
		log.Fatalln(err)
	}

	b := brave.NewBraveClient(
		brave.WithAPIKey(os.Getenv("BRAVE_API_KEY")),
	)

	o := ollama.NewOllamaClient(
		ollama.WithBaseURL(os.Getenv("OLLAMA_URL")),
	)

	wfclient = workflows.NewWorkflowClient(
		workflows.WithValidator(valid),
		workflows.WithDB(postgres.NewPGDB(
			postgres.WithGormDB(db),
		)),
		workflows.WithBraveClient(b),
		workflows.WithOllamaClient(o),
	)
}

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create workflow client", err)
	}
	defer c.Close()

	w := worker.New(c, "default", worker.Options{})
	w.RegisterWorkflow(wfclient.NewEquipmentWorkflow)
	w.RegisterActivity(wfclient.FetchEquipmentActivity)
	w.RegisterActivity(wfclient.BraveSearchActivity)
	w.RegisterActivity(wfclient.FetchEquipmentComponentsActivity)
	w.RegisterActivity(wfclient.GenerateEquipmentComponentsActivity)

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln(err)
	}
}
