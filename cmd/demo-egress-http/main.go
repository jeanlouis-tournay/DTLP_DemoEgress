package main

import (
	"fmt"
	"net/http"

	"eurocontrol.io/demo/egress/pkg/api"
	"eurocontrol.io/demo/egress/pkg/autoconfig"
	"github.com/gorilla/mux"
)

type Configuration struct {
	RestPort int32 `value:"server.port|8000"`
}


func main() {
	config:=&Configuration{}
	err:=autoconfig.AutoConfigure(config)
	if err !=nil {
		panic(err)
	}
	router := mux.NewRouter()
	h:=api.NewRailAPI()
	h.AddRoute(router)
	server := &http.Server{Addr: fmt.Sprintf(":%d", config.RestPort), Handler: router}
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}