package api

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/require"
)

func TestNewRailClient(t *testing.T) {
	client:=NewRailClient()
	stations,err:=client.GetStations()
	require.Nil(t, err)
	fmt.Println(stations)

}
