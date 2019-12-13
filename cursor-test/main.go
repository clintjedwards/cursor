package main

import (
	"fmt"
	"time"

	cursor "github.com/clintjedwards/cursor-sdk"
)

func start() error {
	time.Sleep(time.Second * 5)
	fmt.Println("Hello!")
	return nil
}

func task1() error {
	time.Sleep(time.Second * 5)
	fmt.Println("Task1")
	return nil
}

func task2() error {
	time.Sleep(time.Second * 5)
	fmt.Println("Task2")
	return nil
}

func task3() error {
	time.Sleep(time.Second * 5)
	fmt.Println("Task3")
	return nil
}

func end() error {
	time.Sleep(time.Second * 5)
	fmt.Println("Goodbye!")
	return nil
}

func main() {
	taskMap := map[string]cursor.Task{
		"start": cursor.Task{
			Name:        "Starting Task",
			Description: "This is a simple example task",
			Handler:     start,
			Children:    []string{"task1"},
		},
		"task1": cursor.Task{
			Name:        "Task One",
			Description: "This is another simple example task",
			Handler:     task1,
			Children:    []string{"task2", "task3"},
		},
		"task2": cursor.Task{
			Name:        "Task Two",
			Description: "This is another simple example task",
			Handler:     task2,
			Children:    []string{},
		},
		"task3": cursor.Task{
			Name:        "Task Three",
			Description: "This is another simple example task",
			Handler:     task3,
			Children:    []string{"end"},
		},
		"end": cursor.Task{
			Name:        "Ending Task",
			Description: "This is another simple example task",
			Handler:     end,
			Children:    []string{},
		},
	}

	cursor.Serve("cursor-test", "start", taskMap)
}
