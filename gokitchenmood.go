package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Page struct {
	Title         string
	ValueLampe0x1 string
	ValueLampe0x2 string
	ValueLampe0x3 string
	ValueLampe0x4 string
	ValueLampe0x5 string
	ValueLampe0x6 string
	ValueLampe0x7 string
	ValueLampe0x8 string
	ValueLampe0x9 string
	ValueLampe0xA string
}

func loadLampValues(title string) (*Page, error) {
	filename := title + ".json"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var p Page
	err = json.Unmarshal(body, &p)
	return &p, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	title := "moodlights"
	p, err := loadLampValues(title)
	if err != nil {
		p = &Page{Title: title}
	}
	t, _ := template.ParseFiles("template.html")
	t.Execute(w, p)
}

func savehandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "moodlights"}
	p.ValueLampe0x1 = r.FormValue("Lampe0x1")
	p.ValueLampe0x2 = r.FormValue("Lampe0x2")
	p.ValueLampe0x3 = r.FormValue("Lampe0x3")
	p.ValueLampe0x4 = r.FormValue("Lampe0x4")
	p.ValueLampe0x5 = r.FormValue("Lampe0x5")
	p.ValueLampe0x6 = r.FormValue("Lampe0x6")
	p.ValueLampe0x7 = r.FormValue("Lampe0x7")
	p.ValueLampe0x8 = r.FormValue("Lampe0x8")
	p.ValueLampe0x9 = r.FormValue("Lampe0x9")
	p.ValueLampe0xA = r.FormValue("Lampe0xA")
	b, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("moodlights.json", b, 0644)
	fmt.Fprintf(w, "<h1>Erfolg</h1>"+
		"Moodlights wurden geändert!"+
		"<form action=\"/\" method=\"POST\">"+
		"<div><input type=\"submit\" value=\"Back\"></div>"+
		"</form>")
	//	time.Sleep(20 * time.Second)
	//	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/save", savehandler)
	http.ListenAndServe(":8080", nil)
}
