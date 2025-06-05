package main

import (
	"fmt"
	"go1f/pkg/server"
	"go_final_project/pkg/db"
)

func main() {

	db.Init("scheduler.db")

	fmt.Println("Start server!")
	server.Run() //http.ListenAndServe(":7541", nil)

	fmt.Println("End of work!")
}
