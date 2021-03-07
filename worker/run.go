package worker

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

//Wait group
var wg sync.WaitGroup

//Task ...:This struct is main struct to maintain tasks
type Task struct {
	ID           string
	IsCompleted  bool
	Status       string    // untouched, completed, failed, timeout
	CreationTime time.Time // when was the task created
	TaskData     string    // field containing data about the task
}

// create queue for passing task to workers
var intasks chan Task = make(chan Task, 20)
var outtasks chan Task = make(chan Task, 20)

//WaitTime ...: In order ot acheive timeout for few task, keeping this to 10 seconds
const WaitTime = 10

//Run ...This the main run function which start adder, executor and cleaner
func Run(ctx context.Context) {
	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg.Add(3)
	//Adder
	go Adder(workerCtx, intasks)
	//Executor
	go Executor(workerCtx, intasks, outtasks)
	//Cleaner
	go Cleaner(workerCtx, intasks, outtasks)

	wg.Wait()

}

//Executor ...This function will recives the task and execute them
// In order to demonstrate sending of email, just printing on console
func Executor(ctx context.Context, intasks, outtasks chan Task) {
	defer wg.Done()
	for {
		select {

		case task, ok := <-intasks:
			if !ok {
				log.Println("Channel is closed")
				return
			}
			rand.Seed(time.Now().UnixNano())
			if randBoolean() {
				// Adding a sleep to demonstrate that sending email is a blocking and may consume a millisecond
				time.Sleep(time.Millisecond)
				task.Status = "completed"
				task.IsCompleted = true
				log.Printf("Email Sent successfully for task id %s and task has been %s \n", task.ID, task.Status)
			} else {
				task.Status = "failed"
				log.Printf("Email sending has failed for task id  %s \n", task.ID)
			}
			outtasks <- task

		case <-ctx.Done():
			return
		}
	}
}

//Cleaner ...: This function will remove tasks from queue if its timeout or completed
// all other task status will be resubmiited to reprocess
func Cleaner(ctx context.Context, intasks, outtasks chan Task) {
	defer wg.Done()
	for {
		select {

		case task, ok := <-outtasks:
			if !ok {
				log.Println("Channel is closed")
				return
			}
			//timeout then remove from queue and log
			if time.Now().Sub(task.CreationTime).Seconds() > WaitTime {
				log.Printf("task %s has been time out and removed from queue \n", task.ID)
				//Completed job should be removed from queue
			} else if task.IsCompleted {
				fmt.Printf("task %s has been Completed  and removed from queue \n", task.ID)
				//All task which are not completed or timeout should be resubmitted to queue for processing
			} else {
				log.Printf("task %s has been re submitted in queue \n", task.ID)
				intasks <- task
			}

		case <-ctx.Done():
			return
		}
	}

}

//Adder ...
func Adder(ctx context.Context, tasks chan Task) {
	defer wg.Done()
	i := 0
	for {
		i++
		task := &Task{
			ID:           fmt.Sprintf("task %s", strconv.Itoa(i)),
			IsCompleted:  false,
			Status:       "untouched",
			CreationTime: time.Now(),
			TaskData:     fmt.Sprintf("Dummy data for task %s ", strconv.Itoa(i)),
		}
		tasks <- *task
	}
}

//Generate a randomBoolean
//30% task should be failed sending
func randBoolean() bool {
	return rand.Float32() < 0.7
}
