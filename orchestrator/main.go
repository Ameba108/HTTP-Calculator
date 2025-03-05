package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	http.HandleFunc("/api/v1/expressions", expressionsHandler)
	http.HandleFunc("/api/v1/expressions/", expressionHandler)
	http.HandleFunc("/internal/task", taskHandler)

	log.Println("Orchestrator is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
