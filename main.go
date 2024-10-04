package main

import (
	"fmt"
	"groupie/functions"
	"net/http"
	"strings"
	"text/template"
)


type Data struct {
	artist []string
	members []string
	albums []string
	albumyears []string
	locations []string
	concertdates []string
}


func main () {

// handler functions
http.HandleFunc("/", homepage)
http.HandleFunc("/mumford", Mumford)
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
func Mumford(w http.ResponseWriter, r *http.Request) {
	fileName := "mumford"

	fileLines := functions.Read(fileName) 

Mumford := Data{
artist : strings.Split(fileLines[0],","),
members : strings.Split(fileLines[1],","),
albums : strings.Split(fileLines[2],","),
	albumyears : strings.Split(fileLines[3],","),
	locations : strings.Split(fileLines[4],","),
	concertdates : strings.Split(fileLines[5],","),
}	

fmt.Println(Mumford.members)
	// Parse the HTML template again for the resultpage
	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Error parsing html", http.StatusInternalServerError)
		return
	}

	// Render the template with the result
	err = t.Execute(w, Mumford)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
