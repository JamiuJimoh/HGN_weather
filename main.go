package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

type ResPayload struct {
	ClientIP string `json:"client_ip"`
	Location string `json:"location"`
	Greeting string `json:"greeting"`
}

func main() {
	http.HandleFunc("/api/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIpAddr(r)
	name := r.URL.Query().Get("visitor_name")
	if len(name) < 1 {
		http.Error(w, "query param: visitor_name is required", http.StatusNotFound)
	}
	res := ResPayload{
		// ClientIP: "127.0.0.1",
		ClientIP: fmt.Sprintf("%s", ip),
		Location: "Lagos",
		Greeting: fmt.Sprintf("Hello, %s!, the temperature is 11 degrees Celcius in New York", name),
	}

	payload, err := json.Marshal(&res)
	if err != nil {
		return
	}

	w.Write(payload)
	// fmt.Fprintln(w, payload)
}
func getIpAddr(r *http.Request) string {
	var userIP string
	if len(r.Header.Get("CF-Connecting-IP")) > 1 {
		userIP = r.Header.Get("CF-Connecting-IP")
		fmt.Println(net.ParseIP(userIP))
	} else if len(r.Header.Get("X-Forwarded-For")) > 1 {
		userIP = r.Header.Get("X-Forwarded-For")
		fmt.Println(net.ParseIP(userIP))
	} else if len(r.Header.Get("X-Real-IP")) > 1 {
		userIP = r.Header.Get("X-Real-IP")
		fmt.Println(net.ParseIP(userIP))
	} else {
		userIP = r.RemoteAddr
		if strings.Contains(userIP, ":") {
			fmt.Println(net.ParseIP(strings.Split(userIP, ":")[0]))
		} else {
			fmt.Println(net.ParseIP(userIP))
		}
	}
	return userIP
}

// func getIpAddr(r *http.Request) string {
// 	clientIp := r.Header.Get("X-FORWARDED-FOR")
// 	if clientIp != "" {
// 		return clientIp
// 	}
// 	return r.RemoteAddr
// }
