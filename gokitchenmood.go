package main

import (
	"fmt"
	"gokitchenmood/lampen"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var filetowrite string

func handler(w http.ResponseWriter, r *http.Request) {
	title := "moodlights"
	p := &lampen.Lampen{}
	err := p.LoadLampValues(title)
	if err != nil {
		log.Printf("Error loading Config File")
		for i := range p.Values {
			p.Values[i] = "000000"

		}
	}
	t, _ := template.ParseFiles("template.html")
	t.Execute(w, p)

}

func savehandler(w http.ResponseWriter, r *http.Request) {
	p := &lampen.Lampen{}
	var broken bool
	for i := 0; i < 10; i++ {
		err := p.Parse(r.FormValue("Lampe"+strconv.Itoa(i)), i)
		if err != nil {
			broken = true
			break
		}
	}
	p.Port = filetowrite
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
		p.Send()
		fmt.Fprintf(w, "<h1>Erfolg</h1>"+
			"Moodlights wurden geändert!"+
			"<form action=\"/\" method=\"POST\">"+
			"<div><input type=\"submit\" value=\"Back\"></div>"+
			"</form>")
	}
}

func statichandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s File [-f] \n", os.Args[0])
		return
	}
	lampen.File = false
	filetowrite = os.Args[1]
	if len(os.Args) > 2 {
		if os.Args[2] == "-f" {
			lampen.File = true
		}
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/save", savehandler)
	http.HandleFunc("/static/", statichandler)
	http.ListenAndServe(":8080", nil)
}
