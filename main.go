package main

import (
	"fmt"
	"groupie/functions"
	"net/http"
	"strings"
	"text/template"
)


type Data struct {
	Artist      []string
    Members     []string
    Albums      []string
    AlbumYears  []string
    Locations   []string
    ConcertDates []string
}


func main () {

// handler functions
http.HandleFunc("/", homepage)
http.HandleFunc("/result", result)
http.ListenAndServe(":8080", nil)

}
func homepage (w http.ResponseWriter, r *http.Request) {
// where form value is collected for artist name annd fed innto the relavent funnction
//ie- artist = r.FormValue = mumford
t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Error parsing html", http.StatusInternalServerError)
		return
	}

	// execute the HTML template
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

// different function for each artist 
func result(w http.ResponseWriter, r *http.Request) {
	fileName := "artists/"+r.FormValue("artist")
fmt.Println(fileName)
	fileLines := functions.Read(fileName) 

name := Data{
Artist : strings.Split(fileLines[0],","),
Members : strings.Split(fileLines[1],","),
Albums : strings.Split(fileLines[2],","),
	AlbumYears : strings.Split(fileLines[3],","),
	Locations : strings.Split(fileLines[4],","),
	ConcertDates : strings.Split(fileLines[5],","),
}	

fmt.Println(name.Albums)
	// Parse the HTML template again for the resultpage
	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Error parsing html", http.StatusInternalServerError)
		return
	}

	// Render the template with the result
	err = t.Execute(w, name)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
