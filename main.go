package main

import (
	"fmt"
	"strconv"
	"strings"

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
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	RelUrl       string   `json:"relations"`
}

type Urelles struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

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

func LoadUrelles(Url string) ([]Urelles, error) {
	response, err := http.Get(Url)
	if err != nil {
		fmt.Println("Error getting response from url")
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading response body")
	}

	// Unmarshal the data into a map of Urelles to extract index
	var data map[string][]Urelles
	errr := json.Unmarshal(body, &data)
	if errr != nil {
		fmt.Print("Error storing body in address of character(unnmarshal issue)")
	}

	// Extract index from the data map which is the type urelles struct
	// ok is the boolean whhich indicates which checks if index exists
	urelles, ok := data["index"]
	if !ok {
	 fmt.Println("error: key 'index' not found in API response")
	}

	return urelles, nil
}

func homepage(w http.ResponseWriter, r *http.Request) {
	// where form value is collected for artist name annd fed innto the relavent funnction
	//ie- artist = r.FormValue = mumford

	character, _ := LoadData("https://groupietrackers.herokuapp.com/api/artists")
	// for i := 0; i < 52; i++ {
	// 	character[i] = Data{
	// 		ID:           character[i].ID,
	// 		Image:        character[i].Image,
	// 		Artist:       character[i].Artist,
	// 		Members:      character[i].Members,
	// 		AlbumYear:    character[i].AlbumYear,
	// 		Album1:       character[i].Album1,
	// 		Locations:    character[i].Locations,
	// 		ConcertDates: character[i].ConcertDates,
	// 		RelUrl:       character[i].RelUrl,
	// 	}

	// }
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
	character, err := LoadData("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(wr, "Failed to load artist data", http.StatusInternalServerError)
		return
	}
	if len(character) == 0 {
		http.Error(wr, "No artist data available", http.StatusInternalServerError)
		return
	}

	charData, err := LoadUrelles("https://groupietrackers.herokuapp.com/api/relation")
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

	//fmt.Println("chhard[i],datesloc:", charData[i].DatesLocations)
	//fmt.Println("chhardat", charData)


	cdata := ""
	x:='1'
	for location, date := range charData[i].DatesLocations {
		cdata += string(x) +") " +strings.ReplaceAll(string(location),"-",", ") +": "+ strings.Join(date,", ") + "; " 
		x++
	}

	cdata = strings.ReplaceAll(cdata,"_"," ") 

//fmt.Println("datedata:",cdata)

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
