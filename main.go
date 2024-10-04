package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

fileName := "data"
	file, err := os.Open(fileName + ".txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 
	}
	defer file.Close() 
	Scanner := bufio.NewScanner(file)
	var fileLines []string
	for Scanner.Scan() {
		fileLines = append(fileLines, Scanner.Text())
	}
	if err := Scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

Mumford := Data{
artist : strings.Split(fileLines[0],","),
members : strings.Split(fileLines[1],","),
albums : strings.Split(fileLines[2],","),
	albumyears : strings.Split(fileLines[3],","),
	locations : strings.Split(fileLines[4],","),
	concertdates : strings.Split(fileLines[5],","),
}	

fmt.Println(Mumford.members)


}