package main

import (
	"fmt"
	"gokitchenmood/lampen"
	"html/template"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	title := "moodlights"
	p := &lampen.Lampen{}
	err := p.LoadLampValues(title)
	if err != nil {
		log.Printf("Error loading Config File")
		p = &lampen.Lampen{
			ValueLampe0x1: "0",
			ValueLampe0x2: "0",
			ValueLampe0x3: "0",
			ValueLampe0x4: "0",
			ValueLampe0x5: "0",
			ValueLampe0x6: "0",
			ValueLampe0x7: "0",
			ValueLampe0x8: "0",
			ValueLampe0x9: "0",
			ValueLampe0xA: "0"}
	}
	t, _ := template.ParseFiles("template.html")
	t.Execute(w, p)

}

func savehandler(w http.ResponseWriter, r *http.Request) {
	p := &lampen.Lampen{}
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
	err := p.WriteLampValues("moodlights")
	if err != nil {
		log.Fatal(err)
	}
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
