package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

func Start(s *http.Server) {
	log.Info("Starting HTTP server", "host", "localhost", "port", 8080)

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Could not start server", "error", err)
	}

}

func newServer() *http.Server {

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	fmt.Println(fs)

	mux.Handle("GET /", fs)

	mux.HandleFunc("GET /health", healthCheck)
	mux.HandleFunc("GET /distanceToLa/{ip}", handleIPDirection)
	// mux.HandleFunc("POST /email", handleEmail)
	// mux.HandleFunc("GET /{page}", subPageHandler)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return s

}

func main() {
	srv := newServer()
	Start(srv)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("healthy")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("HEALTHY"))
}

type ClientLocationResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

type DistanceAndDirection struct {
	Distance  float64 `json:"distance"`
	Direction string  `json:"direction"`
}

func handleIPDirection(w http.ResponseWriter, r *http.Request) {
	LALatRads := 34.052235 * math.Pi / 180
	LALonRads := -118.243683 * math.Pi / 180

	clientIp := r.PathValue("ip")
	valid := validateIPAddress(clientIp)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// fmt.Println(valid)

	clientLatLon, err := http.Get(fmt.Sprint("http://ip-api.com/json/", clientIp))
	clientLocation := ClientLocationResponse{}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer clientLatLon.Body.Close()

	body, err := io.ReadAll(clientLatLon.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &clientLocation); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if clientLocation.Status == "fail" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	clientLocation.Lat *= (math.Pi / 180)
	clientLocation.Lon *= (math.Pi / 180)

	deltaLon := LALonRads - clientLocation.Lon

	y := math.Sin(deltaLon) * math.Cos(clientLocation.Lat)
	// x := math.Cos(LALatRads)*math.Sin(clientLocation.Lat) - math.Sin(LALatRads)*math.Cos(clientLocation.Lat)*math.Cos(deltaLon)
	x := math.Cos(clientLocation.Lat)*math.Sin(LALatRads) - math.Sin(clientLocation.Lat)*math.Cos(LALatRads)*math.Cos(deltaLon)

	bearing := math.Atan2(y, x) * (180 / math.Pi)
	bearing = 360 - math.Abs(bearing)
	directions := []string{"north", "north east", "east", "south east", "south", "south west", "west", "north west"}

	index := int(math.Floor(bearing/45)) % 8

	km := math.Acos(math.Sin(clientLocation.Lat)*math.Sin(LALatRads)+(math.Cos(clientLocation.Lat)*math.Cos(LALatRads)*math.Cos(deltaLon))) * 6371

	fmt.Println(km, directions[index])
	distance := DistanceAndDirection{
		Distance:  km,
		Direction: directions[index],
	}

	jsonResp, err := json.Marshal(distance)

	if err != nil {
		w.Write([]byte("failed to get latitude and longitude"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(jsonResp)

}

func validateIPAddress(ip string) bool {
	substrings := strings.Split(ip, ".")
	if len(substrings) < 3 {
		return false
	}
	for _, value := range substrings {
		num, err := strconv.Atoi(value)
		if err != nil || num < 0 || num > 255 {
			return false
		}

	}

	return true
}

// type Email struct {
// 	From string `json:"from"`
// 	Name string `json:"name"`
// 	Msg  string `json:"msg"`
// }

// func handleEmail(w http.ResponseWriter, r *http.Request) {
// 	var email Email

// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()

// 	err := decoder.Decode(&email)

// 	if err != nil {
// 		fmt.Println("Error", err)
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return

// 	}
// 	_, err = mail.ParseAddress(email.From)

// 	if err != nil || email.Msg == "" {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	msg := []byte("Subject: ***SENT FROM WEBSITE***\r\n" + "\r\n" + "SENDER: " + email.From + "\r\n\n\n\n" + "NAME: " + email.Name + "\r\n\n\n\n" + "MESSAGE: " + email.Msg)

// 	password := os.Getenv("EMAIL_PASSWORD")
// 	fromAddress := os.Getenv("EMAIL_ADDRESS")
// 	personalEmail := os.Getenv("PERSONAL_EMAIL")

// 	auth := smtp.PlainAuth("", fromAddress, password, "smtp.gmail.com")

// 	to := []string{personalEmail}

// 	err = smtp.SendMail("smtp.gmail.com:587", auth, fromAddress, to, msg)

// 	if err != nil {

// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return

// 	}
// }
