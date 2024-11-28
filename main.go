package main

import (
	"groupie/functions"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"fmt"
)

type Final struct {
	ID        int
	Image     string
	Artist    string
	Members   string
	AlbumYear int
	Album1    string
	Locations []string
}

func main() {
	// handler functions
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))
	http.HandleFunc("/", homepage)
	http.HandleFunc("/result", result)
	http.ListenAndServe(":8080", nil)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderErrorPage(w, 404)
		return
	}

	character, _ := functions.LoadData("https://groupietrackers.herokuapp.com/api/artists")

	t, err := template.ParseFiles("index.html")
	if err != nil {
		renderErrorPage(w, 500)
		return
	}

	err = t.Execute(w, character)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func result(wr http.ResponseWriter, r *http.Request) {
	// For resultpage the request is always POST not GET
	if r.Method != http.MethodPost {
		renderErrorPage(wr, 405)
		return
	}

	// if url is not for result page error handle
	if r.URL.Path != "/result" {
		renderErrorPage(wr, 404)
		return
	}
	artId := r.FormValue("artist")
	for _, ch := range artId {
		if ch != 10 && ch != 13 && (ch < 32 || ch > 126) {
			renderErrorPage(wr, 400)
			return
		}
	}
	iint, err := strconv.Atoi(artId)
	if err != nil || iint <= 0 || iint > 52 {
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

	var cdata []string
	x := '1'
	d := ""
	for location, date := range charData[i].DatesLocations {
		d = string(x) + ") " + strings.ReplaceAll(string(location), "-", ", ") + ": " + strings.Join(date, ", ")
		d = strings.ReplaceAll(d, "_", " ")
		cdata = append(cdata,d)
		x++
	}

	FFinal := Final{
		ID:        character[i].ID,
		Image:     character[i].Image,
		Artist:    character[i].Artist,
		Members:   members,
		AlbumYear: character[i].AlbumYear,
		Album1:    character[i].Album1,
		Locations:  cdata,
	}

	t, err := template.ParseFiles("result.html")
	if err != nil {
		renderErrorPage(wr, 500)
		return
	}

	errr := t.Execute(wr, FFinal)
	if errr != nil {
		http.Error(wr, "Error executing template", http.StatusInternalServerError)
	}
}

func renderErrorPage(w http.ResponseWriter, code int) {

	// Set the status coded
	if code == 500{
		w.WriteHeader(http.StatusInternalServerError)
	} else  if code == 404 {
		w.WriteHeader(http.StatusNotFound)
	} else if code == 400 {
		w.WriteHeader(http.StatusBadRequest)
	}

	// Generate ASCII art for the error code with the "Standard" style
	
	// Parse and render the custom 404 template
	t, err := template.ParseFiles(fmt.Sprintf("style/%d.html", code))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing %d HTML", code), http.StatusInternalServerError)
		return
	}

	// Render the template with the result
	err = t.Execute(w, nil )
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing %d template", code), http.StatusInternalServerError)
	}
}
