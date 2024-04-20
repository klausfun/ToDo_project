package main

import (
	todo "ToDo_project"
	"log"
	//"github.com/zhashkevych/todo-app"
)

func main() {
	srv := new(todo.Server)
	if err := srv.Run("8080"); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
