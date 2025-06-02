package main

import (
	"fmt"
	"go_final_project/pkg/server"
	"net/http"
)

func main() {
	fmt.Println("Start server!")
	http.Handle("/", http.FileServer(http.Dir("./web")))
	server.Run()
	err := http.ListenAndServe(":7541", nil)
	if err != nil {
		panic(err)
	}

}
