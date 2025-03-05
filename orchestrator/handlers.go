package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Knetic/govaluate"
)

var (
	expressions = make(map[string]*Expression) // Хранилище выражений
	tasks       = make(map[string]*Task)       // Хранилище задач
	mutex       = &sync.Mutex{}
)

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	expr, err := govaluate.NewEvaluableExpression(req.Expression)
	if err != nil {
		http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
		return
	}

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	mutex.Lock()
	defer mutex.Unlock()
	expressions[id] = &Expression{
		ID:     id,
		Expr:   req.Expression,
		Status: "pending",
	}

	go evaluateExpression(id, expr)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func evaluateExpression(id string, expr *govaluate.EvaluableExpression) {
	result, err := expr.Evaluate(nil)
	if err != nil {
		mutex.Lock()
		expressions[id].Status = "error"
		mutex.Unlock()
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	expressions[id].Status = "done"
	expressions[id].Result = result.(float64)
}

func expressionsHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var exprs []map[string]interface{}
	for _, expr := range expressions {
		exprs = append(exprs, map[string]interface{}{
			"id":     expr.ID,
			"status": expr.Status,
			"result": expr.Result,
		})
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": exprs})
}

func expressionHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/expressions/"):]
	mutex.Lock()
	defer mutex.Unlock()

	expr, exists := expressions[id]
	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	switch r.Method {
	case http.MethodGet:
		for _, task := range tasks {
			if task.Arg1 != 0 && task.Arg2 != 0 && task.Operation != "" {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"task": task,
				})
				return
			}
		}
		http.Error(w, "No tasks available", http.StatusNotFound)

	case http.MethodPost:
		var req struct {
			ID     string  `json:"id"`
			Result float64 `json:"result"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
			return
		}

		task, exists := tasks[req.ID]
		if !exists {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		task.Arg1 = req.Result
		task.Arg2 = 0
		task.Operation = ""

		for _, expr := range expressions {
			if expr.ID == task.ID {
				expr.Status = "done"
				expr.Result = req.Result
				break
			}
		}

		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
