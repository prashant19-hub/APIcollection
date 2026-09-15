package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// OpenWeatherMap API response ka structure
type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: weather-app <city-name>")
		fmt.Println("Example: weather-app Kanpur")
		return
	}

	city := os.Args[1]
	apiKey := "YOUR_API_KEY_HERE" // yahan apna OpenWeatherMap API key daalo

	weather, err := getWeather(city, apiKey)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	printWeather(weather)
}

func getWeather(city, apiKey string) (*WeatherResponse, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s",
		city,
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}

func printWeather(w *WeatherResponse) {
	fmt.Println("----------------------------")
	fmt.Printf("City: %s\n", w.Name)
	fmt.Printf("Temperature: %.1f°C\n", w.Main.Temp)
	fmt.Printf("Feels Like: %.1f°C\n", w.Main.FeelsLike)
	fmt.Printf("Humidity: %d%%\n", w.Main.Humidity)
	fmt.Printf("Wind Speed: %.1f m/s\n", w.Wind.Speed)
	if len(w.Weather) > 0 {
		fmt.Printf("Condition: %s\n", w.Weather[0].Description)
	}
	fmt.Println("----------------------------")
}