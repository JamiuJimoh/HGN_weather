package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ResPayload struct {
	ClientIP string `json:"client_ip"`
	Location string `json:"location"`
	Greeting string `json:"greeting"`
}

type WeatherData struct {
	Location WeatherLocation `json:"location"`
	Current  WeatherCurrent  `json:"current"`
}

type WeatherLocation struct {
	Region  string `json:"region"`
	Country string `json:"country"`
}

type WeatherCurrent struct {
	Temp float64 `json:"temp_c"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	p := os.Getenv("PORT")
	port := fmt.Sprintf(":%s", p)
	http.HandleFunc("/api/hello", helloHandler)
	log.Fatal(http.ListenAndServe(port, nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// ip := "197.210.8.121"
	ip := getIpAddr(r)
	name := r.URL.Query().Get("visitor_name")
	if len(name) < 1 {
		http.Error(w, "query param: visitor_name is required", http.StatusNotFound)
		return
	}

	data, err := getWeatherData(ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	greeting := fmt.Sprintf("Hello, %s!, the temperature is %.2f degrees Celcius in %s", name, data.Current.Temp, data.Location.Region)
	res := ResPayload{
		ClientIP: ip,
		Location: data.Location.Region,
		Greeting: greeting,
	}

	w.Header().Set("Content-type", "Application/json")
	payload, _ := json.Marshal(res)
	w.Write(payload)
}

func getIpAddr(r *http.Request) string {
	var userIP string
	if len(r.Header.Get("CF-Connecting-IP")) > 1 {
		userIP = r.Header.Get("CF-Connecting-IP")
	} else if len(r.Header.Get("X-Forwarded-For")) > 1 {
		userIP = r.Header.Get("X-Forwarded-For")
	} else if len(r.Header.Get("X-Real-IP")) > 1 {
		userIP = r.Header.Get("X-Real-IP")
	} else {
		userIP = r.RemoteAddr
		if strings.Contains(userIP, ":") {
			userIP = strings.Split(userIP, ":")[0]
		}
	}

	temp := strings.Split(userIP, ", ")
	if len(temp) == 0 {
		return ""
	}

	return temp[0]
}

func getWeatherData(ip string) (WeatherData, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, ip)

	r, err := http.Get(url)
	if err != nil {
		return WeatherData{}, err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return WeatherData{}, err
	}

	var parsed WeatherData
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return WeatherData{}, err
	}

	return parsed, err
}
