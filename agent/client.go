package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

func fetchTask() (*Task, error) {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch task: %s", resp.Status)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func executeTask(task *Task) float64 {
	resp, err := http.Post("http://localhost:8081/calculate", "application/json", nil)
	if err != nil {
		log.Println("Error sending task to calculator:", err)
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Calculator returned an error:", resp.Status)
		return 0
	}

	var result float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Error decoding calculator response:", err)
		return 0
	}

	return result
}

func sendResult(taskID string, result float64) error {
	resp, err := http.Post("http://localhost:8080/internal/task", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send result: %s", resp.Status)
	}

	return nil
}
