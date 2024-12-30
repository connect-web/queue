package main

import (
	"fmt"
	"github.com/google/uuid"
)

var (
	Tasks = map[string]struct{}{}
)

func CreateTask() (taskId string) {
	id := uuid.New()
	fmt.Printf("Generated UUID: %s\n", id.String())
	_, exists := Tasks[id.String()]
	for exists {
		id = uuid.New()
		fmt.Printf("Generated UUID: %s\n", id.String())
		_, exists = Tasks[id.String()]
	}

	// Adding task
	Tasks[id.String()] = struct{}{}
	return id.String()
}
