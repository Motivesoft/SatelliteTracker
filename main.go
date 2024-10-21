package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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
	if len(os.Args) < 2 {
		fmt.Println("Missing satellite ID")
		return
	}

	satelliteId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid satellite ID: '%s'. %v\n", os.Args[1], err)
		return
	}

	raw, err := performTle(satelliteId)
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
	return performTleLive(satelliteId)
}

func performTleDebug(satelliteId int) ([]byte, error) {
	// If debug, read from file
	data, err := os.ReadFile(fmt.Sprintf("./examples/tle-%d.json", satelliteId))
	if err != nil {
		return nil, err
	}

	return data, err
}

func performTleLive(satelliteId int) ([]byte, error) {
	// Base URL of the API
	baseURL := "https://api.n2yo.com/rest/v1/satellite"

	// Query parameters - will be used to add the apiKey and maybe other details
	queryParams := url.Values{}

	// Read apiKey and anything else relevant from a local file
	headers, err := readHeadersFromDotfile(".env")
	if err != nil {
		fmt.Println("Failed to read headers from .env:", err)
		return nil, err
	}

	// Put all values read from the dotfile as header entries
	for key, value := range headers {
		queryParams.Add(key, value)
	}

	// Construct the full URL
	fullURL := fmt.Sprintf("%s/tle/%d?%s", baseURL, satelliteId, queryParams.Encode())

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, err
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return nil, err
	}

	return body, nil
}

func printTle(tleStruct TLEStructure) {
	fmt.Printf("Satellite name     : %s\n", tleStruct.Info.SatelliteName)
	fmt.Printf("Satellite ID       : %d\n", tleStruct.Info.SatelliteId)
	fmt.Printf("Transactions Count : %d\n", tleStruct.Info.TransactionsCount)
	fmt.Printf("TLE                : \n%s\n", tleStruct.TLE)
}

// Read a dotfile, formatted as a property file, into a string map
func readHeadersFromDotfile(filename string) (map[string]string, error) {
	headers := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore empty lines and comment lines (starting with #)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return headers, nil
}
