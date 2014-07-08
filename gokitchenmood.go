package main

import (
	"fmt"
	"gokitchenmood/lampen"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	title := "moodlights"
	p := &lampen.Lampen{}
	err := p.LoadLampValues(title)
	if err != nil {
		log.Printf("Error loading Config File")
		for i := range p.Values {
			p.Values[i] = "#000000"

		}
	}
	t, _ := template.ParseFiles("template.html")
	t.Execute(w, p)

}

func savehandler(w http.ResponseWriter, r *http.Request) {
	p := &lampen.Lampen{}
	var broken bool
	for i := range p.Values {
		err := p.Parse(r.FormValue("Lampe"+strconv.Itoa(i)), i)
		if err != nil {
			broken = true
			break
		}
	}
	if broken {
		fmt.Fprintf(w, "<h1>Fehler</h1>"+
			"Farben entsprechen nicht dem Format \"#CCCCCC\" oder \"CCCCCC\"!"+
			"<form action=\"/\" method=\"POST\">"+
			"<div><input type=\"submit\" value=\"Back\"></div>"+
			"</form>")
	} else {
		err := p.WriteLampValues("moodlights")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "<h1>Erfolg</h1>"+
			"Moodlights wurden geändert!"+
			"<form action=\"/\" method=\"POST\">"+
			"<div><input type=\"submit\" value=\"Back\"></div>"+
			"</form>")
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/save", savehandler)
	http.ListenAndServe(":8080", nil)
}
