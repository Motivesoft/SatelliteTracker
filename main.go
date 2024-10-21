package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type InfoStructure struct {
	SatelliteName     string `json:"satname"`
	SatelliteId       int    `json:"satid"`
	TransactionsCount int    `json:"transactionscount"`
}

type TLEStructure struct {
	Info InfoStructure `json:"info"`
	TLE  string        `json:"tle"`
}

func main() {
	raw, err := performTle(25544)
	if err != nil {
		fmt.Println("Error performing TLE:", err)
		return
	}

	// Unmarshal the JSON data into the struct
	var tleStruct TLEStructure
	err = json.Unmarshal(raw, &tleStruct)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	printTle(tleStruct)
}

func performTle(satelliteId int) ([]byte, error) {
	// If debug, read from file
	data, err := os.ReadFile(fmt.Sprintf("./examples/tle-%d.json", satelliteId))
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	return data, err
}

func printTle(tleStruct TLEStructure) {
	fmt.Printf("Satellite name     : %s\n", tleStruct.Info.SatelliteName)
	fmt.Printf("Satellite ID       : %d\n", tleStruct.Info.SatelliteId)
	fmt.Printf("Transactions Count : %d\n", tleStruct.Info.TransactionsCount)
	fmt.Printf("TLE                : \n%s\n", tleStruct.TLE)
}
