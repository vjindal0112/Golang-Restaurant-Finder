package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

var results Results // struct that contains the results of our search

var tpl = template.Must(template.ParseFiles("index.html")) // creates the index.html page

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil) // create template
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String()) // gets the string from the URL and splits it up
	if err != nil {                     // check for errors
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()                                // get the parts at the end of the url
	searchKey := params.Get("q")                       // get the search Query
	location := params.Get("location")                 // get the inputted location
	searchKey = strings.ReplaceAll(searchKey, " ", "") // take out spaces
	location = strings.ReplaceAll(location, " ", "")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.yelp.com/v3/businesses/search?location="+location+"&term="+searchKey, nil) // create request for yelp api
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authorization", "Bearer C-ZhCFL4oGiZ0wiZMaTVpNXE5aHIDElgTXkxTSJUr9625yHYg8-rgaX7MZnraTafauK-qXAIgBnZ-8Xir8TjXUjfZc2VzxETKL0e5ihOIxDwKIHUTO5VLILBNqMKX3Yx") // set authorization token
	response, err := client.Do(req)                                                                                                                                            // do the request
	if err != nil {
		fmt.Println(err) // print the error if there is one
	}

	err = json.NewDecoder(response.Body).Decode(&results) // parse the response
	fmt.Println(results.Businesses)                       // print the businesses
	err = tpl.Execute(w, results)                         // run the template, passing in the results from API

	fmt.Println("Search Query is: ", searchKey)
	fmt.Println("Location is: ", location)
}

func main() {
	fmt.Println("App Started")

	mux := http.NewServeMux() // helps to call the correct handler based on the URL

	fs := http.FileServer(http.Dir("assets")) // put the "static" files on the server
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", indexHandler)        // what function to call on main page
	mux.HandleFunc("/search", searchHandler) // function to call on search
	http.ListenAndServe(":3000", mux)        // start a localserver to run files
}

// structs to parse the incoming data
type Business struct {
	Rating     int    `json:"rating"`
	Price      string `json:"price"`
	Phone      string `json:"phone"`
	ID         string `json:"id"`
	Alias      string `json:"alias"`
	IsClosed   bool   `json:"is_closed"`
	Categories []struct {
		Alias string `json:"alias"`
		Title string `json:"title"`
	} `json:"categories"`
	ReviewCount int    `json:"review_count"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	ImageURL string `json:"image_url"`
	Location struct {
		City     string `json:"city"`
		Country  string `json:"country"`
		Address2 string `json:"address2"`
		Address3 string `json:"address3"`
		State    string `json:"state"`
		Address1 string `json:"address1"`
		ZipCode  string `json:"zip_code"`
	} `json:"location"`
	Distance     float64  `json:"distance"`
	Transactions []string `json:"transactions"`
}

type Results struct {
	Total      int        `json:"total"`
	Businesses []Business `json:"businesses"`
}
