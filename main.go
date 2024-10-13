package main

import (
	"fmt"
	//"strconv"
	//"groupie/functions"
	"io"

	//"image/jpeg"
	//"image"
	"encoding/json"
	"net/http"

	//"strings"
	"text/template"
)

type Data struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Artist       string   `json:"name"`
	Members      []string `json:"members"`
	AlbumYear    int      `json:"creationDate"`
	Album1       string   `json:"firstAlbum"`
	Locations    []string `json:"locations"`
	ConcertDates []string `json:"concertDates"`
	RelUrl       string   `json:"relations"`
}

type Origin struct {
	Name string
	URL  string
}

func main() {

	// handler functions
	http.HandleFunc("/", homepage)
	http.HandleFunc("/result", result)
	// http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	http.ListenAndServe(":8080", nil)
}

func LoadData() ([]Data, error) {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error getting response from url")
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error cannot store get resposne in body")
	}

	var character []Data
	err = json.Unmarshal(body, &character)
	if err != nil {
		fmt.Print("Error storing body in address of character")
	}

	return character, nil
}

func homepage(w http.ResponseWriter, r *http.Request) {
	// where form value is collected for artist name annd fed innto the relavent funnction
	//ie- artist = r.FormValue = mumford

	character, _ := LoadData()
	for i := 0; i < 52; i++ {
		character[i] = Data{
			ID:           character[i].ID,
			Image:        character[i].Image,
			Artist:       character[i].Artist,
			Members:      character[i].Members,
			AlbumYear:    character[i].AlbumYear,
			Album1:       character[i].Album1,
			Locations:    character[i].Locations,
			ConcertDates: character[i].ConcertDates,
			RelUrl:       character[i].RelUrl,
		}

	}
	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Error parsing html", http.StatusInternalServerError)
		return
	}

	
	err = t.Execute(w, character)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}

}

func result(w http.ResponseWriter, r *http.Request) {
	artId := r.FormValue("artist")
	fmt.Println(r.FormValue("artist"))
	fmt.Println(artId)

	character, _ := LoadData()
	for i := 0; i < 52; i++ {
		character[i] = Data{
			ID:           character[i].ID,
			Image:        character[i].Image,
			Artist:       character[i].Artist,
			Members:      character[i].Members,
			AlbumYear:    character[i].AlbumYear,
			Album1:       character[i].Album1,
			Locations:    character[i].Locations,
			ConcertDates: character[i].ConcertDates,
			RelUrl:       character[i].RelUrl,
		}
	}

	
	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Error parsing html", http.StatusInternalServerError)
		return
	}


	err = t.Execute(w, character)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
