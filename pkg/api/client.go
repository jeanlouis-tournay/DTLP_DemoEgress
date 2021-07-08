package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type RailClient struct {
	client *http.Client
}


func NewRailClient() RailClient {
	return RailClient{
		client: &http.Client{
			Timeout: 5 *time.Second,
		},
	}
}

func (r RailClient) getData(url string) (map[string]interface{},error) {
	res,err:=r.client.Get(url)
	if err !=nil {
		return nil,err
	}
	if res.StatusCode > 399 {
		return nil,err
	}
	v:=make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&v)
	if err !=nil {
		return nil,err
	}
	return v,nil
}

func (r RailClient) GetStations() (map[string]interface{},error) {
	url:="https://api.irail.be/stations/?format=json"
	return r.getData(url)
}

