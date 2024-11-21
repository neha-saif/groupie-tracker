package main

import (
	"groupie/functions"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type Final struct {
	ID           int
	Image        string
	Artist       string
	Members      string
	AlbumYear    int
	Album1       string
	Locations    string
}

func main() {
	// handler functions
	http.HandleFunc("/", homepage)
	http.HandleFunc("/result", result)
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))
	http.ListenAndServe(":8080", nil)
}


func homepage(w http.ResponseWriter, r *http.Request) {

	character, _ := functions.LoadData("https://groupietrackers.herokuapp.com/api/artists")
	
	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Error parsing html", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, character)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func result(wr http.ResponseWriter, r *http.Request) {
	artId := r.FormValue("artist")
	iint, err := strconv.Atoi(artId)
	if err != nil || iint <= 0 {
		http.Error(wr, "Invalid artist ID", http.StatusBadRequest)
		return
	}
	i := iint - 1

	// Load artist data
	character, err := functions.LoadData("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(wr, "Failed to load artist data", http.StatusInternalServerError)
		return
	}
	
	if len(character) == 0 {
		http.Error(wr, "No artist data available", http.StatusInternalServerError)
		return
	}

	charData, err := functions.LoadUrelles("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		http.Error(wr, "Failed to load relations data", http.StatusInternalServerError)
		return
	}

	if len(charData) == 0 {
		http.Error(wr, "No data available", http.StatusInternalServerError)
		return
	}

	members := "No members available"
	if len(character[i].Members) > 0 {
		members = strings.Join(character[i].Members, ", ")
	}

	cdata := ""
	x:='1'
	for location, date := range charData[i].DatesLocations {
		cdata += string(x) +") " +strings.ReplaceAll(string(location),"-",", ") +": "+ strings.Join(date,", ") + "; " 
		x++
	}

	cdata = strings.ReplaceAll(cdata,"_"," ") 

FFinal := Final{
		ID:           character[i].ID,
		Image:        character[i].Image,
		Artist:       character[i].Artist,
		Members:      members,
		AlbumYear:    character[i].AlbumYear,
		Album1:       character[i].Album1,
		Locations:    cdata,
	}

	t, err := template.ParseFiles("result.html")
	if err != nil {
		http.Error(wr, "Error parsing result.html template", http.StatusInternalServerError)
		return
	}

	errr := t.Execute(wr, FFinal)
	if errr != nil {
		http.Error(wr, "Error executing template", http.StatusInternalServerError)
	}
}
