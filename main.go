package main

import (
	"context"
	"os"
	"time"

	"github.com/SanjeevChoubey/task-management/worker"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	ctx := context.Background()
	connctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	//Handle abrupt closer
	worker.SetupCloseHandler()
	// this the main function which controlls all different go routines
	worker.Run(connctx)

}
