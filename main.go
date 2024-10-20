package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

type unixTimestamp struct {
	Epoch       int64
	Text        string
	Milli       bool
	Valid       bool
	Local       string
	UTC         string
	JustStarted bool
}

func loadPage(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	ts := unixTimestamp{
		Epoch: now.Unix(),
		Text:  now.Format(time.RFC3339),
		Milli: false,
		Valid: true,
		Local: now.Format(time.RFC3339),
		UTC:   now.UTC().Format(time.RFC3339),
	}
	tmpl := template.Must(template.ParseFiles("index.gohtml"))
	err := tmpl.Execute(w, ts)
	if err != nil {
		return
	}
}

func epochHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("epochHandler")
	seconds, err := strconv.Atoi(r.PostFormValue("seconds"))
	result := unixTimestamp{Text: "Invalid input - must be a positive integer", Valid: true}
	if err != nil {
		result.Valid = false
	}
	s := int64(seconds)
	if s < 0 {
		result.Valid = false
	}
	digits := int(math.Ceil(math.Log10(float64(s))))
	var unixTime time.Time
	if digits > 10 {
		unixTime = time.UnixMilli(s)
		result.Milli = true
		result.Text = "Assuming timestamp is in milliseconds"
	} else {
		unixTime = time.Unix(s, 0)
	}
	result.Epoch = unixTime.Unix()
	format := r.PostFormValue("result-format")
	var localTime, utcTime string
	switch format {
	case "RFC3339":
		localTime = unixTime.Format(time.RFC3339)
		utcTime = unixTime.UTC().Format(time.RFC3339)
	case "RFC1123":
		localTime = unixTime.Format(time.RFC1123)
		utcTime = unixTime.UTC().Format(time.RFC1123)
	case "RFC822":
		localTime = unixTime.Format(time.RFC822)
		utcTime = unixTime.UTC().Format(time.RFC822)
	case "RFC850":
		localTime = unixTime.Format(time.RFC850)
		utcTime = unixTime.UTC().Format(time.RFC850)
	case "ANSIC":
		localTime = unixTime.Format(time.ANSIC)
		utcTime = unixTime.UTC().Format(time.ANSIC)
	case "Unix":
		localTime = unixTime.Format(time.UnixDate)
		utcTime = unixTime.UTC().Format(time.UnixDate)
	}
	if result.Valid {
		result.Local = localTime
		result.UTC = utcTime
	}
	tmpl := template.Must(template.ParseFiles("index.gohtml"))
	err = tmpl.ExecuteTemplate(w, "table", result)
	if err != nil {
		return
	}
}

func timestampHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("timestampHandler")
	timestamp := r.PostFormValue("timestamp")
	format := r.PostFormValue("input-format")
	var ts time.Time
	var err error
	switch format {
	case "RFC3339":
		ts, err = time.Parse(time.RFC3339, timestamp)
	case "RFC1123":
		ts, err = time.Parse(time.RFC1123, timestamp)
	case "RFC822":
		ts, err = time.Parse(time.RFC822, timestamp)
	case "RFC850":
		ts, err = time.Parse(time.RFC850, timestamp)
	case "ANSIC":
		ts, err = time.Parse(time.ANSIC, timestamp)
	case "Unix":
		ts, err = time.Parse(time.UnixDate, timestamp)
	}
	result := unixTimestamp{Text: "Invalid timestamp", Valid: true}
	if err != nil {
		result.Valid = false
	}
	result.Epoch = ts.Unix()
	result.Text = timestamp
	fmt.Println(timestamp, format, ts, result)
	tmpl := template.Must(template.ParseFiles("index.gohtml"))
	err = tmpl.ExecuteTemplate(w, "timestamp", result)
	if err != nil {
		return
	}
}

func main() {

	fmt.Println("Starting app server")

	// define handlers
	http.HandleFunc("/", loadPage)
	http.HandleFunc("/epoch/", epochHandler)
	http.HandleFunc("/timestamp/", timestampHandler)

	// start server
	log.Fatal(http.ListenAndServe(":9000", nil))
}
