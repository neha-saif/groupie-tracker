package main

import (
	"fmt"
	"reflect"
	"strconv"

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
	Members      []string `json:"members"` //
	AlbumYear    int      `json:"creationDate"`
	Album1       string   `json:"firstAlbum"`
	Locations    []string `json:"locations"`    //
	ConcertDates string   `json:"concertDates"` //
	RelUrl       string   `json:"relations"`
}

type Urelles struct {
	ID             int `json:"id"`
	DatesLocations [][]string `json:"datesLocations"`
}

type Final struct {
	ID           int    `json:"id"`
	Image        string `json:"image"`
	Artist       string `json:"name"`
	Members      string `json:"members"` //
	AlbumYear    string `json:"creationDate"`
	Album1       string `json:"firstAlbum"`
	Locations    string `json:"locations"`    //
	ConcertDates string `json:"concertDates"` //
}

func main() {

	// handler functions
	http.HandleFunc("/", homepage)
	http.HandleFunc("/result", result)
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))
	http.ListenAndServe(":8080", nil)
}

func LoadData(Url string) ([]Data, error) {
	response, err := http.Get(Url)
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

	character, _ := LoadData("https://groupietrackers.herokuapp.com/api/artists")
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
		return
	}

}

func result(wr http.ResponseWriter, r *http.Request) {
	artId := r.FormValue("artist")
	fmt.Println(r.FormValue("artist"))
	fmt.Println(artId)
	fmt.Println(reflect.TypeOf(artId))

	character, err := LoadData("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println("err:", err)
	}

	// i, err := strconv.Atoi(artId)
	// if err != nil {
	// 	http.Error(w, "Invalid artist ID", http.StatusBadRequest)
	// 	return
	// }

	iint, err := strconv.Atoi(artId)
	if err != nil {
		http.Error(wr, "Invalid artist ID", http.StatusBadRequest)
		return
	}
	i := iint - 1
	fmt.Println("i:", i)
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

	members := ""
	clocations := ""
	for j, _ := range character[i].Members {
		members += character[i].Members[j] + ", "
	}

	CharData, _ := LoadData("https://groupietrackers.herokuapp.com/api/relation")

	CharData[i] = Urelles{
		ID:             CharData[i].ID,
		DatesLocations: CharData[i].DatesLocations,
	}

	for j, _ := range CharData.DatesLocations {
		clocations += CharData.DatesLocations[j] + ", "
	}

	fmt.Println("concerrtd:", concertdates)
	FFinal := Final{
		ID:           character[i].ID,
		Image:        character[i].Image,
		Artist:       character[i].Artist,
		Members:      members,
		ConcertDates: clocations,
	}

	fmt.Println("charM:", character[i].Members[0])
	t, err := template.ParseFiles("result.html")
	if err != nil {
		http.Error(wr, "Error parsing html", http.StatusInternalServerError)
		return
	}

	err = t.Execute(wr, FFinal)
	if err != nil {
		http.Error(wr, "Error executing template ya", http.StatusInternalServerError)
		return
	}
}
