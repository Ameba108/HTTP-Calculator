package main

import (
	"log"
	"time"
)

func main() {
	for {
		task, err := fetchTask()
		if err != nil {
			log.Println("Error fetching task:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		result := executeTask(task)
		if err := sendResult(task.ID, result); err != nil {
			log.Println("Error sending result:", err)
		}
	}
}
