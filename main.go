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
	SatelliteId       int    `json:"satid"`
	SatelliteName     string `json:"satname"`
	TransactionsCount int    `json:"transactionscount"`
	PassesCount       int    `json:"passescount"`
}

type TLEStructure struct {
	Info InfoStructure `json:"info"`
	TLE  string        `json:"tle"`
}

type PassStructure struct {
	StartAz         float64 `json:"startAz"`
	StartAzCompass  string  `json:"startAzCompass"`
	StartEl         float64 `json:"startEl"`
	StartUTC        int64   `json:"startUTC"`
	MaxAz           float64 `json:"maxAz"`
	MaxAzCompass    string  `json:"maxAzCompass"`
	MaxEl           float64 `json:"maxEl"`
	MaxUTC          int64   `json:"maxUTC"`
	EndAz           float64 `json:"endAz"`
	EndAzCompass    string  `json:"endAzCompass"`
	EndEl           float64 `json:"endEl"`
	EndUTC          int64   `json:"endUTC"`
	Mag             float64 `json:"mag"`
	Duration        int     `json:"duration"`
	StartVisibility int64   `json:"startVisibility"`
}

type VisualPassesStructure struct {
	Info   InfoStructure   `json:"info"`
	Passes []PassStructure `json:"passes"`
}

var DEBUG bool

func main() {
	DEBUG = true

	if len(os.Args) < 2 {
		fmt.Println("Missing satellite ID")
		return
	}

	satelliteId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid satellite ID: '%s'. %v\n", os.Args[1], err)
		return
	}

	// Two Line Elements

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

	// Visual Passes

	raw, err = performVisualPasses(satelliteId)
	if err != nil {
		fmt.Println("Error performing Visual Passes:", err)
		return
	}

	// Unmarshal the JSON data into the struct
	var visualPassesStruct VisualPassesStructure
	err = json.Unmarshal(raw, &visualPassesStruct)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	printVisualPasses(visualPassesStruct)
}

func performVisualPasses(satelliteId int) ([]byte, error) {
	if DEBUG {
		return performVisualPassesDebug(satelliteId)
	}

	return performVisualPassesLive(satelliteId)
}

func performVisualPassesDebug(satelliteId int) ([]byte, error) {
	// If debug, read from file
	data, err := os.ReadFile(fmt.Sprintf("./examples/visualpasses-%d.json", satelliteId))
	if err != nil {
		return nil, err
	}

	return data, err
}

func performVisualPassesLive(satelliteId int) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func performTle(satelliteId int) ([]byte, error) {
	if DEBUG {
		return performTleDebug(satelliteId)
	}

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
	env, err := readHeadersFromDotfile(".env")
	if err != nil {
		fmt.Println("Failed to read headers from .env:", err)
		return nil, err
	}

	// Put all values read from the dotfile as header entries
	for key, value := range env {
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

func printTle(structure TLEStructure) {
	fmt.Printf("Satellite name     : %s\n", structure.Info.SatelliteName)
	fmt.Printf("Satellite ID       : %d\n", structure.Info.SatelliteId)
	fmt.Printf("Transactions Count : %d\n", structure.Info.TransactionsCount)
	fmt.Printf("TLE                : \n%s\n", structure.TLE)
}

func printVisualPasses(structure VisualPassesStructure) {
	fmt.Printf("Satellite name     : %s\n", structure.Info.SatelliteName)
	fmt.Printf("Satellite ID       : %d\n", structure.Info.SatelliteId)
	fmt.Printf("Transactions Count : %d\n", structure.Info.TransactionsCount)
	fmt.Printf("Passes Count       : %d\n", structure.Info.PassesCount)

	for i := 0; i < structure.Info.PassesCount; i++ {
		fmt.Printf("Pass %2d:\n", i)
		fmt.Printf("  StartUTC         : %d\n", structure.Passes[i].StartUTC)
		fmt.Printf("  StartAz          : %f\n", structure.Passes[i].StartAz)
		fmt.Printf("  StartAzCompass   : %s\n", structure.Passes[i].StartAzCompass)
		fmt.Printf("  StartEl          : %f\n", structure.Passes[i].StartEl)
		fmt.Printf("  MaxUTC           : %d\n", structure.Passes[i].MaxUTC)
		fmt.Printf("  MaxAz            : %f\n", structure.Passes[i].MaxAz)
		fmt.Printf("  MaxAzCompass     : %s\n", structure.Passes[i].MaxAzCompass)
		fmt.Printf("  MaxEl            : %f\n", structure.Passes[i].MaxEl)
		fmt.Printf("  EndUTC           : %d\n", structure.Passes[i].EndUTC)
		fmt.Printf("  EndAz            : %f\n", structure.Passes[i].EndAz)
		fmt.Printf("  EndAzCompass     : %s\n", structure.Passes[i].EndAzCompass)
		fmt.Printf("  EndEl            : %f\n", structure.Passes[i].EndEl)
		fmt.Printf("  Mag              : %f\n", structure.Passes[i].Mag)
		fmt.Printf("  Duration         : %d\n", structure.Passes[i].Duration)
		fmt.Printf("  Start Visibility : %d\n", structure.Passes[i].StartVisibility)
	}
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
