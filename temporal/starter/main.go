package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create temporal client", err)
	}
	defer c.Close()

	wfOpts := client.StartWorkflowOptions{
		TaskQueue: "default",
		ID:        "equip_test",
	}

	//test for new equipment workflow
	run, err := c.ExecuteWorkflow(context.Background(), wfOpts, "NewEquipmentWorkflow", uuid.MustParse("a9a05684-3c2a-4867-9f89-4a6c06880c2b"))
	if err != nil {
		log.Fatalln("unable to execute workflow", err)
	}

	fmt.Printf("%v:%v", run.GetID(), run.GetRunID())
}
