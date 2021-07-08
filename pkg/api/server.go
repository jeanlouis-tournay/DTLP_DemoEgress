package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type RailApi interface {
	AddRoute(router *mux.Router)
}


type railAPI struct {
	client RailClient
}


func NewRailAPI() RailApi {
	return &railAPI{NewRailClient()}
}

func (ra *railAPI) stations() func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET STATIONS")
		s,err:=ra.client.GetStations()
		if err !=nil {
			fmt.Printf("error %v \n",err)
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(s)
	}
}

func (ra *railAPI) AddRoute(router *mux.Router) {
	router.HandleFunc("/stations", ra.stations()).Methods("GET")
}