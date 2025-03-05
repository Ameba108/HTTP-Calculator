package main

import "time"

type Expression struct {
	ID     string
	Expr   string
	Status string
	Result float64
}

type Task struct {
	ID            string
	Arg1          float64
	Arg2          float64
	Operation     string
	OperationTime time.Duration
}
