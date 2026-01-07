package main

import (
	"fmt"
	"go-todoList/http"
	"go-todoList/todo"
)

func main() {
	todoList := todo.NewList()
	httpHandles := http.NewHTTPHandlers(todoList)
	httpServer := http.NewHTTPServer(httpHandles)

	if err := httpServer.StartServer(); err != nil {
		fmt.Println("failed to start HTTP server", err)
	}
}
