package worker

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS.
func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Intrupt from OS,Closing the Task management System, Thank you for using this application")
		os.Exit(0)
	}()
}
