package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type WeatherRequest struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type WeatherResponse struct {
	Id    int `json:"id"`
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

var randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
var client = &http.Client{}

func main() {
	for {
		requestData := WeatherRequest{
			Water: generateRandomValue(1, 20),
			Wind:  generateRandomValue(1, 20),
		}

		responseData, err := sendToAPI(requestData)
		if err != nil {
			log.Panicf("Error processing POST request: %s", err)
		}

		waterStatus, windStatus := determineStatus(responseData)

		fmt.Printf("%+v\n", responseData)
		fmt.Printf("Status Water: %s\n", waterStatus)
		fmt.Printf("Status Wind : %s\n", windStatus)
		fmt.Println()

		time.Sleep(time.Second * 2)
	}
}

func generateRandomValue(min, max int) (value int) {
	value = randomGenerator.Intn((max-min)+1) + min
	return
}

func sendToAPI(requestData WeatherRequest) (responseData WeatherResponse, err error) {
	apiUrl := "https://jsonplaceholder.typicode.com/posts"

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		err = fmt.Errorf("failed to marshal request data: %w", err)
		return
	}

	request, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		err = fmt.Errorf("failed to create http request: %w", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("failed to perform http request: %w", err)
		return
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("failed to read response body: %w", err)
		return
	}

	err = json.Unmarshal(responseBody, &responseData)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal response data: %w", err)
		return
	}

	return
}

func determineStatus(weatherData WeatherResponse) (waterStatus, windStatus string) {
	switch {
	case weatherData.Water > 8:
		waterStatus = "Bahaya"
	case weatherData.Water >= 6:
		waterStatus = "Siaga"
	default:
		waterStatus = "Aman"
	}

	switch {
	case weatherData.Wind > 15:
		windStatus = "Bahaya"
	case weatherData.Wind >= 7:
		windStatus = "Siaga"
	default:
		windStatus = "Aman"
	}

	return
}
