package main

import (
	"context"
	"time"

	worker "github.com/sanjeevchoubey/task-management/worker"
)

func main() {
	ctx := context.Background()
	connctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	//Handle abrupt closer
	worker.SetupCloseHandler()
	// this the main function which controlls all different go routines
	worker.Run(connctx)

}
