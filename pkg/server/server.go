package server

import (
	"go_final_project/pkg/api"
)

func Run() {

	api.Init()
	//port := 7541
	//http.Handle("/", http.FileServer(http.Dir("web")))
	//return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
